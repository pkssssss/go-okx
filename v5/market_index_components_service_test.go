package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMarketIndexComponentsService_Do(t *testing.T) {
	t.Run("missing_index", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewMarketIndexComponentsService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMarketIndexComponentsMissingIndex {
			t.Fatalf("error = %v, want %v", err, errMarketIndexComponentsMissingIndex)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/market/index-components"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "index=BTC-USD"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":{"components":[{"symbol":"BTC/USDT","symPx":"52733.2","wgt":"0.25","cnvPx":"52733.2","exch":"OKEx"}],"last":"52735.4123234925","index":"BTC-USD","ts":"1630985335599"}}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewMarketIndexComponentsService().Index("BTC-USD").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.Index != "BTC-USD" || got.Last != "52735.4123234925" || got.TS != 1630985335599 {
			t.Fatalf("data = %#v", got)
		}
		if len(got.Components) != 1 || got.Components[0].Exch != "OKEx" || got.Components[0].Symbol != "BTC/USDT" {
			t.Fatalf("components = %#v", got.Components)
		}
	})

	t.Run("empty_object", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":{}}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		_, err := c.NewMarketIndexComponentsService().Index("BTC-USD").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errEmptyMarketIndexComponentsResponse {
			t.Fatalf("error = %v, want %v", err, errEmptyMarketIndexComponentsResponse)
		}
	})
}
