package okx

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWSClient_SprdPlaceOrder_RequiresBusinessPrivate(t *testing.T) {
	c := NewClient(WithCredentials(Credentials{APIKey: "k", SecretKey: "s", Passphrase: "p"}))

	t.Run("public_business_requires_private", func(t *testing.T) {
		ws := c.NewWSBusiness()
		_, err := ws.SprdPlaceOrder(context.Background(), WSSprdPlaceOrderArg{
			SprdId:  "BTC-USDT_BTC-USDT-SWAP",
			Side:    "buy",
			OrdType: "market",
			Sz:      "1",
		})
		if !errors.Is(err, errWSPrivateRequired) {
			t.Fatalf("error = %v, want errWSPrivateRequired", err)
		}
	})

	t.Run("private_non_business_requires_business", func(t *testing.T) {
		ws := c.NewWSPrivate()
		_, err := ws.SprdPlaceOrder(context.Background(), WSSprdPlaceOrderArg{
			SprdId:  "BTC-USDT_BTC-USDT-SWAP",
			Side:    "buy",
			OrdType: "market",
			Sz:      "1",
		})
		if !errors.Is(err, errWSBusinessRequired) {
			t.Fatalf("error = %v, want errWSBusinessRequired", err)
		}
	})
}

func TestWSClient_SprdPlaceOrder_WSOpReply(t *testing.T) {
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

	type opReq struct {
		ID   string          `json:"id"`
		Op   string          `json:"op"`
		Args json.RawMessage `json:"args"`
	}

	opReqCh := make(chan opReq, 1)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		// op=sprd-order
		_, msg, err := c.ReadMessage()
		if err != nil {
			t.Fatalf("server read op: %v", err)
		}
		var req opReq
		if err := json.Unmarshal(msg, &req); err != nil {
			t.Fatalf("unmarshal op: %v", err)
		}
		opReqCh <- req

		resp := `{"id":"` + req.ID + `","op":"sprd-order","code":"0","msg":"","data":[{"clOrdId":"c1","ordId":"o1","sCode":"0","sMsg":"","ts":"1700000000000"}]}`
		_ = c.WriteMessage(websocket.TextMessage, []byte(resp))

		for {
			if _, _, err := c.ReadMessage(); err != nil {
				return
			}
		}
	}))
	t.Cleanup(srv.Close)

	wsURL := "ws" + srv.URL[len("http"):] + "/ws/v5/business"

	client := NewClient(WithCredentials(Credentials{
		APIKey:     "mykey",
		SecretKey:  "mysecret",
		Passphrase: "mypass",
	}))
	ws := client.NewWSBusinessPrivate(WithWSURL(wsURL))

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	if err := ws.Start(ctx, nil, nil); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	t.Cleanup(ws.Close)

	opCtx, opCancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(opCancel)

	ack, err := ws.SprdPlaceOrder(opCtx, WSSprdPlaceOrderArg{
		SprdId:  "BTC-USDT_BTC-USDT-SWAP",
		Side:    "buy",
		OrdType: "market",
		Sz:      "1",
		ClOrdId: "c1",
	})
	if err != nil {
		t.Fatalf("SprdPlaceOrder() error = %v", err)
	}
	if ack == nil || ack.OrdId != "o1" || ack.ClOrdId != "c1" || ack.SCode != "0" || ack.TS != 1700000000000 {
		t.Fatalf("ack = %#v", ack)
	}

	select {
	case req := <-opReqCh:
		if req.ID == "" || req.Op != "sprd-order" {
			t.Fatalf("op req = %#v", req)
		}
		var args []WSSprdPlaceOrderArg
		if err := json.Unmarshal(req.Args, &args); err != nil {
			t.Fatalf("unmarshal args: %v", err)
		}
		if len(args) != 1 || args[0].SprdId != "BTC-USDT_BTC-USDT-SWAP" || args[0].Side != "buy" || args[0].OrdType != "market" || args[0].Sz != "1" {
			t.Fatalf("args = %#v", args)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting op req")
	}
}

func TestWSClient_SprdMassCancel_WSOpReply(t *testing.T) {
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

	type opReq struct {
		ID   string          `json:"id"`
		Op   string          `json:"op"`
		Args json.RawMessage `json:"args"`
	}

	opReqCh := make(chan opReq, 1)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		resp := `{"id":"` + req.ID + `","op":"sprd-mass-cancel","code":"0","msg":"","data":[{"result":true}]}`
		_ = c.WriteMessage(websocket.TextMessage, []byte(resp))

		for {
			if _, _, err := c.ReadMessage(); err != nil {
				return
			}
		}
	}))
	t.Cleanup(srv.Close)

	wsURL := "ws" + srv.URL[len("http"):] + "/ws/v5/business"

	client := NewClient(WithCredentials(Credentials{
		APIKey:     "mykey",
		SecretKey:  "mysecret",
		Passphrase: "mypass",
	}))
	ws := client.NewWSBusinessPrivate(WithWSURL(wsURL))

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	if err := ws.Start(ctx, nil, nil); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	t.Cleanup(ws.Close)

	opCtx, opCancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(opCancel)

	ack, err := ws.SprdMassCancel(opCtx, WSSprdMassCancelArg{SprdId: "BTC-USDT_BTC-USDT-SWAP"})
	if err != nil {
		t.Fatalf("SprdMassCancel() error = %v", err)
	}
	if ack == nil || !ack.Result {
		t.Fatalf("ack = %#v", ack)
	}

	select {
	case req := <-opReqCh:
		if req.ID == "" || req.Op != "sprd-mass-cancel" {
			t.Fatalf("op req = %#v", req)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting op req")
	}
}
