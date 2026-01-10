package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublicEstimatedPriceService_Do(t *testing.T) {
	t.Run("missing_instId", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewPublicEstimatedPriceService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errPublicEstimatedPriceMissingInstId {
			t.Fatalf("error = %v, want %v", err, errPublicEstimatedPriceMissingInstId)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/public/estimated-price"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "instId=BTC-USD-200214"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instType":"FUTURES","instId":"BTC-USD-200214","settlePx":"200","ts":"1597026383085"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewPublicEstimatedPriceService().InstId("BTC-USD-200214").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].InstId != "BTC-USD-200214" || got[0].TS != 1597026383085 {
			t.Fatalf("data = %#v", got)
		}
		if got[0].SettlePx != "200" {
			t.Fatalf("SettlePx = %q, want %q", got[0].SettlePx, "200")
		}
	})
}
