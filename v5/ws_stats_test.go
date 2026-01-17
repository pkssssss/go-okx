package okx

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWSClient_Stats_Snapshot(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	type subMsg struct {
		ID   string  `json:"id"`
		Op   string  `json:"op"`
		Args []WSArg `json:"args"`
	}

	subCh := make(chan subMsg, 2)
	var connCount atomic.Int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade error: %v", err)
		}
		defer c.Close()

		thisConn := connCount.Add(1)

		_, msg, err := c.ReadMessage()
		if err != nil {
			return
		}
		var sm subMsg
		if err := json.Unmarshal(msg, &sm); err != nil {
			t.Fatalf("unmarshal subscribe: %v", err)
		}
		subCh <- sm

		for _, a := range sm.Args {
			ev := WSEvent{
				ID:     sm.ID,
				Event:  "subscribe",
				Arg:    &a,
				ConnID: "x",
			}
			b, _ := json.Marshal(ev)
			_ = c.WriteMessage(websocket.TextMessage, b)
		}

		if thisConn == 1 {
			_ = c.WriteMessage(websocket.TextMessage, []byte(`{"event":"notice","code":"64008","msg":"upgrade","connId":"x"}`))
		}

		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				return
			}
		}
	}))
	t.Cleanup(srv.Close)

	wsURL := "ws" + srv.URL[len("http"):]

	c := NewClient()
	ws := c.NewWSPublic(WithWSURL(wsURL))
	_ = ws.Subscribe(WSArg{Channel: WSChannelTickers, InstId: "BTC-USDT"})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	if err := ws.Start(ctx, nil, nil); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	for i := 0; i < 2; i++ {
		select {
		case sm := <-subCh:
			if sm.Op != "subscribe" || len(sm.Args) != 1 || sm.Args[0].Channel != WSChannelTickers {
				t.Fatalf("subscribe msg = %#v", sm)
			}
		case <-time.After(3 * time.Second):
			t.Fatalf("timeout waiting subscribe #%d", i+1)
		}
	}

	var st WSStats
	deadline := time.Now().Add(2 * time.Second)
	for {
		st = ws.Stats()
		if st.Connects >= 2 && st.Reconnects >= 1 && st.SubscribeOK >= 2 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("stats not updated in time: %#v", st)
		}
		time.Sleep(10 * time.Millisecond)
	}

	if !st.Started {
		t.Fatalf("Started = false, want true")
	}
	if !st.Connected {
		t.Fatalf("Connected = false, want true")
	}
	if st.Reconnects < 1 || st.Connects < 2 {
		t.Fatalf("Reconnects/Connects = %d/%d, want >=1 and >=2", st.Reconnects, st.Connects)
	}
	if st.SubscribeOK < 2 {
		t.Fatalf("SubscribeOK = %d, want >=2", st.SubscribeOK)
	}
	if st.DesiredSubscriptions != 1 {
		t.Fatalf("DesiredSubscriptions = %d, want %d", st.DesiredSubscriptions, 1)
	}
	if st.LastRecv.IsZero() {
		t.Fatalf("LastRecv is zero, want non-zero")
	}
	if st.LastError.Message == "" || !strings.Contains(st.LastError.Message, "notice 64008 reconnect") {
		t.Fatalf("LastError = %#v, want contains notice 64008 reconnect", st.LastError)
	}

	ws.Close()
	select {
	case <-ws.Done():
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting Done")
	}
}
