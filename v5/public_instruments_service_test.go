package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublicInstrumentsService_Do(t *testing.T) {
	t.Run("missing_instType", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewPublicInstrumentsService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errPublicInstrumentsMissingInstType {
			t.Fatalf("error = %v, want %v", err, errPublicInstrumentsMissingInstType)
		}
	})

	t.Run("ok_minimal", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/public/instruments"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "instType=SPOT"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instType":"SPOT","instId":"BTC-USDT","instIdCode":211874,"baseCcy":"BTC","quoteCcy":"USDT","tickSz":"0.1","lotSz":"0.0001","minSz":"0.0001","state":"live"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewPublicInstrumentsService().InstType("SPOT").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("len = %d, want %d", len(got), 1)
		}
		if got[0].InstId != "BTC-USDT" {
			t.Fatalf("InstId = %q, want %q", got[0].InstId, "BTC-USDT")
		}
		if got[0].InstIdCode == nil {
			t.Fatalf("InstIdCode = nil, want %d", 211874)
		}
		if *got[0].InstIdCode != 211874 {
			t.Fatalf("InstIdCode = %d, want %d", *got[0].InstIdCode, 211874)
		}
		if got[0].TickSz != "0.1" {
			t.Fatalf("TickSz = %q, want %q", got[0].TickSz, "0.1")
		}
	})

	t.Run("with_filters", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/public/instruments"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "instFamily=BTC-USD&instId=BTC-USD-SWAP&instType=SWAP&uly=BTC-USD"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instType":"SWAP","instId":"BTC-USD-SWAP","instFamily":"BTC-USD","uly":"BTC-USD","settleCcy":"BTC","ctVal":"0.01","ctValCcy":"BTC","tickSz":"0.1","lotSz":"1","minSz":"1","state":"live"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewPublicInstrumentsService().
			InstType("SWAP").
			Uly("BTC-USD").
			InstFamily("BTC-USD").
			InstId("BTC-USD-SWAP").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("len = %d, want %d", len(got), 1)
		}
		if got[0].InstFamily != "BTC-USD" {
			t.Fatalf("InstFamily = %q, want %q", got[0].InstFamily, "BTC-USD")
		}
		if got[0].CtVal != "0.01" {
			t.Fatalf("CtVal = %q, want %q", got[0].CtVal, "0.01")
		}
	})
}
