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

func TestRFQCreateRFQService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	t.Run("signed_request_and_body", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/rfq/create-rfq"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}

			bodyBytes, _ := io.ReadAll(r.Body)
			wantBody := `{"counterparties":["Trader1"],"anonymous":true,"allowPartialExecution":false,"clRfqId":"rfq01","tag":"123456","legs":[{"instId":"BTC-USDT","tdMode":"cross","ccy":"USDT","sz":"25","side":"buy","posSide":"long"},{"instId":"ETH-USDT","tdMode":"cross","ccy":"USDT","sz":"150","side":"buy","posSide":"long","tgtCcy":"base_ccy","tradeQuoteCcy":"USDT"}]}`
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
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"cTime":"1611033737572","uTime":"1611033737572","rfqId":"22534","clRfqId":"rfq01","state":"active","validUntil":"1611033857557","allowPartialExecution":false,"counterparties":["Trader1"],"legs":[{"instId":"BTC-USDT","sz":"25","side":"buy"},{"instId":"ETH-USDT","sz":"150","side":"buy"}]}]}`))
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

		rfq, err := c.NewRFQCreateRFQService().
			Counterparties([]string{"Trader1"}).
			Anonymous(true).
			AllowPartialExecution(false).
			ClRfqId("rfq01").
			Tag("123456").
			Legs([]RFQLeg{
				{InstId: "BTC-USDT", TdMode: "cross", Ccy: "USDT", Sz: "25", Side: "buy", PosSide: "long"},
				{InstId: "ETH-USDT", TdMode: "cross", Ccy: "USDT", Sz: "150", Side: "buy", PosSide: "long", TgtCcy: "base_ccy", TradeQuoteCcy: "USDT"},
			}).
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if rfq == nil || rfq.RfqId != "22534" || rfq.ClRfqId != "rfq01" || rfq.State != "active" {
			t.Fatalf("rfq = %#v", rfq)
		}
	})

	t.Run("missing_required", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewRFQCreateRFQService().Do(context.Background())
		if !errors.Is(err, errRFQCreateRFQMissingCounterparties) {
			t.Fatalf("expected errRFQCreateRFQMissingCounterparties, got %T: %v", err, err)
		}
	})
}
