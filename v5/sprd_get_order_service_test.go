package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSprdGetOrderService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	t.Run("signed_request_and_query_ordId", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/sprd/order"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "ordId=2510789768709120"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-KEY"), "mykey"; got != want {
				t.Fatalf("OK-ACCESS-KEY = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"sprdId":"BTC-USD-SWAP_BTC-USD-200329","ordId":"312269865356374016","clOrdId":"b1","tag":"","px":"999","sz":"3","ordType":"limit","side":"buy","fillSz":"0","fillPx":"","tradeId":"","accFillSz":"0","pendingFillSz":"2","pendingSettleSz":"1","canceledSz":"1","state":"live","avgPx":"0","cancelSource":"","uTime":"1597026383085","cTime":"1597026383085"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
			WithCredentials(Credentials{
				APIKey:     "mykey",
				SecretKey:  "mysecret",
				Passphrase: "mypass",
			}),
			WithNowFunc(func() time.Time { return fixedNow }),
		)

		got, err := c.NewSprdGetOrderService().OrdId("2510789768709120").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.OrdId != "312269865356374016" || got.SprdId == "" || got.State != "live" {
			t.Fatalf("order = %#v", got)
		}
		if got.UTime != 1597026383085 || got.CTime != 1597026383085 {
			t.Fatalf("time = %#v", got)
		}
	})
}
