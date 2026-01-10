package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublicBlockTradesService_Do(t *testing.T) {
	t.Run("missing_instId", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewPublicBlockTradesService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errPublicBlockTradesMissingInstId {
			t.Fatalf("error = %v, want %v", err, errPublicBlockTradesMissingInstId)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/public/block-trades"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "instId=BTC-USDT"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"fillVol":"","fwdPx":"","groupId":"","idxPx":"","instId":"BTC-USDT","markPx":"","px":"65000","side":"buy","sz":"0.1","tradeId":"1","ts":"1697181568974"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewPublicBlockTradesService().InstId("BTC-USDT").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].InstId != "BTC-USDT" || got[0].TradeId != "1" || got[0].TS != 1697181568974 {
			t.Fatalf("data = %#v", got)
		}
		if got[0].Px == "" || got[0].Sz == "" || got[0].Side == "" {
			t.Fatalf("trade = %#v", got[0])
		}
	})
}
