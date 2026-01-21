package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestTradeAccountRateLimit_AutoPrimedOnFirstTradeRequest(t *testing.T) {
	var (
		mu    sync.Mutex
		paths []string
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		paths = append(paths, r.URL.Path)
		mu.Unlock()

		switch r.URL.Path {
		case "/api/v5/trade/account-rate-limit":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"accRateLimit":"1","fillRatio":"0","mainFillRatio":"0","nextAccRateLimit":"1","ts":"1"}]}`))
		case "/api/v5/trade/order":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"ordId":"1","sCode":"0","sMsg":""}]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)

	c := NewClient(
		WithBaseURL(srv.URL),
		WithHTTPClient(srv.Client()),
		WithCredentials(Credentials{APIKey: "k", SecretKey: "s", Passphrase: "p"}),
	)

	_, err := c.NewPlaceOrderService().
		InstId("BTC-USDT").
		TdMode("cash").
		Side("buy").
		OrdType("market").
		Sz("1").
		Do(context.Background())
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if got, want := len(paths), 2; got != want {
		t.Fatalf("paths = %v, want %d entries", paths, want)
	}
	if got, want := paths[0], "/api/v5/trade/account-rate-limit"; got != want {
		t.Fatalf("first path = %q, want %q (auto prime)", got, want)
	}
	if got, want := paths[1], "/api/v5/trade/order"; got != want {
		t.Fatalf("second path = %q, want %q", got, want)
	}
}
