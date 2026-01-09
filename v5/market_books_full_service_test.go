package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMarketBooksFullService_Do(t *testing.T) {
	t.Run("missing_instId", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewMarketBooksFullService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMarketBooksFullMissingInstId {
			t.Fatalf("error = %v, want %v", err, errMarketBooksFullMissingInstId)
		}
	})

	t.Run("empty_response", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		_, err := c.NewMarketBooksFullService().InstId("BTC-USDT").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errEmptyMarketBooksFullResponse {
			t.Fatalf("error = %v, want %v", err, errEmptyMarketBooksFullResponse)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/market/books-full"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			q := r.URL.Query()
			if got, want := q.Get("instId"), "BTC-USDT"; got != want {
				t.Fatalf("instId = %q, want %q", got, want)
			}
			if got, want := q.Get("sz"), "20"; got != want {
				t.Fatalf("sz = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"asks":[["41006.8","0.60038921","1"]],"bids":[["41006.3","0.30178218","2"]],"ts":"1629966436396"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		ob, err := c.NewMarketBooksFullService().InstId("BTC-USDT").Sz(20).Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if ob.TS != 1629966436396 {
			t.Fatalf("TS = %d, want %d", ob.TS, 1629966436396)
		}
		if len(ob.Asks) != 1 || ob.Asks[0].Px != "41006.8" || ob.Asks[0].Sz != "0.60038921" {
			t.Fatalf("asks = %#v, want 1 item with px=41006.8 sz=0.60038921", ob.Asks)
		}
		if ob.Asks[0].LiqOrd != "" || ob.Asks[0].NumOrders != "1" {
			t.Fatalf("asks[0] liqOrd/numOrders = %q/%q, want %q/%q", ob.Asks[0].LiqOrd, ob.Asks[0].NumOrders, "", "1")
		}
		if len(ob.Bids) != 1 || ob.Bids[0].Px != "41006.3" || ob.Bids[0].Sz != "0.30178218" {
			t.Fatalf("bids = %#v, want 1 item with px=41006.3 sz=0.30178218", ob.Bids)
		}
		if ob.Bids[0].LiqOrd != "" || ob.Bids[0].NumOrders != "2" {
			t.Fatalf("bids[0] liqOrd/numOrders = %q/%q, want %q/%q", ob.Bids[0].LiqOrd, ob.Bids[0].NumOrders, "", "2")
		}
	})
}
