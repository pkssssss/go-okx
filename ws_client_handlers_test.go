package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWSClient_OnOrders_RoutesData(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade error: %v", err)
		}
		defer c.Close()

		_ = c.WriteMessage(websocket.TextMessage, []byte(`{"arg":{"channel":"orders"},"data":[{"instType":"SWAP","instId":"BTC-USDT-SWAP","ordId":"o1","clOrdId":"c1","side":"buy","ordType":"market","state":"live","px":"0","sz":"1","avgPx":"0","fillPx":"0","fillSz":"0","accFillSz":"0","uTime":"1700000000000","cTime":"1700000000000"}]}`))

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

	gotCh := make(chan TradeOrder, 1)
	ws.OnOrders(func(order TradeOrder) {
		select {
		case gotCh <- order:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	if err := ws.Start(ctx, nil, nil); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	t.Cleanup(ws.Close)

	select {
	case o := <-gotCh:
		if o.OrdId != "o1" || o.InstId != "BTC-USDT-SWAP" || o.Side != "buy" {
			t.Fatalf("order = %#v", o)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting order")
	}
}

func TestWSClient_OnFills_RoutesData(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade error: %v", err)
		}
		defer c.Close()

		_ = c.WriteMessage(websocket.TextMessage, []byte(`{"arg":{"channel":"fills"},"data":[{"instId":"BTC-USDT-SWAP","fillSz":"1","fillPx":"100","side":"buy","ts":"1700000000000","ordId":"o1","clOrdId":"c1","tradeId":"t1","execType":"T","count":"1"}]}`))

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

	gotCh := make(chan WSFill, 1)
	ws.OnFills(func(fill WSFill) {
		select {
		case gotCh <- fill:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	if err := ws.Start(ctx, nil, nil); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	t.Cleanup(ws.Close)

	select {
	case f := <-gotCh:
		if f.OrdId != "o1" || f.TradeId != "t1" || f.InstId != "BTC-USDT-SWAP" {
			t.Fatalf("fill = %#v", f)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting fill")
	}
}

func TestWSClient_OnOpReply_RoutesReply(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade error: %v", err)
		}
		defer c.Close()

		_ = c.WriteMessage(websocket.TextMessage, []byte(`{"id":"99","op":"order","code":"0","msg":"","data":[],"inTime":"1","outTime":"2"}`))

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

	gotCh := make(chan WSOpReply, 1)
	ws.OnOpReply(func(reply WSOpReply, raw []byte) {
		_ = raw
		select {
		case gotCh <- reply:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	if err := ws.Start(ctx, nil, nil); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	t.Cleanup(ws.Close)

	select {
	case r := <-gotCh:
		if r.ID != "99" || r.Op != "order" || r.Code != "0" {
			t.Fatalf("reply = %#v", r)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting op reply")
	}
}
