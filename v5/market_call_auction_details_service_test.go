package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMarketCallAuctionDetailsService_Do(t *testing.T) {
	t.Run("missing_instId", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewMarketCallAuctionDetailsService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMarketCallAuctionDetailsMissingInstId {
			t.Fatalf("error = %v, want %v", err, errMarketCallAuctionDetailsMissingInstId)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/market/call-auction-details"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "instId=ONDO-USDC"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instId":"ONDO-USDC","unmatchedSz":"9988764","eqPx":"0.6","matchedSz":"44978","state":"continuous_trading","auctionEndTime":"1726542000000","ts":"1726542000007"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewMarketCallAuctionDetailsService().InstId("ONDO-USDC").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.InstId != "ONDO-USDC" || got.EqPx != "0.6" || got.MatchedSz != "44978" || got.UnmatchedSz != "9988764" {
			t.Fatalf("details = %#v", got)
		}
		if got.AuctionEndTime != 1726542000000 || got.TS != 1726542000007 {
			t.Fatalf("times = %#v", got)
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

		_, err := c.NewMarketCallAuctionDetailsService().InstId("ONDO-USDC").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errEmptyMarketCallAuctionDetailsResponse {
			t.Fatalf("error = %v, want %v", err, errEmptyMarketCallAuctionDetailsResponse)
		}
	})
}
