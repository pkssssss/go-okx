package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMarketSprdCandlesService_Do(t *testing.T) {
	t.Run("missing_sprdId", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewMarketSprdCandlesService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMarketSprdCandlesMissingSprdId {
			t.Fatalf("error = %v, want %v", err, errMarketSprdCandlesMissingSprdId)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/market/sprd-candles"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			q := r.URL.Query()
			if got, want := q.Get("sprdId"), "BTC-USDT_BTC-USDT-SWAP"; got != want {
				t.Fatalf("sprdId = %q, want %q", got, want)
			}
			if got, want := q.Get("bar"), "1m"; got != want {
				t.Fatalf("bar = %q, want %q", got, want)
			}
			if got, want := q.Get("limit"), "2"; got != want {
				t.Fatalf("limit = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			// sprd-candles: [ts,o,h,l,c,vol,confirm]
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[["1597026383085","3.721","3.743","3.677","3.708","8422410","0"],["1597026383086","3.731","3.799","3.494","3.72","24912403","1"]]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewMarketSprdCandlesService().SprdId("BTC-USDT_BTC-USDT-SWAP").Bar("1m").Limit(2).Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 2 {
			t.Fatalf("len = %d, want %d", len(got), 2)
		}
		if got[0].TS != 1597026383085 || got[0].Confirm != "0" {
			t.Fatalf("candle = %#v", got[0])
		}
		if got[0].VolCcy != "" || got[0].VolCcyQuote != "" {
			t.Fatalf("unexpected optional fields: VolCcy=%q VolCcyQuote=%q", got[0].VolCcy, got[0].VolCcyQuote)
		}
	})
}
