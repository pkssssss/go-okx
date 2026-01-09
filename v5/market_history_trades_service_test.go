package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMarketHistoryTradesService_Do(t *testing.T) {
	t.Run("missing_instId", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewMarketHistoryTradesService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMarketHistoryTradesMissingInstId {
			t.Fatalf("error = %v, want %v", err, errMarketHistoryTradesMissingInstId)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/market/history-trades"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			q := r.URL.Query()
			if got, want := q.Get("instId"), "BTC-USDT"; got != want {
				t.Fatalf("instId = %q, want %q", got, want)
			}
			if got, want := q.Get("type"), "2"; got != want {
				t.Fatalf("type = %q, want %q", got, want)
			}
			if got, want := q.Get("after"), "100"; got != want {
				t.Fatalf("after = %q, want %q", got, want)
			}
			if got, want := q.Get("before"), "200"; got != want {
				t.Fatalf("before = %q, want %q", got, want)
			}
			if got, want := q.Get("limit"), "3"; got != want {
				t.Fatalf("limit = %q, want %q", got, want)
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

		got, err := c.NewMarketHistoryTradesService().
			InstId("BTC-USDT").
			Type("2").
			After("100").
			Before("200").
			Limit(3).
			Do(context.Background())
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
