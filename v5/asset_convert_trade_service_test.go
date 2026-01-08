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

func TestAssetConvertTradeService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("missing_required", func(t *testing.T) {
		c := NewClient(
			WithNowFunc(func() time.Time { return fixedNow }),
		)
		_, err := c.NewAssetConvertTradeService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errAssetConvertTradeMissingRequired {
			t.Fatalf("error = %v, want %v", err, errAssetConvertTradeMissingRequired)
		}
	})

	t.Run("signed_request_and_body", func(t *testing.T) {
		timestamp := sign.TimestampISO8601Millis(fixedNow)
		wantBody := `{"quoteId":"quoterETH-USDT16461885104612381","baseCcy":"ETH","quoteCcy":"USDT","side":"buy","sz":"30","szCcy":"USDT"}`
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodPost, "/api/v5/asset/convert/trade", wantBody))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bodyBytes, _ := io.ReadAll(r.Body)
			if got, want := string(bodyBytes), wantBody; got != want {
				t.Fatalf("body = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), wantSig; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"baseCcy":"ETH","clTReqId":"","fillBaseSz":"0.01023052","fillPx":"2932.40104429","fillQuoteSz":"30","instId":"ETH-USDT","quoteCcy":"USDT","quoteId":"quoterETH-USDT16461885104612381","side":"buy","state":"fullyFilled","tradeId":"trader16461885203381437","ts":"1646188520338"}]}`))
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

		got, err := c.NewAssetConvertTradeService().
			QuoteId("quoterETH-USDT16461885104612381").
			BaseCcy("ETH").
			QuoteCcy("USDT").
			Side("buy").
			Sz("30").
			SzCcy("USDT").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.TradeId == "" || got.State != "fullyFilled" || got.FillQuoteSz != "30" {
			t.Fatalf("data = %#v", got)
		}
	})

	t.Run("missing_credentials", func(t *testing.T) {
		c := NewClient(
			WithNowFunc(func() time.Time { return fixedNow }),
		)
		_, err := c.NewAssetConvertTradeService().
			QuoteId("quoterETH-USDT16461885104612381").
			BaseCcy("ETH").
			QuoteCcy("USDT").
			Side("buy").
			Sz("30").
			SzCcy("USDT").
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMissingCredentials {
			t.Fatalf("error = %v, want %v", err, errMissingCredentials)
		}
	})
}
