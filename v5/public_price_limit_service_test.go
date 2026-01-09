package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublicPriceLimitService_Do(t *testing.T) {
	t.Run("missing_instId", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewPublicPriceLimitService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errPublicPriceLimitMissingInstId {
			t.Fatalf("error = %v, want %v", err, errPublicPriceLimitMissingInstId)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/public/price-limit"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "instId=BTC-USDT-SWAP"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instType":"SWAP","instId":"BTC-USDT-SWAP","buyLmt":"17057.9","sellLmt":"16388.9","ts":"1597026383085","enabled":true}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewPublicPriceLimitService().InstId("BTC-USDT-SWAP").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("len = %d, want %d", len(got), 1)
		}
		if got[0].BuyLmt != "17057.9" {
			t.Fatalf("BuyLmt = %q, want %q", got[0].BuyLmt, "17057.9")
		}
		if got[0].Enabled != true {
			t.Fatalf("Enabled = %v, want %v", got[0].Enabled, true)
		}
		if got[0].TS != 1597026383085 {
			t.Fatalf("TS = %d, want %d", got[0].TS, 1597026383085)
		}
	})
}
