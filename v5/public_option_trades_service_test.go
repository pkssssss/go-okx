package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublicOptionTradesService_Do(t *testing.T) {
	t.Run("missing_instId_or_instFamily", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewPublicOptionTradesService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errPublicOptionTradesMissingInstIdOrInstFamily {
			t.Fatalf("error = %v, want %v", err, errPublicOptionTradesMissingInstIdOrInstFamily)
		}
	})

	t.Run("ok_instFamily", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/public/option-trades"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "instFamily=BTC-USD&optType=P"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"fillVol":"0.24415013671875","fwdPx":"16676.907614127158","idxPx":"16667","instFamily":"BTC-USD","instId":"BTC-USD-221230-16600-P","markPx":"0.006308943261227884","optType":"P","px":"0.005","side":"sell","sz":"30","tradeId":"65","ts":"1672225112048"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewPublicOptionTradesService().
			InstFamily("BTC-USD").
			OptType("P").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].InstFamily != "BTC-USD" || got[0].OptType != "P" || got[0].TS != 1672225112048 {
			t.Fatalf("data = %#v", got)
		}
		if got[0].InstId == "" || got[0].TradeId == "" || got[0].Px == "" {
			t.Fatalf("trade = %#v", got[0])
		}
	})

	t.Run("ok_instId_prefer", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.URL.RawQuery, "instId=BTC-USD-221230-16600-P&optType=P"; got != want {
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

		_, err := c.NewPublicOptionTradesService().
			InstId("BTC-USD-221230-16600-P").
			InstFamily("BTC-USD").
			OptType("P").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
	})
}
