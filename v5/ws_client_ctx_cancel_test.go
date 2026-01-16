package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWSClient_CancelCtx_ClosesConnAndDone(t *testing.T) {
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

	connected := make(chan struct{}, 1)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade error: %v", err)
		}
		defer c.Close()

		// 读取 subscribe（确保客户端已建立连接并进入 readLoop）
		_, _, err = c.ReadMessage()
		if err != nil {
			return
		}
		select {
		case connected <- struct{}{}:
		default:
		}

		// 保持静默：不发送任何消息，模拟低频/无推送场景
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				return
			}
		}
	}))
	t.Cleanup(srv.Close)

	wsURL := "ws" + srv.URL[len("http"):]

	client := NewClient()
	ws := client.NewWSPublic(WithWSURL(wsURL), WithWSHeartbeat(0))
	_ = ws.Subscribe(WSArg{Channel: WSChannelTickers, InstId: "BTC-USDT"})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	if err := ws.Start(ctx, nil, nil); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	t.Cleanup(ws.Close)

	select {
	case <-connected:
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting server connected")
	}

	cancel()

	select {
	case <-ws.Done():
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting Done after ctx cancel")
	}
}
