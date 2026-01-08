package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMarketTradesService_Do(t *testing.T) {
	t.Run("missing_instId", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewMarketTradesService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMarketTradesMissingInstId {
			t.Fatalf("error = %v, want %v", err, errMarketTradesMissingInstId)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/market/trades"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "instId=BTC-USDT&limit=2"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instId":"BTC-USDT","tradeId":"1","px":"100","sz":"1","side":"buy","ts":"1597026383085"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewMarketTradesService().InstId("BTC-USDT").Limit(2).Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("len = %d, want %d", len(got), 1)
		}
		if got[0].TradeId != "1" {
			t.Fatalf("TradeId = %q, want %q", got[0].TradeId, "1")
		}
		if got[0].TS != 1597026383085 {
			t.Fatalf("TS = %d, want %d", got[0].TS, 1597026383085)
		}
	})
}
