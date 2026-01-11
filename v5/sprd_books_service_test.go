package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSprdBooksService_Do(t *testing.T) {
	t.Run("missing_sprdId", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewSprdBooksService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errSprdBooksMissingSprdId {
			t.Fatalf("error = %v, want %v", err, errSprdBooksMissingSprdId)
		}
	})

	t.Run("ok_default_sz", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/sprd/books"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "sprdId=BTC-USDT_BTC-USDT-SWAP"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
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

		ob, err := c.NewSprdBooksService().
			SprdId("BTC-USDT_BTC-USDT-SWAP").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if ob == nil || ob.TS != 1629966436396 {
			t.Fatalf("orderBook = %#v", ob)
		}
		if len(ob.Asks) != 1 || ob.Asks[0].Px != "41006.8" || ob.Asks[0].NumOrders != "1" {
			t.Fatalf("asks = %#v", ob.Asks)
		}
		if len(ob.Bids) != 1 || ob.Bids[0].Px != "41006.3" || ob.Bids[0].NumOrders != "2" {
			t.Fatalf("bids = %#v", ob.Bids)
		}
	})

	t.Run("ok_with_sz", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/sprd/books"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "sprdId=BTC-USDT_BTC-USDT-SWAP&sz=20"; got != want {
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

		_, err := c.NewSprdBooksService().
			SprdId("BTC-USDT_BTC-USDT-SWAP").
			Sz(20).
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errEmptySprdBooksResponse {
			t.Fatalf("error = %v, want %v", err, errEmptySprdBooksResponse)
		}
	})
}
