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

func TestWSClient_SubscribeAndWait(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	type subMsg struct {
		ID   string  `json:"id"`
		Op   string  `json:"op"`
		Args []WSArg `json:"args"`
	}

	t.Run("success_multi_args", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				t.Fatalf("upgrade error: %v", err)
			}
			defer c.Close()

			_, msg, err := c.ReadMessage()
			if err != nil {
				t.Fatalf("server read subscribe: %v", err)
			}
			var sm subMsg
			if err := json.Unmarshal(msg, &sm); err != nil {
				t.Fatalf("unmarshal subscribe: %v", err)
			}
			if sm.ID == "" || sm.Op != "subscribe" || len(sm.Args) != 2 {
				t.Fatalf("subscribe msg = %#v", sm)
			}

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

			for {
				if _, _, err := c.ReadMessage(); err != nil {
					return
				}
			}
		}))
		t.Cleanup(srv.Close)

		wsURL := "ws" + srv.URL[len("http"):]
		c := NewClient()
		ws := c.NewWSPublic(WithWSURL(wsURL))

		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)
		if err := ws.Start(ctx, nil, nil); err != nil {
			t.Fatalf("Start() error = %v", err)
		}
		t.Cleanup(ws.Close)

		waitCtx, waitCancel := context.WithTimeout(context.Background(), 2*time.Second)
		t.Cleanup(waitCancel)

		if err := ws.SubscribeAndWait(waitCtx,
			WSArg{Channel: "tickers", InstId: "BTC-USDT"},
			WSArg{Channel: "tickers", InstId: "ETH-USDT"},
		); err != nil {
			t.Fatalf("SubscribeAndWait() error = %v", err)
		}
	})

	t.Run("error_event", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				t.Fatalf("upgrade error: %v", err)
			}
			defer c.Close()

			_, msg, err := c.ReadMessage()
			if err != nil {
				t.Fatalf("server read subscribe: %v", err)
			}
			var sm subMsg
			if err := json.Unmarshal(msg, &sm); err != nil {
				t.Fatalf("unmarshal subscribe: %v", err)
			}
			if sm.ID == "" || sm.Op != "subscribe" {
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
		c := NewClient()
		ws := c.NewWSPublic(WithWSURL(wsURL))

		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)
		if err := ws.Start(ctx, nil, nil); err != nil {
			t.Fatalf("Start() error = %v", err)
		}
		t.Cleanup(ws.Close)

		waitCtx, waitCancel := context.WithTimeout(context.Background(), 2*time.Second)
		t.Cleanup(waitCancel)

		err := ws.SubscribeAndWait(waitCtx, WSArg{Channel: "tickers", InstId: "BTC-USDT"})
		if err == nil {
			t.Fatalf("expected error")
		}
		if !strings.Contains(err.Error(), "code=60012") {
			t.Fatalf("error = %v, want contains code=60012", err)
		}
	})
}
