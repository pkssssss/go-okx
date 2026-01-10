package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublicPositionTiersService_Do(t *testing.T) {
	t.Run("missing_required", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewPublicPositionTiersService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errPublicPositionTiersMissingRequired {
			t.Fatalf("error = %v, want %v", err, errPublicPositionTiersMissingRequired)
		}
	})

	t.Run("missing_instFamily_for_derivatives", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewPublicPositionTiersService().InstType("SWAP").TdMode("cross").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errPublicPositionTiersMissingInstFamily {
			t.Fatalf("error = %v, want %v", err, errPublicPositionTiersMissingInstFamily)
		}
	})

	t.Run("missing_instId_or_ccy_for_margin", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewPublicPositionTiersService().InstType("MARGIN").TdMode("cross").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errPublicPositionTiersMissingInstIdOrCcy {
			t.Fatalf("error = %v, want %v", err, errPublicPositionTiersMissingInstIdOrCcy)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/public/position-tiers"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			wantQuery := "instFamily=BTC-USDT&instType=SWAP&tdMode=cross&tier=1"
			if got, want := r.URL.RawQuery, wantQuery; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"baseMaxLoan":"","imr":"0.01","instFamily":"BTC-USDT","instId":"","maxLever":"100","maxSz":"1000","minSz":"0","mmr":"0.004","optMgnFactor":"0","quoteMaxLoan":"","tier":"1","uly":"BTC-USDT"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewPublicPositionTiersService().
			InstType("SWAP").
			TdMode("cross").
			InstFamily("BTC-USDT").
			Tier("1").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].Tier != "1" || got[0].MaxLever != "100" || got[0].MMR != "0.004" {
			t.Fatalf("data = %#v", got)
		}
	})
}
