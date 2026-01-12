package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMarketTickerService_Do(t *testing.T) {
	t.Run("missing_instId", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewMarketTickerService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMarketTickerMissingInstId {
			t.Fatalf("error = %v, want %v", err, errMarketTickerMissingInstId)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/market/ticker"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "instId=BTC-USDT"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}

			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instType":"SPOT","instId":"BTC-USDT","last":"100","lastSz":"1","askPx":"101","askSz":"2","bidPx":"99","bidSz":"3","open24h":"90","high24h":"110","low24h":"80","volCcy24h":"1","vol24h":"2","sodUtc0":"0","sodUtc8":"0","ts":"1597026383085"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		tk, err := c.NewMarketTickerService().InstId("BTC-USDT").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if tk.InstId != "BTC-USDT" {
			t.Fatalf("InstId = %q, want %q", tk.InstId, "BTC-USDT")
		}
		if tk.TS != 1597026383085 {
			t.Fatalf("TS = %d, want %d", tk.TS, 1597026383085)
		}
	})

	t.Run("empty_data", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		_, err := c.NewMarketTickerService().InstId("BTC-USDT").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errEmptyMarketTickerResponse {
			t.Fatalf("error = %v, want %v", err, errEmptyMarketTickerResponse)
		}
	})
}

func TestMarketTickersService_Do(t *testing.T) {
	t.Run("missing_instType", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewMarketTickersService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMarketTickersMissingInstType {
			t.Fatalf("error = %v, want %v", err, errMarketTickersMissingInstType)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.URL.Path, "/api/v5/market/tickers"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "instType=SWAP"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instType":"SWAP","instId":"A","last":"1","ts":"1"},{"instType":"SWAP","instId":"B","last":"2","ts":"2"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewMarketTickersService().InstType("SWAP").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 2 {
			t.Fatalf("len = %d, want %d", len(got), 2)
		}
	})

	t.Run("with_instFamily", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.URL.Path, "/api/v5/market/tickers"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "instFamily=BTC-USD&instType=SWAP"; got != want {
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

		_, err := c.NewMarketTickersService().InstType("SWAP").InstFamily("BTC-USD").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
	})
}

func TestMarketBlockTickerService_Do(t *testing.T) {
	t.Run("missing_instId", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewMarketBlockTickerService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMarketBlockTickerMissingInstId {
			t.Fatalf("error = %v, want %v", err, errMarketBlockTickerMissingInstId)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/market/block-ticker"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "instId=BTC-USD-SWAP"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}

			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instType":"SWAP","instId":"BTC-USD-SWAP","volCcy24h":"2222","vol24h":"333","ts":"1597026383085"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		tk, err := c.NewMarketBlockTickerService().InstId("BTC-USD-SWAP").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if tk.InstId != "BTC-USD-SWAP" || tk.InstType != "SWAP" {
			t.Fatalf("ticker = %#v", tk)
		}
		if tk.TS != 1597026383085 {
			t.Fatalf("TS = %d, want %d", tk.TS, 1597026383085)
		}
	})

	t.Run("empty_data", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		_, err := c.NewMarketBlockTickerService().InstId("BTC-USD-SWAP").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errEmptyMarketBlockTickerResponse {
			t.Fatalf("error = %v, want %v", err, errEmptyMarketBlockTickerResponse)
		}
	})
}

func TestMarketBlockTickersService_Do(t *testing.T) {
	t.Run("missing_instType", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewMarketBlockTickersService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMarketBlockTickersMissingInstType {
			t.Fatalf("error = %v, want %v", err, errMarketBlockTickersMissingInstType)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.URL.Path, "/api/v5/market/block-tickers"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "instType=SWAP"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instType":"SWAP","instId":"A","volCcy24h":"1","vol24h":"2","ts":"1"},{"instType":"SWAP","instId":"B","volCcy24h":"3","vol24h":"4","ts":"2"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewMarketBlockTickersService().InstType("SWAP").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 2 {
			t.Fatalf("len = %d, want %d", len(got), 2)
		}
	})

	t.Run("with_instFamily", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.URL.Path, "/api/v5/market/block-tickers"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "instFamily=BTC-USD&instType=SWAP"; got != want {
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

		_, err := c.NewMarketBlockTickersService().InstType("SWAP").InstFamily("BTC-USD").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
	})
}

func TestMarketBooksService_Do(t *testing.T) {
	t.Run("missing_instId", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewMarketBooksService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMarketBooksMissingInstId {
			t.Fatalf("error = %v, want %v", err, errMarketBooksMissingInstId)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.URL.Path, "/api/v5/market/books"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "instId=BTC-USDT&sz=2"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"asks":[["41006.8","0.60038921","0","1"],["41007.0","1","0","2"]],"bids":[["41006.3","0.30178218","0","2"]],"ts":"1629966436396"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		ob, err := c.NewMarketBooksService().InstId("BTC-USDT").Sz(2).Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(ob.Asks) != 2 {
			t.Fatalf("asks len = %d, want %d", len(ob.Asks), 2)
		}
		if got, want := ob.Asks[0].Px, "41006.8"; got != want {
			t.Fatalf("ask[0].Px = %q, want %q", got, want)
		}
		if got, want := ob.Bids[0].NumOrders, "2"; got != want {
			t.Fatalf("bid[0].NumOrders = %q, want %q", got, want)
		}
		if ob.TS != 1629966436396 {
			t.Fatalf("ts = %d, want %d", ob.TS, 1629966436396)
		}
	})

	t.Run("empty_data", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		_, err := c.NewMarketBooksService().InstId("BTC-USDT").Sz(2).Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errEmptyMarketBooksResponse {
			t.Fatalf("error = %v, want %v", err, errEmptyMarketBooksResponse)
		}
	})
}
