package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublicFundingRateHistoryService_Do(t *testing.T) {
	t.Run("missing_instId", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewPublicFundingRateHistoryService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errPublicFundingRateHistoryMissingInstId {
			t.Fatalf("error = %v, want %v", err, errPublicFundingRateHistoryMissingInstId)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/public/funding-rate-history"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "after=1&before=2&instId=BTC-USD-SWAP&limit=3"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instType":"SWAP","instId":"BTC-USD-SWAP","formulaType":"noRate","fundingRate":"0.0000746604960499","realizedRate":"0.0002279755647389","fundingTime":"1703059200000","method":"next_period"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewPublicFundingRateHistoryService().
			InstId("BTC-USD-SWAP").
			After("1").
			Before("2").
			Limit(3).
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].InstId != "BTC-USD-SWAP" || got[0].FundingTime != 1703059200000 {
			t.Fatalf("data = %#v", got)
		}
		if got[0].FundingRate == "" || got[0].RealizedRate == "" {
			t.Fatalf("rate = %#v", got[0])
		}
	})
}
