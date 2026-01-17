package okx

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWSClient_AutoResubscribe_ErrorEventVisible(t *testing.T) {
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

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

		var sm struct {
			ID   string  `json:"id"`
			Op   string  `json:"op"`
			Args []WSArg `json:"args"`
		}
		if err := json.Unmarshal(msg, &sm); err != nil {
			t.Fatalf("unmarshal subscribe: %v", err)
		}
		if sm.Op != "subscribe" || sm.ID == "" || len(sm.Args) != 1 {
			t.Fatalf("subscribe msg = %#v", sm)
		}

		_ = c.WriteMessage(websocket.TextMessage, []byte(`{"id":"`+sm.ID+`","event":"error","code":"60012","msg":"Invalid request","connId":"x"}`))

		for {
			if _, _, err := c.ReadMessage(); err != nil {
				return
			}
		}
	}))
	t.Cleanup(srv.Close)

	wsURL := "ws" + srv.URL[len("http"):]

	errCh := make(chan error, 1)

	client := NewClient()
	ws := client.NewWSPublic(WithWSURL(wsURL), WithWSHeartbeat(0), WithWSResubscribeWaitTimeout(2*time.Second))
	_ = ws.Subscribe(WSArg{Channel: WSChannelTickers, InstId: "BTC-USDT"})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	if err := ws.Start(ctx, nil, func(err error) {
		select {
		case errCh <- err:
		default:
		}
	}); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	t.Cleanup(ws.Close)

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatalf("expected error")
		}
		if !strings.Contains(err.Error(), "op=subscribe") || !strings.Contains(err.Error(), "code=60012") {
			t.Fatalf("error = %v, want contains op=subscribe and code=60012", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting errHandler called")
	}
}
