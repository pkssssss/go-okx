package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSprdSpreadsService_Do(t *testing.T) {
	t.Run("ok_empty_query", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/sprd/spreads"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, ""; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"sprdId":"BTC-USDT_BTC-USDT-SWAP","sprdType":"linear","state":"live","baseCcy":"BTC","szCcy":"BTC","quoteCcy":"USDT","tickSz":"0.0001","minSz":"0.001","lotSz":"1","listTime":"1597026383085","expTime":"1597029999085","uTime":"1597028888085","legs":[{"instId":"BTC-USDT","side":"sell"},{"instId":"BTC-USDT-SWAP","side":"buy"}]}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewSprdSpreadsService().Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].SprdId == "" || got[0].State != "live" {
			t.Fatalf("data = %#v", got)
		}
		if got[0].ListTime != 1597026383085 || got[0].UTime != 1597028888085 {
			t.Fatalf("time = %#v", got[0])
		}
		if len(got[0].Legs) != 2 || got[0].Legs[0].InstId == "" || got[0].Legs[1].InstId == "" {
			t.Fatalf("legs = %#v", got[0].Legs)
		}
	})

	t.Run("ok_with_filters", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/sprd/spreads"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "baseCcy=BTC&instId=BTC-USDT&sprdId=BTC-USDT_BTC-USDT-SWAP&state=live"; got != want {
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

		_, err := c.NewSprdSpreadsService().
			BaseCcy("BTC").
			InstId("BTC-USDT").
			SprdId("BTC-USDT_BTC-USDT-SWAP").
			State("live").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
	})
}
