package okx

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWSClient_PrivateLogin_Subscribe_PingPong(t *testing.T) {
	fixedNow := time.Unix(1538054050, 0).UTC()

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	type loginMsg struct {
		Op   string `json:"op"`
		Args []struct {
			APIKey     string `json:"apiKey"`
			Passphrase string `json:"passphrase"`
			Timestamp  string `json:"timestamp"`
			Sign       string `json:"sign"`
		} `json:"args"`
	}

	type subMsg struct {
		Op   string  `json:"op"`
		Args []WSArg `json:"args"`
	}

	pongCh := make(chan string, 1)
	loginCh := make(chan loginMsg, 1)
	subCh := make(chan subMsg, 1)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade error: %v", err)
		}
		defer c.Close()

		c.SetPongHandler(func(appData string) error {
			pongCh <- appData
			return nil
		})

		// 触发客户端 pong 回包（需要复制 payload）
		if err := c.WriteControl(websocket.PingMessage, []byte("11446744073709551615"), time.Now().Add(2*time.Second)); err != nil {
			t.Fatalf("server write ping: %v", err)
		}

		// 读 login
		_, msg, err := c.ReadMessage()
		if err != nil {
			t.Fatalf("server read login: %v", err)
		}
		var lm loginMsg
		if err := json.Unmarshal(msg, &lm); err != nil {
			t.Fatalf("unmarshal login: %v", err)
		}
		loginCh <- lm

		// 回登录成功
		_ = c.WriteMessage(websocket.TextMessage, []byte(`{"event":"login","code":"0","msg":"","connId":"x"}`))

		// 读 subscribe
		_, msg, err = c.ReadMessage()
		if err != nil {
			t.Fatalf("server read subscribe: %v", err)
		}
		var sm subMsg
		if err := json.Unmarshal(msg, &sm); err != nil {
			t.Fatalf("unmarshal subscribe: %v", err)
		}
		subCh <- sm
	}))
	t.Cleanup(srv.Close)

	// 将 http://127... 转成 ws://127...
	wsURL := "ws" + srv.URL[len("http"):]

	c := NewClient(
		WithCredentials(Credentials{
			APIKey:     "mykey",
			SecretKey:  "mysecret",
			Passphrase: "mypass",
		}),
		WithNowFunc(func() time.Time { return fixedNow }),
	)

	ws := c.NewWSPrivate(WithWSURL(wsURL))
	_ = ws.Subscribe(WSArg{Channel: "orders", InstType: "SWAP"})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	if err := ws.Start(ctx, nil, nil); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	t.Cleanup(ws.Close)

	select {
	case payload := <-pongCh:
		if payload != "11446744073709551615" {
			t.Fatalf("pong payload = %q, want %q", payload, "11446744073709551615")
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting pong")
	}

	select {
	case lm := <-loginCh:
		if lm.Op != "login" || len(lm.Args) != 1 {
			t.Fatalf("login msg = %#v", lm)
		}
		if got, want := lm.Args[0].APIKey, "mykey"; got != want {
			t.Fatalf("apiKey = %q, want %q", got, want)
		}
		if got, want := lm.Args[0].Passphrase, "mypass"; got != want {
			t.Fatalf("passphrase = %q, want %q", got, want)
		}
		if got, want := lm.Args[0].Timestamp, "1538054050"; got != want {
			t.Fatalf("timestamp = %q, want %q", got, want)
		}
		if got, want := lm.Args[0].Sign, "m+lzVL6siKIpimAa/6y8lHpWZe0SCpehAqymC8Nel0A="; got != want {
			t.Fatalf("sign = %q, want %q", got, want)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting login")
	}

	select {
	case sm := <-subCh:
		if sm.Op != "subscribe" || len(sm.Args) != 1 {
			t.Fatalf("subscribe msg = %#v", sm)
		}
		if got, want := sm.Args[0].Channel, "orders"; got != want {
			t.Fatalf("channel = %q, want %q", got, want)
		}
		if got, want := sm.Args[0].InstType, "SWAP"; got != want {
			t.Fatalf("instType = %q, want %q", got, want)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting subscribe")
	}
}
