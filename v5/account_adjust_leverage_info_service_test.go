package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pkssssss/go-okx/v5/internal/sign"
)

func TestAccountAdjustLeverageInfoService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("missing_required", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewAccountAdjustLeverageInfoService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errAccountAdjustLeverageInfoMissingRequired {
			t.Fatalf("error = %v, want %v", err, errAccountAdjustLeverageInfoMissingRequired)
		}
	})

	t.Run("signed_request", func(t *testing.T) {
		timestamp := sign.TimestampISO8601Millis(fixedNow)
		wantPath := "/api/v5/account/adjust-leverage-info?instId=BTC-USDT&instType=MARGIN&lever=3&mgnMode=isolated"
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodGet, wantPath, ""))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/account/adjust-leverage-info"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "instId=BTC-USDT&instType=MARGIN&lever=3&mgnMode=isolated"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}

			if got, want := r.Header.Get("OK-ACCESS-TIMESTAMP"), timestamp; got != want {
				t.Fatalf("OK-ACCESS-TIMESTAMP = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), wantSig; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"estAvailQuoteTrans":"","estAvailTrans":"1.1398040558348279","estLiqPx":"","estMaxAmt":"10.6095865868904898","estMgn":"0.0701959441651721","estQuoteMaxAmt":"176889.6871254563042714","estQuoteMgn":"","existOrd":false,"maxLever":"10","minLever":"0.01"}]}`))
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

		got, err := c.NewAccountAdjustLeverageInfoService().
			InstType("MARGIN").
			MgnMode("isolated").
			Lever("3").
			InstId("BTC-USDT").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.MaxLever != "10" || got.MinLever != "0.01" || got.EstMgn != "0.0701959441651721" || got.ExistOrd {
			t.Fatalf("info = %#v", got)
		}
	})

	t.Run("empty_response", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[]}`))
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

		_, err := c.NewAccountAdjustLeverageInfoService().
			InstType("MARGIN").
			MgnMode("isolated").
			Lever("3").
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errEmptyAccountAdjustLeverageInfo {
			t.Fatalf("error = %v, want %v", err, errEmptyAccountAdjustLeverageInfo)
		}
	})

	t.Run("missing_credentials", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewAccountAdjustLeverageInfoService().
			InstType("MARGIN").
			MgnMode("isolated").
			Lever("3").
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMissingCredentials {
			t.Fatalf("error = %v, want %v", err, errMissingCredentials)
		}
	})
}
