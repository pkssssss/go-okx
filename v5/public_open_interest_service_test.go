package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublicOpenInterestService_Do(t *testing.T) {
	t.Run("missing_instType", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewPublicOpenInterestService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errPublicOpenInterestMissingInstType {
			t.Fatalf("error = %v, want %v", err, errPublicOpenInterestMissingInstType)
		}
	})

	t.Run("ok_minimal", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/public/open-interest"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "instType=SWAP"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instId":"BTC-USDT-SWAP","instType":"SWAP","oi":"1","oiCcy":"2","oiUsd":"3","ts":"4"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewPublicOpenInterestService().InstType("SWAP").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("len = %d, want %d", len(got), 1)
		}
		if got[0].OI != "1" {
			t.Fatalf("OI = %q, want %q", got[0].OI, "1")
		}
		if got[0].TS != 4 {
			t.Fatalf("TS = %d, want %d", got[0].TS, 4)
		}
	})

	t.Run("with_filters", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.URL.Path, "/api/v5/public/open-interest"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "instFamily=BTC-USD&instId=BTC-USD-SWAP&instType=SWAP&uly=BTC-USD"; got != want {
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

		_, err := c.NewPublicOpenInterestService().
			InstType("SWAP").
			Uly("BTC-USD").
			InstFamily("BTC-USD").
			InstId("BTC-USD-SWAP").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
	})
}
