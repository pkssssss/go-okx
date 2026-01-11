package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRFQRfqsService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, http.MethodGet; got != want {
			t.Fatalf("method = %q, want %q", got, want)
		}
		if got, want := r.URL.Path, "/api/v5/rfq/rfqs"; got != want {
			t.Fatalf("path = %q, want %q", got, want)
		}

		if got, want := r.URL.Query().Get("rfqId"), "22534"; got != want {
			t.Fatalf("rfqId = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("state"), "active"; got != want {
			t.Fatalf("state = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("limit"), "2"; got != want {
			t.Fatalf("limit = %q, want %q", got, want)
		}

		if got, want := r.Header.Get("OK-ACCESS-KEY"), "mykey"; got != want {
			t.Fatalf("OK-ACCESS-KEY = %q, want %q", got, want)
		}
		if got, want := r.Header.Get("OK-ACCESS-PASSPHRASE"), "mypass"; got != want {
			t.Fatalf("OK-ACCESS-PASSPHRASE = %q, want %q", got, want)
		}
		if got, want := r.Header.Get("OK-ACCESS-TIMESTAMP"), "2020-03-28T12:21:41.274Z"; got != want {
			t.Fatalf("OK-ACCESS-TIMESTAMP = %q, want %q", got, want)
		}
		if got := r.Header.Get("OK-ACCESS-SIGN"); got == "" {
			t.Fatalf("OK-ACCESS-SIGN empty")
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"cTime":"1611033737572","uTime":"1611033737572","rfqId":"22534","clRfqId":"rfq01","state":"active","validUntil":"1611033857557","allowPartialExecution":false,"counterparties":["Trader1"],"legs":[{"instId":"BTC-USDT","sz":"25","side":"buy"}]}]}`))
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

	rfqs, err := c.NewRFQRfqsService().RfqId("22534").State("active").Limit(2).Do(context.Background())
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if len(rfqs) != 1 || rfqs[0].RfqId != "22534" {
		t.Fatalf("rfqs = %#v", rfqs)
	}
}
