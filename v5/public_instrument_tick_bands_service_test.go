package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublicInstrumentTickBandsService_Do(t *testing.T) {
	t.Run("missing_instType", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewPublicInstrumentTickBandsService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errPublicInstrumentTickBandsMissingInstType {
			t.Fatalf("error = %v, want %v", err, errPublicInstrumentTickBandsMissingInstType)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/public/instrument-tick-bands"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "instFamily=BTC-USD&instType=OPTION"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instType":"OPTION","instFamily":"BTC-USD","tickBand":[{"minPx":"0","maxPx":"100","tickSz":"0.1"},{"minPx":"100","maxPx":"10000","tickSz":"1"}]}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewPublicInstrumentTickBandsService().InstType("OPTION").InstFamily("BTC-USD").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].InstFamily != "BTC-USD" || got[0].InstType != "OPTION" || len(got[0].TickBands) != 2 {
			t.Fatalf("data = %#v", got)
		}
		if got[0].TickBands[0].TickSz != "0.1" || got[0].TickBands[1].MaxPx != "10000" {
			t.Fatalf("tickBands = %#v", got[0].TickBands)
		}
	})
}
