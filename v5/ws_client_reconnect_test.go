package okx

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWSClient_Close_ClosesConnAndDone(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	type subMsg struct {
		Op   string  `json:"op"`
		Args []WSArg `json:"args"`
	}

	subCh := make(chan subMsg, 1)

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
		var sm subMsg
		if err := json.Unmarshal(msg, &sm); err != nil {
			t.Fatalf("unmarshal subscribe: %v", err)
		}
		subCh <- sm

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

	select {
	case sm := <-subCh:
		if sm.Op != "subscribe" || len(sm.Args) != 1 || sm.Args[0].Channel != WSChannelTickers {
			t.Fatalf("subscribe msg = %#v", sm)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting subscribe")
	}

	ws.Close()

	select {
	case <-ws.Done():
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting Done after Close")
	}
}

func TestWSClient_Notice64008_ReconnectAndResubscribe(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	type subMsg struct {
		Op   string  `json:"op"`
		Args []WSArg `json:"args"`
	}

	var connNum atomic.Int32
	subCh := make(chan subMsg, 2)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade error: %v", err)
		}
		defer c.Close()

		n := connNum.Add(1)

		_, msg, err := c.ReadMessage()
		if err != nil {
			return
		}
		var sm subMsg
		if err := json.Unmarshal(msg, &sm); err != nil {
			t.Fatalf("unmarshal subscribe: %v", err)
		}
		subCh <- sm

		if n == 1 {
			_ = c.WriteMessage(websocket.TextMessage, []byte(`{"event":"notice","code":"64008","msg":"upgrade","connId":"x"}`))
			for {
				_, _, err := c.ReadMessage()
				if err != nil {
					return
				}
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

	ws.Close()
	select {
	case <-ws.Done():
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting Done")
	}
}
