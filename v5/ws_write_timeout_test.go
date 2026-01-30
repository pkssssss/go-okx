package okx

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

type writeDeadlineConn struct {
	net.Conn
	lastWriteDeadline atomic.Value // time.Time
}

func (c *writeDeadlineConn) SetWriteDeadline(t time.Time) error {
	c.lastWriteDeadline.Store(t)
	return c.Conn.SetWriteDeadline(t)
}

func (c *writeDeadlineConn) LastWriteDeadline() time.Time {
	if c == nil {
		return time.Time{}
	}
	if v := c.lastWriteDeadline.Load(); v != nil {
		if t, ok := v.(time.Time); ok {
			return t
		}
	}
	return time.Time{}
}

func TestWSClient_WriteDeadline_SetOnWriteJSONAndWriteText(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	type opReq struct {
		ID   string  `json:"id"`
		Op   string  `json:"op"`
		Args []WSArg `json:"args"`
	}

	sendPing := make(chan struct{})
	gotPong := make(chan struct{}, 1)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade error: %v", err)
		}
		defer c.Close()

		_, msg, err := c.ReadMessage()
		if err != nil {
			return
		}
		var req opReq
		if err := json.Unmarshal(msg, &req); err != nil {
			t.Fatalf("unmarshal subscribe op: %v", err)
		}
		if req.Op != "subscribe" || len(req.Args) != 1 {
			t.Fatalf("subscribe op req = %#v", req)
		}

		ev := struct {
			ID     string `json:"id"`
			Event  string `json:"event"`
			Arg    WSArg  `json:"arg"`
			ConnID string `json:"connId"`
		}{
			ID:     req.ID,
			Event:  "subscribe",
			Arg:    req.Args[0],
			ConnID: "x",
		}
		ack, _ := json.Marshal(ev)
		_ = c.WriteMessage(websocket.TextMessage, ack)

		<-sendPing
		_ = c.WriteMessage(websocket.TextMessage, []byte("ping"))

		_, pong, err := c.ReadMessage()
		if err == nil && isWSPongMessage(pong) {
			gotPong <- struct{}{}
		}
	}))
	t.Cleanup(srv.Close)

	wsURL := "ws" + srv.URL[len("http"):]

	var dconn *writeDeadlineConn
	dialer := &websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
		NetDialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			c, err := (&net.Dialer{}).DialContext(ctx, network, addr)
			if err != nil {
				return nil, err
			}
			dc := &writeDeadlineConn{Conn: c}
			dconn = dc
			return dc, nil
		},
	}

	writeTimeout := 1234 * time.Millisecond

	client := NewClient()
	ws := client.NewWSPublic(
		WithWSURL(wsURL),
		WithWSDialer(dialer),
		WithWSHeartbeat(0),
		WithWSWriteTimeout(writeTimeout),
	)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	if err := ws.Start(ctx, nil, nil); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	t.Cleanup(ws.Close)

	subCtx, subCancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(subCancel)

	subStart := time.Now()
	if err := ws.SubscribeAndWait(subCtx, WSArg{Channel: WSChannelTickers, InstId: "BTC-USDT"}); err != nil {
		t.Fatalf("SubscribeAndWait() error = %v", err)
	}
	if dconn == nil {
		t.Fatalf("expected dial conn")
	}

	dl1 := dconn.LastWriteDeadline()
	if dl1.IsZero() {
		t.Fatalf("expected write deadline to be set on subscribe write")
	}
	if got := dl1.Sub(subStart); got < writeTimeout {
		t.Fatalf("subscribe write deadline delta=%s, want >= %s", got, writeTimeout)
	}

	pingStart := time.Now()
	close(sendPing)

	select {
	case <-gotPong:
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting pong")
	}

	dl2 := dconn.LastWriteDeadline()
	if dl2.IsZero() {
		t.Fatalf("expected write deadline to be set on pong write")
	}
	if got := dl2.Sub(pingStart); got < writeTimeout {
		t.Fatalf("pong write deadline delta=%s, want >= %s", got, writeTimeout)
	}
}
