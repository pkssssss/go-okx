package okx

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWSClient_PrivateStart_MissingCredentials_FailFastNoDial(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	var dials atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dials.Add(1)
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		_ = conn.Close()
	}))
	t.Cleanup(srv.Close)

	wsURL := "ws" + srv.URL[len("http"):]
	client := NewClient()

	tests := []struct {
		name string
		new  func(*Client) *WSClient
	}{
		{
			name: "private",
			new: func(c *Client) *WSClient {
				return c.NewWSPrivate(WithWSURL(wsURL))
			},
		},
		{
			name: "business_private",
			new: func(c *Client) *WSClient {
				return c.NewWSBusinessPrivate(WithWSURL(wsURL))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ws := tt.new(client)
			err := ws.Start(context.Background(), nil, nil)
			if !errors.Is(err, errMissingCredentials) {
				t.Fatalf("Start() error = %v, want %v", err, errMissingCredentials)
			}
			if ws.started.Load() {
				t.Fatalf("started=true after Start() error")
			}
			select {
			case <-ws.Done():
			default:
				t.Fatalf("Done() should be closed when Start() failed")
			}
		})
	}

	time.Sleep(200 * time.Millisecond)
	if got := dials.Load(); got != 0 {
		t.Fatalf("dials = %d, want 0", got)
	}
}
