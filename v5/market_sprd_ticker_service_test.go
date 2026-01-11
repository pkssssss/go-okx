package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMarketSprdTickerService_Do(t *testing.T) {
	t.Run("missing_sprdId", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewMarketSprdTickerService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMarketSprdTickerMissingSprdId {
			t.Fatalf("error = %v, want %v", err, errMarketSprdTickerMissingSprdId)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/market/sprd-ticker"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			q := r.URL.Query()
			if got, want := q.Get("sprdId"), "BTC-USDT_BTC-USDT-SWAP"; got != want {
				t.Fatalf("sprdId = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"sprdId":"BTC-USDT_BTC-USDT-SWAP","last":"14.5","lastSz":"0.5","askPx":"8.5","askSz":"12.0","bidPx":"0.5","bidSz":"12.0","open24h":"4","high24h":"14.5","low24h":"-2.2","vol24h":"6.67","ts":"1715331406485"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewMarketSprdTickerService().SprdId("BTC-USDT_BTC-USDT-SWAP").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.SprdId != "BTC-USDT_BTC-USDT-SWAP" || got.Last != "14.5" || got.AskPx != "8.5" || got.BidPx != "0.5" {
			t.Fatalf("ticker = %#v", got)
		}
		if got.TS != 1715331406485 {
			t.Fatalf("TS = %d, want %d", got.TS, 1715331406485)
		}
	})
}
