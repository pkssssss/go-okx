package rest

import (
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
