package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublicSettlementHistoryService_Do(t *testing.T) {
	t.Run("missing_instFamily", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewPublicSettlementHistoryService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errPublicSettlementHistoryMissingInstFamily {
			t.Fatalf("error = %v, want %v", err, errPublicSettlementHistoryMissingInstFamily)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/public/settlement-history"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "after=1&before=2&instFamily=BTC-USDT&limit=3"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"details":[{"instId":"XRP-USDT-250307","settlePx":"2.5192078615298715"}],"ts":"1741161600000"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewPublicSettlementHistoryService().
			InstFamily("BTC-USDT").
			After("1").
			Before("2").
			Limit(3).
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].TS != 1741161600000 {
			t.Fatalf("data = %#v", got)
		}
		if len(got[0].Details) != 1 || got[0].Details[0].InstId == "" || got[0].Details[0].SettlePx == "" {
			t.Fatalf("details = %#v", got[0])
		}
	})
}
