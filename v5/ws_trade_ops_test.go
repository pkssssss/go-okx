package okx

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWSClient_PlaceOrder_RequiresPrivate(t *testing.T) {
	c := NewClient()
	ws := c.NewWSPublic()

	_, err := ws.PlaceOrder(context.Background(), WSPlaceOrderArg{
		InstId:     "BTC-USDT",
		TdMode:     "cash",
		Side:       "buy",
		OrdType:    "market",
		Sz:         "1",
		ClOrdId:    "x",
		InstIdCode: 0,
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, errWSPrivateRequired) {
		t.Fatalf("error = %v, want errWSPrivateRequired", err)
	}
}

func TestWSClient_PlaceOrder_RejectsBusinessPrivate(t *testing.T) {
	c := NewClient()
	ws := c.NewWSBusinessPrivate()

	_, err := ws.PlaceOrder(context.Background(), WSPlaceOrderArg{
		InstId:  "BTC-USDT",
		TdMode:  "cash",
		Side:    "buy",
		OrdType: "market",
		Sz:      "1",
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, errWSPrivateRequired) {
		t.Fatalf("error = %v, want errWSPrivateRequired", err)
	}
}

func TestWSClient_CancelOrder_RejectsBothOrdIDAndClOrdID(t *testing.T) {
	c := NewClient()
	ws := c.NewWSPrivate()

	_, err := ws.CancelOrder(context.Background(), WSCancelOrderArg{
		InstId:  "BTC-USDT",
		OrdId:   "o1",
		ClOrdId: "c1",
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "exactly one of ordId or clOrdId") {
		t.Fatalf("error = %v, want contains %q", err, "exactly one of ordId or clOrdId")
	}
}

func TestWSClient_AmendOrder_RejectsBothOrdIDAndClOrdID(t *testing.T) {
	c := NewClient()
	ws := c.NewWSPrivate()

	_, err := ws.AmendOrder(context.Background(), WSAmendOrderArg{
		InstId:  "BTC-USDT",
		OrdId:   "o1",
		ClOrdId: "c1",
		NewSz:   "2",
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "exactly one of ordId or clOrdId") {
		t.Fatalf("error = %v, want contains %q", err, "exactly one of ordId or clOrdId")
	}
}

func TestWSClient_PlaceOrder_WSOpReply(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	type opReq struct {
		ID   string          `json:"id"`
		Op   string          `json:"op"`
		Args json.RawMessage `json:"args"`
	}

	t.Run("success", func(t *testing.T) {
		opReqCh := make(chan opReq, 1)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if handleTradeAccountRateLimitMock(w, r) {
				return
			}
			c, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				t.Fatalf("upgrade error: %v", err)
			}
			defer c.Close()

			// login
			if _, _, err := c.ReadMessage(); err != nil {
				t.Fatalf("server read login: %v", err)
			}
			_ = c.WriteMessage(websocket.TextMessage, []byte(`{"event":"login","code":"0","msg":"","connId":"x"}`))

			// op=order
			_, msg, err := c.ReadMessage()
			if err != nil {
				t.Fatalf("server read op: %v", err)
			}
			var req opReq
			if err := json.Unmarshal(msg, &req); err != nil {
				t.Fatalf("unmarshal op: %v", err)
			}
			opReqCh <- req

			resp := `{"id":"` + req.ID + `","op":"order","code":"0","msg":"","data":[{"clOrdId":"c1","ordId":"o1","ts":"1700000000000","sCode":"0","sMsg":""}],"inTime":"1","outTime":"2"}`
			_ = c.WriteMessage(websocket.TextMessage, []byte(resp))

			for {
				if _, _, err := c.ReadMessage(); err != nil {
					return
				}
			}
		}))
		t.Cleanup(srv.Close)

		wsURL := "ws" + srv.URL[len("http"):]

		client := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
			WithCredentials(Credentials{
				APIKey:     "mykey",
				SecretKey:  "mysecret",
				Passphrase: "mypass",
			}),
		)
		ws := client.NewWSPrivate(WithWSURL(wsURL))

		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)
		if err := ws.Start(ctx, nil, nil); err != nil {
			t.Fatalf("Start() error = %v", err)
		}
		t.Cleanup(ws.Close)

		opCtx, opCancel := context.WithTimeout(context.Background(), 2*time.Second)
		t.Cleanup(opCancel)

		ack, err := ws.PlaceOrder(opCtx, WSPlaceOrderArg{
			InstId:  "BTC-USDT",
			TdMode:  "cash",
			Side:    "buy",
			OrdType: "market",
			Sz:      "1",
			ClOrdId: "c1",
		})
		if err != nil {
			t.Fatalf("PlaceOrder() error = %v", err)
		}
		if ack == nil || ack.OrdId != "o1" || ack.ClOrdId != "c1" || ack.SCode != "0" || ack.TS != 1700000000000 {
			t.Fatalf("ack = %#v", ack)
		}

		select {
		case req := <-opReqCh:
			if req.ID == "" || req.Op != "order" {
				t.Fatalf("op req = %#v", req)
			}
			var args []WSPlaceOrderArg
			if err := json.Unmarshal(req.Args, &args); err != nil {
				t.Fatalf("unmarshal args: %v", err)
			}
			if len(args) != 1 || args[0].InstId != "BTC-USDT" || args[0].Side != "buy" || args[0].OrdType != "market" || args[0].Sz != "1" {
				t.Fatalf("args = %#v", args)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("timeout waiting op req")
		}
	})

	t.Run("top_level_code_error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if handleTradeAccountRateLimitMock(w, r) {
				return
			}
			c, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				t.Fatalf("upgrade error: %v", err)
			}
			defer c.Close()

			if _, _, err := c.ReadMessage(); err != nil {
				t.Fatalf("server read login: %v", err)
			}
			_ = c.WriteMessage(websocket.TextMessage, []byte(`{"event":"login","code":"0","msg":"","connId":"x"}`))

			_, msg, err := c.ReadMessage()
			if err != nil {
				t.Fatalf("server read op: %v", err)
			}
			var req opReq
			if err := json.Unmarshal(msg, &req); err != nil {
				t.Fatalf("unmarshal op: %v", err)
			}

			resp := `{"id":"` + req.ID + `","op":"order","code":"50000","msg":"bad","data":[],"inTime":"1","outTime":"2"}`
			_ = c.WriteMessage(websocket.TextMessage, []byte(resp))

			for {
				if _, _, err := c.ReadMessage(); err != nil {
					return
				}
			}
		}))
		t.Cleanup(srv.Close)

		wsURL := "ws" + srv.URL[len("http"):]

		client := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
			WithCredentials(Credentials{
				APIKey:     "mykey",
				SecretKey:  "mysecret",
				Passphrase: "mypass",
			}),
		)
		ws := client.NewWSPrivate(WithWSURL(wsURL))

		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)
		if err := ws.Start(ctx, nil, nil); err != nil {
			t.Fatalf("Start() error = %v", err)
		}
		t.Cleanup(ws.Close)

		opCtx, opCancel := context.WithTimeout(context.Background(), 2*time.Second)
		t.Cleanup(opCancel)

		_, err := ws.PlaceOrder(opCtx, WSPlaceOrderArg{
			InstId:  "BTC-USDT",
			TdMode:  "cash",
			Side:    "buy",
			OrdType: "market",
			Sz:      "1",
		})
		if err == nil {
			t.Fatalf("expected error")
		}
		var oe *WSTradeOpError
		if !errors.As(err, &oe) || oe.Code != "50000" || oe.Op != "order" {
			t.Fatalf("error = %#v", err)
		}
	})

	t.Run("data_sCode_error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if handleTradeAccountRateLimitMock(w, r) {
				return
			}
			c, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				t.Fatalf("upgrade error: %v", err)
			}
			defer c.Close()

			if _, _, err := c.ReadMessage(); err != nil {
				t.Fatalf("server read login: %v", err)
			}
			_ = c.WriteMessage(websocket.TextMessage, []byte(`{"event":"login","code":"0","msg":"","connId":"x"}`))

			_, msg, err := c.ReadMessage()
			if err != nil {
				t.Fatalf("server read op: %v", err)
			}
			var req opReq
			if err := json.Unmarshal(msg, &req); err != nil {
				t.Fatalf("unmarshal op: %v", err)
			}

			resp := `{"id":"` + req.ID + `","op":"order","code":"0","msg":"","data":[{"sCode":"51000","sMsg":"invalid"}],"inTime":"1","outTime":"2"}`
			_ = c.WriteMessage(websocket.TextMessage, []byte(resp))

			for {
				if _, _, err := c.ReadMessage(); err != nil {
					return
				}
			}
		}))
		t.Cleanup(srv.Close)

		wsURL := "ws" + srv.URL[len("http"):]

		client := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
			WithCredentials(Credentials{
				APIKey:     "mykey",
				SecretKey:  "mysecret",
				Passphrase: "mypass",
			}),
		)
		ws := client.NewWSPrivate(WithWSURL(wsURL))

		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)
		if err := ws.Start(ctx, nil, nil); err != nil {
			t.Fatalf("Start() error = %v", err)
		}
		t.Cleanup(ws.Close)

		opCtx, opCancel := context.WithTimeout(context.Background(), 2*time.Second)
		t.Cleanup(opCancel)

		_, err := ws.PlaceOrder(opCtx, WSPlaceOrderArg{
			InstId:  "BTC-USDT",
			TdMode:  "cash",
			Side:    "buy",
			OrdType: "market",
			Sz:      "1",
		})
		if err == nil {
			t.Fatalf("expected error")
		}
		var oe *WSTradeOpError
		if !errors.As(err, &oe) || oe.Code != "0" || oe.SCode != "51000" || !strings.Contains(oe.Error(), "sCode=51000") {
			t.Fatalf("error = %#v", err)
		}
		if len(oe.Raw) == 0 || !strings.Contains(string(oe.Raw), `"sCode":"51000"`) {
			t.Fatalf("raw = %q", string(oe.Raw))
		}
	})
}

func TestWSClient_CancelOrder_WSOpReply(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	type opReq struct {
		ID   string          `json:"id"`
		Op   string          `json:"op"`
		Args json.RawMessage `json:"args"`
	}

	opReqCh := make(chan opReq, 1)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handleTradeAccountRateLimitMock(w, r) {
			return
		}
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade error: %v", err)
		}
		defer c.Close()

		if _, _, err := c.ReadMessage(); err != nil {
			t.Fatalf("server read login: %v", err)
		}
		_ = c.WriteMessage(websocket.TextMessage, []byte(`{"event":"login","code":"0","msg":"","connId":"x"}`))

		_, msg, err := c.ReadMessage()
		if err != nil {
			t.Fatalf("server read op: %v", err)
		}
		var req opReq
		if err := json.Unmarshal(msg, &req); err != nil {
			t.Fatalf("unmarshal op: %v", err)
		}
		opReqCh <- req

		resp := `{"id":"` + req.ID + `","op":"cancel-order","code":"0","msg":"","data":[{"ordId":"o1","ts":"1700000000000","sCode":"0","sMsg":""}],"inTime":"1","outTime":"2"}`
		_ = c.WriteMessage(websocket.TextMessage, []byte(resp))

		for {
			if _, _, err := c.ReadMessage(); err != nil {
				return
			}
		}
	}))
	t.Cleanup(srv.Close)

	wsURL := "ws" + srv.URL[len("http"):]

	client := NewClient(
		WithBaseURL(srv.URL),
		WithHTTPClient(srv.Client()),
		WithCredentials(Credentials{
			APIKey:     "mykey",
			SecretKey:  "mysecret",
			Passphrase: "mypass",
		}),
	)
	ws := client.NewWSPrivate(WithWSURL(wsURL))

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	if err := ws.Start(ctx, nil, nil); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	t.Cleanup(ws.Close)

	opCtx, opCancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(opCancel)

	ack, err := ws.CancelOrder(opCtx, WSCancelOrderArg{
		InstId: "BTC-USDT",
		OrdId:  "o1",
	})
	if err != nil {
		t.Fatalf("CancelOrder() error = %v", err)
	}
	if ack == nil || ack.OrdId != "o1" || ack.SCode != "0" || ack.TS != 1700000000000 {
		t.Fatalf("ack = %#v", ack)
	}

	select {
	case req := <-opReqCh:
		if req.ID == "" || req.Op != "cancel-order" {
			t.Fatalf("op req = %#v", req)
		}
		var args []WSCancelOrderArg
		if err := json.Unmarshal(req.Args, &args); err != nil {
			t.Fatalf("unmarshal args: %v", err)
		}
		if len(args) != 1 || args[0].InstId != "BTC-USDT" || args[0].OrdId != "o1" {
			t.Fatalf("args = %#v", args)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting op req")
	}
}

func TestWSClient_AmendOrder_WSOpReply(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	type opReq struct {
		ID   string          `json:"id"`
		Op   string          `json:"op"`
		Args json.RawMessage `json:"args"`
	}

	opReqCh := make(chan opReq, 1)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handleTradeAccountRateLimitMock(w, r) {
			return
		}
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade error: %v", err)
		}
		defer c.Close()

		if _, _, err := c.ReadMessage(); err != nil {
			t.Fatalf("server read login: %v", err)
		}
		_ = c.WriteMessage(websocket.TextMessage, []byte(`{"event":"login","code":"0","msg":"","connId":"x"}`))

		_, msg, err := c.ReadMessage()
		if err != nil {
			t.Fatalf("server read op: %v", err)
		}
		var req opReq
		if err := json.Unmarshal(msg, &req); err != nil {
			t.Fatalf("unmarshal op: %v", err)
		}
		opReqCh <- req

		resp := `{"id":"` + req.ID + `","op":"amend-order","code":"0","msg":"","data":[{"ordId":"o1","reqId":"r1","ts":"1700000000000","sCode":"0","sMsg":""}],"inTime":"1","outTime":"2"}`
		_ = c.WriteMessage(websocket.TextMessage, []byte(resp))

		for {
			if _, _, err := c.ReadMessage(); err != nil {
				return
			}
		}
	}))
	t.Cleanup(srv.Close)

	wsURL := "ws" + srv.URL[len("http"):]

	client := NewClient(
		WithBaseURL(srv.URL),
		WithHTTPClient(srv.Client()),
		WithCredentials(Credentials{
			APIKey:     "mykey",
			SecretKey:  "mysecret",
			Passphrase: "mypass",
		}),
	)
	ws := client.NewWSPrivate(WithWSURL(wsURL))

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	if err := ws.Start(ctx, nil, nil); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	t.Cleanup(ws.Close)

	opCtx, opCancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(opCancel)

	ack, err := ws.AmendOrder(opCtx, WSAmendOrderArg{
		InstId: "BTC-USDT",
		OrdId:  "o1",
		NewSz:  "2",
		ReqId:  "r1",
	})
	if err != nil {
		t.Fatalf("AmendOrder() error = %v", err)
	}
	if ack == nil || ack.OrdId != "o1" || ack.ReqId != "r1" || ack.SCode != "0" || ack.TS != 1700000000000 {
		t.Fatalf("ack = %#v", ack)
	}

	select {
	case req := <-opReqCh:
		if req.ID == "" || req.Op != "amend-order" {
			t.Fatalf("op req = %#v", req)
		}
		var args []WSAmendOrderArg
		if err := json.Unmarshal(req.Args, &args); err != nil {
			t.Fatalf("unmarshal args: %v", err)
		}
		if len(args) != 1 || args[0].InstId != "BTC-USDT" || args[0].OrdId != "o1" || args[0].NewSz != "2" || args[0].ReqId != "r1" {
			t.Fatalf("args = %#v", args)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting op req")
	}
}

func TestWSClient_TradeOp_EventError(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	type opReq struct {
		ID   string          `json:"id"`
		Op   string          `json:"op"`
		Args json.RawMessage `json:"args"`
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handleTradeAccountRateLimitMock(w, r) {
			return
		}
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade error: %v", err)
		}
		defer c.Close()

		if _, _, err := c.ReadMessage(); err != nil {
			t.Fatalf("server read login: %v", err)
		}
		_ = c.WriteMessage(websocket.TextMessage, []byte(`{"event":"login","code":"0","msg":"","connId":"x"}`))

		_, msg, err := c.ReadMessage()
		if err != nil {
			t.Fatalf("server read op: %v", err)
		}
		var req opReq
		if err := json.Unmarshal(msg, &req); err != nil {
			t.Fatalf("unmarshal op: %v", err)
		}

		_ = c.WriteMessage(websocket.TextMessage, []byte(`{"id":"`+req.ID+`","event":"error","code":"60012","msg":"Invalid request","connId":"x"}`))

		for {
			if _, _, err := c.ReadMessage(); err != nil {
				return
			}
		}
	}))
	t.Cleanup(srv.Close)

	wsURL := "ws" + srv.URL[len("http"):]

	client := NewClient(
		WithBaseURL(srv.URL),
		WithHTTPClient(srv.Client()),
		WithCredentials(Credentials{
			APIKey:     "mykey",
			SecretKey:  "mysecret",
			Passphrase: "mypass",
		}),
	)
	ws := client.NewWSPrivate(WithWSURL(wsURL))

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	if err := ws.Start(ctx, nil, nil); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	t.Cleanup(ws.Close)

	opCtx, opCancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(opCancel)

	_, err := ws.PlaceOrder(opCtx, WSPlaceOrderArg{
		InstId:  "BTC-USDT",
		TdMode:  "cash",
		Side:    "buy",
		OrdType: "market",
		Sz:      "1",
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "code=60012") {
		t.Fatalf("error = %v, want contains code=60012", err)
	}
}

func TestWSClient_PlaceOrders_WSOpReply(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	type opReq struct {
		ID   string          `json:"id"`
		Op   string          `json:"op"`
		Args json.RawMessage `json:"args"`
	}

	t.Run("success_multi", func(t *testing.T) {
		opReqCh := make(chan opReq, 1)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if handleTradeAccountRateLimitMock(w, r) {
				return
			}
			c, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				t.Fatalf("upgrade error: %v", err)
			}
			defer c.Close()

			if _, _, err := c.ReadMessage(); err != nil {
				t.Fatalf("server read login: %v", err)
			}
			_ = c.WriteMessage(websocket.TextMessage, []byte(`{"event":"login","code":"0","msg":"","connId":"x"}`))

			_, msg, err := c.ReadMessage()
			if err != nil {
				t.Fatalf("server read op: %v", err)
			}
			var req opReq
			if err := json.Unmarshal(msg, &req); err != nil {
				t.Fatalf("unmarshal op: %v", err)
			}
			opReqCh <- req

			resp := `{"id":"` + req.ID + `","op":"order","code":"0","msg":"","data":[{"clOrdId":"c1","ordId":"o1","ts":"1700000000000","sCode":"0","sMsg":""},{"clOrdId":"c2","ordId":"o2","ts":"1700000000001","sCode":"0","sMsg":""}],"inTime":"1","outTime":"2"}`
			_ = c.WriteMessage(websocket.TextMessage, []byte(resp))

			for {
				if _, _, err := c.ReadMessage(); err != nil {
					return
				}
			}
		}))
		t.Cleanup(srv.Close)

		wsURL := "ws" + srv.URL[len("http"):]

		client := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
			WithCredentials(Credentials{
				APIKey:     "mykey",
				SecretKey:  "mysecret",
				Passphrase: "mypass",
			}),
		)
		ws := client.NewWSPrivate(WithWSURL(wsURL))

		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)
		if err := ws.Start(ctx, nil, nil); err != nil {
			t.Fatalf("Start() error = %v", err)
		}
		t.Cleanup(ws.Close)

		opCtx, opCancel := context.WithTimeout(context.Background(), 2*time.Second)
		t.Cleanup(opCancel)

		acks, err := ws.PlaceOrders(opCtx,
			WSPlaceOrderArg{InstId: "BTC-USDT", TdMode: "cash", Side: "buy", OrdType: "market", Sz: "1", ClOrdId: "c1"},
			WSPlaceOrderArg{InstId: "BTC-USDT", TdMode: "cash", Side: "sell", OrdType: "market", Sz: "2", ClOrdId: "c2"},
		)
		if err != nil {
			t.Fatalf("PlaceOrders() error = %v", err)
		}
		if len(acks) != 2 || acks[0].OrdId != "o1" || acks[1].OrdId != "o2" {
			t.Fatalf("acks = %#v", acks)
		}

		select {
		case req := <-opReqCh:
			if req.ID == "" || req.Op != "order" {
				t.Fatalf("op req = %#v", req)
			}
			var args []WSPlaceOrderArg
			if err := json.Unmarshal(req.Args, &args); err != nil {
				t.Fatalf("unmarshal args: %v", err)
			}
			if len(args) != 2 || args[0].ClOrdId != "c1" || args[1].ClOrdId != "c2" {
				t.Fatalf("args = %#v", args)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("timeout waiting op req")
		}
	})

	t.Run("partial_failure", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if handleTradeAccountRateLimitMock(w, r) {
				return
			}
			c, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				t.Fatalf("upgrade error: %v", err)
			}
			defer c.Close()

			if _, _, err := c.ReadMessage(); err != nil {
				t.Fatalf("server read login: %v", err)
			}
			_ = c.WriteMessage(websocket.TextMessage, []byte(`{"event":"login","code":"0","msg":"","connId":"x"}`))

			_, msg, err := c.ReadMessage()
			if err != nil {
				t.Fatalf("server read op: %v", err)
			}
			var req opReq
			if err := json.Unmarshal(msg, &req); err != nil {
				t.Fatalf("unmarshal op: %v", err)
			}

			resp := `{"id":"` + req.ID + `","op":"order","code":"0","msg":"","data":[{"clOrdId":"c1","ordId":"o1","ts":"1700000000000","sCode":"0","sMsg":""},{"clOrdId":"c2","ordId":"","ts":"1700000000001","sCode":"51000","sMsg":"invalid"}],"inTime":"1","outTime":"2"}`
			_ = c.WriteMessage(websocket.TextMessage, []byte(resp))

			for {
				if _, _, err := c.ReadMessage(); err != nil {
					return
				}
			}
		}))
		t.Cleanup(srv.Close)

		wsURL := "ws" + srv.URL[len("http"):]

		client := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
			WithCredentials(Credentials{
				APIKey:     "mykey",
				SecretKey:  "mysecret",
				Passphrase: "mypass",
			}),
		)
		ws := client.NewWSPrivate(WithWSURL(wsURL))

		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)
		if err := ws.Start(ctx, nil, nil); err != nil {
			t.Fatalf("Start() error = %v", err)
		}
		t.Cleanup(ws.Close)

		opCtx, opCancel := context.WithTimeout(context.Background(), 2*time.Second)
		t.Cleanup(opCancel)

		acks, err := ws.PlaceOrders(opCtx,
			WSPlaceOrderArg{InstId: "BTC-USDT", TdMode: "cash", Side: "buy", OrdType: "market", Sz: "1", ClOrdId: "c1"},
			WSPlaceOrderArg{InstId: "BTC-USDT", TdMode: "cash", Side: "sell", OrdType: "market", Sz: "2", ClOrdId: "c2"},
		)
		if err == nil {
			t.Fatalf("expected error")
		}
		var be *WSTradeOpBatchError
		if !errors.As(err, &be) || be.Op != "order" || be.ID == "" || len(be.Acks) != 2 {
			t.Fatalf("error = %#v", err)
		}
		if len(acks) != 2 || acks[1].SCode != "51000" {
			t.Fatalf("acks = %#v", acks)
		}
		if len(be.Raw) == 0 || !strings.Contains(string(be.Raw), `"sCode":"51000"`) {
			t.Fatalf("raw = %q", string(be.Raw))
		}
	})
}

func TestWSClient_TradeOp_AccountRateLimitPrimeFailure_FailClosed(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	opRecvCh := make(chan []byte, 1)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/api/v5/trade/account-rate-limit" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"code":"1","msg":"fail","data":[]}`))
			return
		}

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade error: %v", err)
		}
		defer c.Close()

		if _, _, err := c.ReadMessage(); err != nil {
			t.Fatalf("server read login: %v", err)
		}
		_ = c.WriteMessage(websocket.TextMessage, []byte(`{"event":"login","code":"0","msg":"","connId":"x"}`))

		_ = c.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
		_, msg, err := c.ReadMessage()
		if err == nil {
			opRecvCh <- msg
		}
	}))
	t.Cleanup(srv.Close)

	wsURL := "ws" + srv.URL[len("http"):]

	client := NewClient(
		WithBaseURL(srv.URL),
		WithHTTPClient(srv.Client()),
		WithCredentials(Credentials{
			APIKey:     "mykey",
			SecretKey:  "mysecret",
			Passphrase: "mypass",
		}),
	)
	ws := client.NewWSPrivate(WithWSURL(wsURL), WithWSHeartbeat(0))

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	if err := ws.Start(ctx, nil, nil); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	t.Cleanup(ws.Close)

	opCtx, opCancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(opCancel)

	_, err := ws.PlaceOrder(opCtx, WSPlaceOrderArg{
		InstId:  "BTC-USDT",
		TdMode:  "cash",
		Side:    "buy",
		OrdType: "market",
		Sz:      "1",
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	var rse *RequestStateError
	if !errors.As(err, &rse) || rse.Stage != RequestStagePreflight || rse.Dispatched {
		t.Fatalf("error = %#v, want *RequestStateError{Stage:preflight, Dispatched:false}", err)
	}

	select {
	case msg := <-opRecvCh:
		t.Fatalf("unexpected op message dispatched: %s", string(msg))
	case <-time.After(250 * time.Millisecond):
	}
}
