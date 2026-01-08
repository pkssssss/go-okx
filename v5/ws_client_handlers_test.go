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

func TestWSClient_OnAccount_RoutesData(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade error: %v", err)
		}
		defer c.Close()

		_ = c.WriteMessage(websocket.TextMessage, []byte(`{"arg":{"channel":"account"},"data":[{"uTime":"1700000000000","totalEq":"1","adjEq":"1","availEq":"1","details":[]}]}`))

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

	gotCh := make(chan AccountBalance, 1)
	ws.OnAccount(func(balance AccountBalance) {
		select {
		case gotCh <- balance:
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
	case b := <-gotCh:
		if b.TotalEq != "1" || b.AvailEq != "1" || b.UTime != 1700000000000 {
			t.Fatalf("balance = %#v", b)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting account")
	}
}

func TestWSClient_OnPositions_RoutesData(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade error: %v", err)
		}
		defer c.Close()

		_ = c.WriteMessage(websocket.TextMessage, []byte(`{"arg":{"channel":"positions"},"data":[{"instType":"SWAP","instId":"BTC-USDT-SWAP","posSide":"long","pos":"1","availPos":"1","avgPx":"100","markPx":"100","lever":"1","mgnMode":"cross","ccy":"USDT","uTime":"1700000000000"}]}`))

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

	gotCh := make(chan AccountPosition, 1)
	ws.OnPositions(func(position AccountPosition) {
		select {
		case gotCh <- position:
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
	case p := <-gotCh:
		if p.InstId != "BTC-USDT-SWAP" || p.Pos != "1" || p.UTime != 1700000000000 {
			t.Fatalf("position = %#v", p)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting position")
	}
}

func TestWSClient_OnBalanceAndPosition_RoutesData(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade error: %v", err)
		}
		defer c.Close()

		_ = c.WriteMessage(websocket.TextMessage, []byte(`{"arg":{"channel":"balance_and_position"},"data":[{"pTime":"1700000000000","eventType":"snapshot","balData":[{"ccy":"USDT","cashBal":"1","uTime":"1700000000000"}],"posData":[{"instType":"SWAP","instId":"BTC-USDT-SWAP","posSide":"long","pos":"1","availPos":"1","avgPx":"100","markPx":"100","lever":"1","mgnMode":"cross","ccy":"USDT","uTime":"1700000000000"}]}]}`))

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

	gotCh := make(chan WSBalanceAndPosition, 1)
	ws.OnBalanceAndPosition(func(data WSBalanceAndPosition) {
		select {
		case gotCh <- data:
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
	case d := <-gotCh:
		if d.PTime != 1700000000000 || d.EventType != "snapshot" || len(d.BalData) != 1 || len(d.PosData) != 1 {
			t.Fatalf("data = %#v", d)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting balance_and_position")
	}
}
