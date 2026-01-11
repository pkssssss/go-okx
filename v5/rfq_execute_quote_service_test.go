package okx

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRFQExecuteQuoteService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	t.Run("signed_request_and_body", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/rfq/execute-quote"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}

			bodyBytes, _ := io.ReadAll(r.Body)
			if got, want := string(bodyBytes), `{"rfqId":"22540","quoteId":"84073"}`; got != want {
				t.Fatalf("body = %q, want %q", got, want)
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
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"cTime":"1650966816550","rfqId":"22540","quoteId":"84073","blockTdId":"439121886014849026","tag":"123456","tTraderCode":"SATS","mTraderCode":"MIKE","isSuccessful":true,"errorCode":"","legs":[{"instId":"BTC-USDT","side":"buy","sz":"0.532","px":"100","tradeId":"439121886014849026","fee":"-0.0266","feeCcy":"USDT","tradeQuoteCcy":"USDT"}]}]}`))
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

		trade, err := c.NewRFQExecuteQuoteService().RfqId("22540").QuoteId("84073").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if trade == nil || trade.BlockTdId != "439121886014849026" || len(trade.Legs) != 1 {
			t.Fatalf("trade = %#v", trade)
		}
		if trade.Legs[0].TradeQuoteCcy != "USDT" {
			t.Fatalf("trade.Leg[0] = %#v", trade.Legs[0])
		}
	})

	t.Run("missing_required", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewRFQExecuteQuoteService().Do(context.Background())
		if !errors.Is(err, errRFQExecuteQuoteMissingRequired) {
			t.Fatalf("expected errRFQExecuteQuoteMissingRequired, got %T: %v", err, err)
		}
	})
}
