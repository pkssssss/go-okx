package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSprdPublicTradesService_Do(t *testing.T) {
	t.Run("ok_empty_query", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/sprd/public-trades"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, ""; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"sprdId":"BTC-USDT_BTC-USDC-SWAP","side":"sell","sz":"0.1","px":"964.1","tradeId":"242720719","ts":"1654161641568"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewSprdPublicTradesService().Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].SprdId == "" || got[0].TradeId == "" || got[0].TS != 1654161641568 {
			t.Fatalf("data = %#v", got)
		}
		if got[0].Px == "" || got[0].Sz == "" || got[0].Side == "" {
			t.Fatalf("trade = %#v", got[0])
		}
	})

	t.Run("ok_with_sprdId", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/sprd/public-trades"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "sprdId=BTC-USDT_BTC-USDT-SWAP"; got != want {
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

		_, err := c.NewSprdPublicTradesService().
			SprdId("BTC-USDT_BTC-USDT-SWAP").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
	})
}
