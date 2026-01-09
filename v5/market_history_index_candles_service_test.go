package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMarketHistoryIndexCandlesService_Do(t *testing.T) {
	t.Run("missing_instId", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewMarketHistoryIndexCandlesService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMarketHistoryIndexCandlesMissingInstId {
			t.Fatalf("error = %v, want %v", err, errMarketHistoryIndexCandlesMissingInstId)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/market/history-index-candles"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			q := r.URL.Query()
			if got, want := q.Get("instId"), "BTC-USD"; got != want {
				t.Fatalf("instId = %q, want %q", got, want)
			}
			if got, want := q.Get("bar"), "1H"; got != want {
				t.Fatalf("bar = %q, want %q", got, want)
			}
			if got, want := q.Get("after"), "1000"; got != want {
				t.Fatalf("after = %q, want %q", got, want)
			}
			if got, want := q.Get("before"), "2000"; got != want {
				t.Fatalf("before = %q, want %q", got, want)
			}
			if got, want := q.Get("limit"), "2"; got != want {
				t.Fatalf("limit = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[["1629966436396","1","2","0.5","1.5","1"]]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewMarketHistoryIndexCandlesService().
			InstId("BTC-USD").
			Bar("1H").
			After("1000").
			Before("2000").
			Limit(2).
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("len = %d, want %d", len(got), 1)
		}
		if got[0].TS != 1629966436396 {
			t.Fatalf("TS = %d, want %d", got[0].TS, 1629966436396)
		}
		if got[0].Confirm != "1" {
			t.Fatalf("Confirm = %q, want %q", got[0].Confirm, "1")
		}
	})
}
