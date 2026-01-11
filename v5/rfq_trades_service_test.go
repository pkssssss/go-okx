package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRFQTradesService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, http.MethodGet; got != want {
			t.Fatalf("method = %q, want %q", got, want)
		}
		if got, want := r.URL.Path, "/api/v5/rfq/trades"; got != want {
			t.Fatalf("path = %q, want %q", got, want)
		}

		if got, want := r.URL.Query().Get("blockTdId"), "439121886014849026"; got != want {
			t.Fatalf("blockTdId = %q, want %q", got, want)
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
		_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"cTime":"1650966816550","rfqId":"22540","quoteId":"84073","blockTdId":"439121886014849026","tag":"123456","tTraderCode":"SATS","mTraderCode":"MIKE","isSuccessful":true,"errorCode":"","legs":[{"instId":"BTC-USDT","side":"sell","sz":"0.532","px":"100","tradeId":"439121886014849026","fee":"-0.0266","feeCcy":"USDT","tradeQuoteCcy":"USDT"}]}]}`))
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

	trades, err := c.NewRFQTradesService().BlockTdId("439121886014849026").Do(context.Background())
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if len(trades) != 1 || trades[0].BlockTdId != "439121886014849026" {
		t.Fatalf("trades = %#v", trades)
	}
	if trades[0].Legs[0].TradeQuoteCcy != "USDT" {
		t.Fatalf("legs[0] = %#v", trades[0].Legs[0])
	}
}
