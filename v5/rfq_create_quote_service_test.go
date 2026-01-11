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

func TestRFQCreateQuoteService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	t.Run("signed_request_and_body", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/rfq/create-quote"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}

			bodyBytes, _ := io.ReadAll(r.Body)
			wantBody := `{"rfqId":"22539","clQuoteId":"q001","tag":"123456","quoteSide":"buy","anonymous":true,"expiresIn":"30","legs":[{"px":"39450.0","sz":"200000","instId":"BTC-USDT-SWAP","tdMode":"cross","ccy":"USDT","side":"buy","posSide":"long"},{"px":"39450.0","sz":"200000","instId":"ETH-USDT","tdMode":"cross","ccy":"USDT","side":"buy","posSide":"long","tgtCcy":"base_ccy","tradeQuoteCcy":"USDT"}]}`
			if got := string(bodyBytes); got != wantBody {
				t.Fatalf("body = %q, want %q", got, wantBody)
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
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"cTime":"1608267227834","uTime":"1608267227834","quoteId":"25092","rfqId":"18753","tag":"123456","quoteSide":"sell","state":"active","reason":"mmp_canceled","clQuoteId":"","clRfqId":"","traderCode":"Aksha","legs":[{"px":"46000","sz":"25","instId":"BTC-USD-220114-25000-C","tdMode":"cross","ccy":"USDT","side":"sell","posSide":"long"},{"px":"4000","sz":"25","instId":"ETH-USD-220114-25000-C","tdMode":"cross","ccy":"USDT","side":"buy","posSide":"long"}]}]}`))
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

		quote, err := c.NewRFQCreateQuoteService().
			RfqId("22539").
			ClQuoteId("q001").
			Tag("123456").
			QuoteSide("buy").
			Anonymous(true).
			ExpiresIn("30").
			Legs([]QuoteLeg{
				{Px: "39450.0", Sz: "200000", InstId: "BTC-USDT-SWAP", TdMode: "cross", Ccy: "USDT", Side: "buy", PosSide: "long"},
				{Px: "39450.0", Sz: "200000", InstId: "ETH-USDT", TdMode: "cross", Ccy: "USDT", Side: "buy", PosSide: "long", TgtCcy: "base_ccy", TradeQuoteCcy: "USDT"},
			}).
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if quote == nil || quote.QuoteId != "25092" || quote.RfqId != "18753" || quote.State != "active" {
			t.Fatalf("quote = %#v", quote)
		}
	})

	t.Run("missing_required", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewRFQCreateQuoteService().Do(context.Background())
		if !errors.Is(err, errRFQCreateQuoteMissingRequired) {
			t.Fatalf("expected errRFQCreateQuoteMissingRequired, got %T: %v", err, err)
		}
	})
}
