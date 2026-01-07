package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublicOptSummaryService_Do(t *testing.T) {
	t.Run("missing_uly", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewPublicOptSummaryService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errPublicOptSummaryMissingUly {
			t.Fatalf("error = %v, want %v", err, errPublicOptSummaryMissingUly)
		}
	})

	t.Run("ok_minimal", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/public/opt-summary"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "uly=BTC-USD"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instType":"OPTION","instId":"BTC-USD-260123-106000-P","uly":"BTC-USD","askVol":"0","bidVol":"0","markVol":"0.4","realVol":"0","delta":"-1","gamma":"0","theta":"0","vega":"0","volLv":"0.4","fwdPx":"1","distance":"0","ts":"2"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewPublicOptSummaryService().Uly("BTC-USD").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("len = %d, want %d", len(got), 1)
		}
		if got[0].InstId != "BTC-USD-260123-106000-P" {
			t.Fatalf("InstId = %q, want %q", got[0].InstId, "BTC-USD-260123-106000-P")
		}
		if got[0].TS != 2 {
			t.Fatalf("TS = %d, want %d", got[0].TS, 2)
		}
	})

	t.Run("with_expTime", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.URL.Path, "/api/v5/public/opt-summary"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "expTime=260123&uly=BTC-USD"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		_, err := c.NewPublicOptSummaryService().Uly("BTC-USD").ExpTime(260123).Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
	})
}
