package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublicFundingRateService_Do(t *testing.T) {
	t.Run("missing_instId", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewPublicFundingRateService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errPublicFundingRateMissingInstId {
			t.Fatalf("error = %v, want %v", err, errPublicFundingRateMissingInstId)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/public/funding-rate"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "instId=BTC-USDT-SWAP"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instType":"SWAP","instId":"BTC-USDT-SWAP","fundingRate":"0.1","fundingTime":"1","nextFundingRate":"","nextFundingTime":"2","minFundingRate":"-0.1","maxFundingRate":"0.1","interestRate":"0.1","premium":"0","method":"m","settState":"s","formulaType":"t","ts":"3"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewPublicFundingRateService().InstId("BTC-USDT-SWAP").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.InstId != "BTC-USDT-SWAP" {
			t.Fatalf("InstId = %q, want %q", got.InstId, "BTC-USDT-SWAP")
		}
		if got.FundingTime != 1 {
			t.Fatalf("FundingTime = %d, want %d", got.FundingTime, 1)
		}
		if got.NextFundingTime != 2 {
			t.Fatalf("NextFundingTime = %d, want %d", got.NextFundingTime, 2)
		}
		if got.TS != 3 {
			t.Fatalf("TS = %d, want %d", got.TS, 3)
		}
	})

	t.Run("empty_data", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		_, err := c.NewPublicFundingRateService().InstId("BTC-USDT-SWAP").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errEmptyFundingRateResponse {
			t.Fatalf("error = %v, want %v", err, errEmptyFundingRateResponse)
		}
	})
}
