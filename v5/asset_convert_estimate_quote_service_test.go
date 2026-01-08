package okx

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pkssssss/go-okx/v5/internal/sign"
)

func TestAssetConvertEstimateQuoteService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("missing_required", func(t *testing.T) {
		c := NewClient(
			WithNowFunc(func() time.Time { return fixedNow }),
		)
		_, err := c.NewAssetConvertEstimateQuoteService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errAssetConvertEstimateQuoteMissingRequired {
			t.Fatalf("error = %v, want %v", err, errAssetConvertEstimateQuoteMissingRequired)
		}
	})

	t.Run("signed_request_and_body", func(t *testing.T) {
		timestamp := sign.TimestampISO8601Millis(fixedNow)
		wantBody := `{"baseCcy":"ETH","quoteCcy":"USDT","side":"buy","rfqSz":"30","rfqSzCcy":"USDT"}`
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodPost, "/api/v5/asset/convert/estimate-quote", wantBody))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bodyBytes, _ := io.ReadAll(r.Body)
			if got, want := string(bodyBytes), wantBody; got != want {
				t.Fatalf("body = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), wantSig; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"baseCcy":"ETH","baseSz":"0.01023052","clQReqId":"","cnvtPx":"2932.40104429","origRfqSz":"30","quoteCcy":"USDT","quoteId":"quoterETH-USDT16461885104612381","quoteSz":"30","quoteTime":"1646188510461","rfqSz":"30","rfqSzCcy":"USDT","side":"buy","ttlMs":"10000"}]}`))
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

		got, err := c.NewAssetConvertEstimateQuoteService().
			BaseCcy("ETH").
			QuoteCcy("USDT").
			Side("buy").
			RfqSz("30").
			RfqSzCcy("USDT").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.QuoteId == "" || got.BaseCcy != "ETH" || got.QuoteCcy != "USDT" {
			t.Fatalf("data = %#v", got)
		}
	})

	t.Run("missing_credentials", func(t *testing.T) {
		c := NewClient(
			WithNowFunc(func() time.Time { return fixedNow }),
		)
		_, err := c.NewAssetConvertEstimateQuoteService().
			BaseCcy("ETH").
			QuoteCcy("USDT").
			Side("buy").
			RfqSz("30").
			RfqSzCcy("USDT").
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMissingCredentials {
			t.Fatalf("error = %v, want %v", err, errMissingCredentials)
		}
	})
}
