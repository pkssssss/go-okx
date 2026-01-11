package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRFQPublicTradesService_Do(t *testing.T) {
	t.Run("ok_empty_query", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/rfq/public-trades"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, ""; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"strategy":"CALL_CALENDAR_SPREAD","cTime":"1650976251241","blockTdId":"439161457415012352","groupId":"","legs":[{"instId":"BTC-USDT","side":"buy","sz":"0.1","px":"10.1","tradeId":"439161457415012356"}]}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewRFQPublicTradesService().Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("len(data) = %d, want 1", len(got))
		}
		if got[0].Strategy == "" || got[0].BlockTdId == "" || got[0].CTime == 0 {
			t.Fatalf("trade = %#v", got[0])
		}
		if len(got[0].Legs) != 1 || got[0].Legs[0].InstId != "BTC-USDT" {
			t.Fatalf("legs = %#v", got[0].Legs)
		}
	})

	t.Run("ok_with_query", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/rfq/public-trades"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "beginId=1&endId=9&limit=10"; got != want {
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

		_, err := c.NewRFQPublicTradesService().
			BeginId("1").
			EndId("9").
			Limit(10).
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
	})
}
