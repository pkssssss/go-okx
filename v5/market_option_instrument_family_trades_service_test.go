package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMarketOptionInstrumentFamilyTradesService_Do(t *testing.T) {
	t.Run("missing_instFamily", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewMarketOptionInstrumentFamilyTradesService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMarketOptionInstrumentFamilyTradesMissingInstFamily {
			t.Fatalf("error = %v, want %v", err, errMarketOptionInstrumentFamilyTradesMissingInstFamily)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/market/option/instrument-family-trades"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "instFamily=BTC-USD"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"vol24h":"103381","tradeInfo":[{"instId":"BTC-USD-221111-17750-C","side":"sell","sz":"1","px":"0.0075","tradeId":"20","ts":"1668090715058"}],"optType":"C"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewMarketOptionInstrumentFamilyTradesService().InstFamily("BTC-USD").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].OptType != "C" || got[0].Vol24h != "103381" {
			t.Fatalf("data = %#v", got)
		}
		if len(got[0].TradeInfo) != 1 || got[0].TradeInfo[0].InstId != "BTC-USD-221111-17750-C" || got[0].TradeInfo[0].TS != 1668090715058 {
			t.Fatalf("tradeInfo = %#v", got[0].TradeInfo)
		}
	})
}
