package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublicInsuranceFundService_Do(t *testing.T) {
	t.Run("missing_instType", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewPublicInsuranceFundService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errPublicInsuranceFundMissingInstType {
			t.Fatalf("error = %v, want %v", err, errPublicInsuranceFundMissingInstType)
		}
	})

	t.Run("missing_ccy_for_margin", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewPublicInsuranceFundService().InstType("MARGIN").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errPublicInsuranceFundMissingCcy {
			t.Fatalf("error = %v, want %v", err, errPublicInsuranceFundMissingCcy)
		}
	})

	t.Run("missing_instFamily_or_uly", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewPublicInsuranceFundService().InstType("SWAP").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errPublicInsuranceFundMissingInstFamilyOrUly {
			t.Fatalf("error = %v, want %v", err, errPublicInsuranceFundMissingInstFamilyOrUly)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/public/insurance-fund"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			wantQuery := "after=1&before=2&instType=SWAP&limit=3&type=regular_update&uly=BTC-USD"
			if got, want := r.URL.RawQuery, wantQuery; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"details":[{"adlType":"","amt":"0","balance":"3849.6616","ccy":"BTC","decRate":"","maxBal":"","maxBalTs":"","ts":"1768029918000","type":"regular_update"}],"instFamily":"BTC-USD","instType":"SWAP","total":"2689789954.6579"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewPublicInsuranceFundService().
			InstType("SWAP").
			Type("regular_update").
			Uly("BTC-USD").
			After("1").
			Before("2").
			Limit(3).
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].InstFamily != "BTC-USD" || got[0].Total == "" || len(got[0].Details) != 1 {
			t.Fatalf("data = %#v", got)
		}
		if got[0].Details[0].Ccy != "BTC" || got[0].Details[0].TS != 1768029918000 {
			t.Fatalf("details = %#v", got[0].Details)
		}
	})
}
