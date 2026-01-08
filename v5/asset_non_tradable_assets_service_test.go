package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pkssssss/go-okx/v5/internal/sign"
)

func TestAssetNonTradableAssetsService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("signed_request", func(t *testing.T) {
		timestamp := sign.TimestampISO8601Millis(fixedNow)
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodGet, "/api/v5/asset/non-tradable-assets?ccy=CELT%2CMEME", ""))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/asset/non-tradable-assets"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "ccy=CELT%2CMEME"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}

			if got, want := r.Header.Get("OK-ACCESS-TIMESTAMP"), timestamp; got != want {
				t.Fatalf("OK-ACCESS-TIMESTAMP = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), wantSig; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"ccy":"CELT","name":"","logoLink":"https://static.example/celt.png","chain":"CELT-OKTC","ctAddr":"","bal":"989.84719571","canWd":true,"minWd":"0.1","fee":"2","feeCcy":"USDT","wdTickSz":"8","wdAll":false,"needTag":false},{"ccy":"MEME","name":"MEME Inu","logoLink":"https://static.example/meme.png","chain":"MEME-ERC20","ctAddr":"09b760","bal":"0.001","canWd":true,"minWd":"0.001","fee":"5","feeCcy":"USDT","wdTickSz":"8","wdAll":false,"needTag":false}]}`))
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

		got, err := c.NewAssetNonTradableAssetsService().Ccy("CELT,MEME").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 2 {
			t.Fatalf("len(data) = %d, want %d", len(got), 2)
		}
		if got[0].Ccy != "CELT" || got[0].Bal != "989.84719571" || got[0].FeeCcy != "USDT" {
			t.Fatalf("data[0] = %#v", got[0])
		}
	})

	t.Run("missing_credentials", func(t *testing.T) {
		c := NewClient(
			WithNowFunc(func() time.Time { return fixedNow }),
		)
		_, err := c.NewAssetNonTradableAssetsService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMissingCredentials {
			t.Fatalf("error = %v, want %v", err, errMissingCredentials)
		}
	})
}
