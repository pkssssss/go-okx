package rest

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestBuildRequestPath(t *testing.T) {
	t.Run("nil_query", func(t *testing.T) {
		if got, want := BuildRequestPath("/api/v5/public/time", nil), "/api/v5/public/time"; got != want {
			t.Fatalf("BuildRequestPath() = %q, want %q", got, want)
		}
	})

	t.Run("with_query", func(t *testing.T) {
		q := url.Values{}
		q.Set("ccy", "BTC")
		q.Set("instId", "BTC-USDT")
		if got, want := BuildRequestPath("/api/v5/account/balance", q), "/api/v5/account/balance?ccy=BTC&instId=BTC-USDT"; got != want {
			t.Fatalf("BuildRequestPath() = %q, want %q", got, want)
		}
	})
}

func TestClientDo_ResponseBodyTooLarge(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("0123456789ABCDEF"))
	}))
	t.Cleanup(srv.Close)

	c := &Client{
		BaseURL:              srv.URL,
		HTTPClient:           srv.Client(),
		DefaultTimeout:       -1,
		MaxResponseBodyBytes: 10,
	}

	_, _, _, err := c.Do(context.Background(), http.MethodGet, "/", nil, nil)
	if err == nil {
		t.Fatalf("expected error")
	}
	var tooLarge *ResponseBodyTooLargeError
	if !errors.As(err, &tooLarge) {
		t.Fatalf("error = %T, want *ResponseBodyTooLargeError", err)
	}
	if tooLarge.MaxBytes != 10 {
		t.Fatalf("MaxBytes = %d, want %d", tooLarge.MaxBytes, 10)
	}
}
