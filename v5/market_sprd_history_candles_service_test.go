package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMarketSprdHistoryCandlesService_Do(t *testing.T) {
	t.Run("missing_sprdId", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewMarketSprdHistoryCandlesService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMarketSprdHistoryCandlesMissingSprdId {
			t.Fatalf("error = %v, want %v", err, errMarketSprdHistoryCandlesMissingSprdId)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/market/sprd-history-candles"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			q := r.URL.Query()
			if got, want := q.Get("sprdId"), "BTC-USDT_BTC-USDT-SWAP"; got != want {
				t.Fatalf("sprdId = %q, want %q", got, want)
			}
			if got, want := q.Get("bar"), "1m"; got != want {
				t.Fatalf("bar = %q, want %q", got, want)
			}
			if got, want := q.Get("limit"), "1"; got != want {
				t.Fatalf("limit = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[["1597026383085","3.721","3.743","3.677","3.708","8422410","1"]]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewMarketSprdHistoryCandlesService().SprdId("BTC-USDT_BTC-USDT-SWAP").Bar("1m").Limit(1).Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].TS != 1597026383085 || got[0].Confirm != "1" {
			t.Fatalf("candles = %#v", got)
		}
	})
}
