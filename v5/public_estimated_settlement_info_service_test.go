package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublicEstimatedSettlementInfoService_Do(t *testing.T) {
	t.Run("missing_instId", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewPublicEstimatedSettlementInfoService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errPublicEstimatedSettlementInfoMissingInstId {
			t.Fatalf("error = %v, want %v", err, errPublicEstimatedSettlementInfoMissingInstId)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/public/estimated-settlement-info"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "instId=XRP-USDT-250307"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"estSettlePx":"2.5666068562369959","instId":"XRP-USDT-250307","nextSettleTime":"1741248000000","ts":"1741246429748"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewPublicEstimatedSettlementInfoService().InstId("XRP-USDT-250307").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].InstId != "XRP-USDT-250307" || got[0].NextSettleTime != 1741248000000 || got[0].TS != 1741246429748 {
			t.Fatalf("data = %#v", got)
		}
		if got[0].EstSettlePx == "" {
			t.Fatalf("price = %#v", got[0])
		}
	})
}
