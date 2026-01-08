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

func TestWSClient_UnsubscribeAndWait(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	type opMsg struct {
		ID   string  `json:"id"`
		Op   string  `json:"op"`
		Args []WSArg `json:"args"`
	}

	t.Run("success_removes_desired", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				t.Fatalf("upgrade error: %v", err)
			}
			defer c.Close()

			for {
				_, msg, err := c.ReadMessage()
				if err != nil {
					return
				}

				var om opMsg
				if err := json.Unmarshal(msg, &om); err != nil {
					t.Fatalf("unmarshal op: %v", err)
				}
				if om.ID == "" || (om.Op != "subscribe" && om.Op != "unsubscribe") {
					t.Fatalf("op msg = %#v", om)
				}

				for _, a := range om.Args {
					ev := WSEvent{
						ID:     om.ID,
						Event:  om.Op,
						Arg:    &a,
						ConnID: "x",
					}
					b, _ := json.Marshal(ev)
					_ = c.WriteMessage(websocket.TextMessage, b)
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

		arg := WSArg{Channel: "tickers", InstId: "BTC-USDT"}
		if err := ws.SubscribeAndWait(waitCtx, arg); err != nil {
			t.Fatalf("SubscribeAndWait() error = %v", err)
		}
		if got := len(ws.snapshotDesired()); got != 1 {
			t.Fatalf("desired = %d, want %d", got, 1)
		}

		waitCtx2, waitCancel2 := context.WithTimeout(context.Background(), 2*time.Second)
		t.Cleanup(waitCancel2)
		if err := ws.UnsubscribeAndWait(waitCtx2, arg); err != nil {
			t.Fatalf("UnsubscribeAndWait() error = %v", err)
		}
		if got := len(ws.snapshotDesired()); got != 0 {
			t.Fatalf("desired = %d, want %d", got, 0)
		}
	})

	t.Run("error_event_keeps_desired", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				t.Fatalf("upgrade error: %v", err)
			}
			defer c.Close()

			for {
				_, msg, err := c.ReadMessage()
				if err != nil {
					return
				}

				var om opMsg
				if err := json.Unmarshal(msg, &om); err != nil {
					t.Fatalf("unmarshal op: %v", err)
				}
				if om.Op == "unsubscribe" {
					_ = c.WriteMessage(websocket.TextMessage, []byte(`{"id":"`+om.ID+`","event":"error","code":"60012","msg":"Invalid request","connId":"x"}`))
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

		ws.Subscribe(WSArg{Channel: "tickers", InstId: "BTC-USDT"})
		if got := len(ws.snapshotDesired()); got != 1 {
			t.Fatalf("desired = %d, want %d", got, 1)
		}

		waitCtx, waitCancel := context.WithTimeout(context.Background(), 2*time.Second)
		t.Cleanup(waitCancel)

		err := ws.UnsubscribeAndWait(waitCtx, WSArg{Channel: "tickers", InstId: "BTC-USDT"})
		if err == nil {
			t.Fatalf("expected error")
		}
		if !strings.Contains(err.Error(), "code=60012") {
			t.Fatalf("error = %v, want contains code=60012", err)
		}
		if got := len(ws.snapshotDesired()); got != 1 {
			t.Fatalf("desired = %d, want %d", got, 1)
		}
	})
}
