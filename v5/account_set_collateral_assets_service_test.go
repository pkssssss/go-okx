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

func TestAccountSetCollateralAssetsService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("missing_required", func(t *testing.T) {
		c := NewClient(
			WithNowFunc(func() time.Time { return fixedNow }),
		)
		_, err := c.NewAccountSetCollateralAssetsService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errAccountSetCollateralAssetsMissingRequired {
			t.Fatalf("error = %v, want %v", err, errAccountSetCollateralAssetsMissingRequired)
		}
	})

	t.Run("custom_missing_ccy_list", func(t *testing.T) {
		c := NewClient(
			WithNowFunc(func() time.Time { return fixedNow }),
		)
		_, err := c.NewAccountSetCollateralAssetsService().Type("custom").CollateralEnabled(true).Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errAccountSetCollateralAssetsMissingCcyList {
			t.Fatalf("error = %v, want %v", err, errAccountSetCollateralAssetsMissingCcyList)
		}
	})

	t.Run("signed_request_all", func(t *testing.T) {
		timestamp := sign.TimestampISO8601Millis(fixedNow)
		wantBody := `{"type":"all","collateralEnabled":true}`
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodPost, "/api/v5/account/set-collateral-assets", wantBody))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/account/set-collateral-assets"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}

			bodyBytes, _ := io.ReadAll(r.Body)
			if got, want := string(bodyBytes), wantBody; got != want {
				t.Fatalf("body = %q, want %q", got, want)
			}

			if got, want := r.Header.Get("OK-ACCESS-TIMESTAMP"), timestamp; got != want {
				t.Fatalf("OK-ACCESS-TIMESTAMP = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), wantSig; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"type":"all","ccyList":["BTC","ETH"],"collateralEnabled":false}]}`))
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

		got, err := c.NewAccountSetCollateralAssetsService().Type("all").CollateralEnabled(true).Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.Type != "all" || got.CollateralEnabled != false || len(got.CcyList) != 2 || got.CcyList[0] != "BTC" {
			t.Fatalf("ack = %#v", got)
		}
	})

	t.Run("signed_request_custom", func(t *testing.T) {
		timestamp := sign.TimestampISO8601Millis(fixedNow)
		wantBody := `{"type":"custom","ccyList":["BTC","ETH"],"collateralEnabled":false}`
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodPost, "/api/v5/account/set-collateral-assets", wantBody))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bodyBytes, _ := io.ReadAll(r.Body)
			if got, want := string(bodyBytes), wantBody; got != want {
				t.Fatalf("body = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), wantSig; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"type":"custom","ccyList":["BTC","ETH"],"collateralEnabled":false}]}`))
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

		got, err := c.NewAccountSetCollateralAssetsService().
			Type("custom").
			CcyList([]string{"BTC", "ETH"}).
			CollateralEnabled(false).
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.Type != "custom" || got.CollateralEnabled != false || len(got.CcyList) != 2 || got.CcyList[1] != "ETH" {
			t.Fatalf("ack = %#v", got)
		}
	})

	t.Run("invalid_ack_response", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{}]}`))
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

		_, err := c.NewAccountSetCollateralAssetsService().Type("all").CollateralEnabled(true).Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errInvalidAccountSetCollateralAssets {
			t.Fatalf("error = %v, want %v", err, errInvalidAccountSetCollateralAssets)
		}
	})

	t.Run("missing_credentials", func(t *testing.T) {
		c := NewClient(
			WithNowFunc(func() time.Time { return fixedNow }),
		)
		_, err := c.NewAccountSetCollateralAssetsService().Type("all").CollateralEnabled(true).Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMissingCredentials {
			t.Fatalf("error = %v, want %v", err, errMissingCredentials)
		}
	})
}
