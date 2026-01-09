package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMarketIndexTickersService_Do(t *testing.T) {
	t.Run("missing_quoteCcy_and_instId", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewMarketIndexTickersService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMarketIndexTickersMissingQuoteCcyOrInstId {
			t.Fatalf("error = %v, want %v", err, errMarketIndexTickersMissingQuoteCcyOrInstId)
		}
	})

	t.Run("ok_by_instId", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/market/index-tickers"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "instId=BTC-USDT"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instId":"BTC-USDT","idxPx":"43350","high24h":"43649.7","sodUtc0":"43444.1","open24h":"43640.8","low24h":"43261.9","sodUtc8":"43328.7","ts":"1649419644492"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewMarketIndexTickersService().InstId("BTC-USDT").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("len = %d, want %d", len(got), 1)
		}
		if got[0].IdxPx != "43350" {
			t.Fatalf("IdxPx = %q, want %q", got[0].IdxPx, "43350")
		}
		if got[0].TS != 1649419644492 {
			t.Fatalf("TS = %d, want %d", got[0].TS, 1649419644492)
		}
	})

	t.Run("ok_by_quoteCcy", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.URL.Path, "/api/v5/market/index-tickers"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "quoteCcy=USDT"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		_, err := c.NewMarketIndexTickersService().QuoteCcy("USDT").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
	})
}
