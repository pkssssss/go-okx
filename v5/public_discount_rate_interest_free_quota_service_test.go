package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublicDiscountRateInterestFreeQuotaService_Do(t *testing.T) {
	t.Run("ok_empty_query", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/public/discount-rate-interest-free-quota"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, ""; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"amt":"0","ccy":"BTC","colRes":"0","collateralRestrict":false,"details":[{"discountRate":"0.98","liqPenaltyRate":"0.02","maxAmt":"20","minAmt":"0","tier":"1","disCcyEq":"1000"}],"discountLv":"1","minDiscountRate":"0"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewPublicDiscountRateInterestFreeQuotaService().Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].Ccy != "BTC" || got[0].Amt != "0" || got[0].ColRes != "0" {
			t.Fatalf("data = %#v", got)
		}
		if len(got[0].Details) != 1 || got[0].Details[0].DiscountRate == "" {
			t.Fatalf("details = %#v", got[0])
		}
	})

	t.Run("ok_with_ccy", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/public/discount-rate-interest-free-quota"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "ccy=BTC&discountLv=1"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"amt":"0","ccy":"BTC","colRes":"0","collateralRestrict":false,"details":[],"discountLv":"1","minDiscountRate":"0"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewPublicDiscountRateInterestFreeQuotaService().Ccy("BTC").DiscountLv("1").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].Ccy != "BTC" || got[0].DiscountLv != "1" {
			t.Fatalf("data = %#v", got)
		}
	})
}
