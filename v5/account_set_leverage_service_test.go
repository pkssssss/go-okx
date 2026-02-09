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

func TestAccountSetLeverageService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("missing_required", func(t *testing.T) {
		c := NewClient(
			WithNowFunc(func() time.Time { return fixedNow }),
		)
		_, err := c.NewAccountSetLeverageService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errAccountSetLeverageMissingRequired {
			t.Fatalf("error = %v, want %v", err, errAccountSetLeverageMissingRequired)
		}
	})

	t.Run("ambiguous_scope", func(t *testing.T) {
		c := NewClient(
			WithNowFunc(func() time.Time { return fixedNow }),
		)
		_, err := c.NewAccountSetLeverageService().
			InstId("BTC-USDT-SWAP").
			Ccy("BTC").
			Lever("5").
			MgnMode("cross").
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errAccountSetLeverageAmbiguousScope {
			t.Fatalf("error = %v, want %v", err, errAccountSetLeverageAmbiguousScope)
		}
	})

	t.Run("invalid_mgn_mode_for_ccy", func(t *testing.T) {
		c := NewClient(
			WithNowFunc(func() time.Time { return fixedNow }),
		)
		_, err := c.NewAccountSetLeverageService().
			Ccy("BTC").
			Lever("5").
			MgnMode("isolated").
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errAccountSetLeverageInvalidMgnModeForCcy {
			t.Fatalf("error = %v, want %v", err, errAccountSetLeverageInvalidMgnModeForCcy)
		}
	})

	t.Run("invalid_pos_side_mgn_mode", func(t *testing.T) {
		c := NewClient(
			WithNowFunc(func() time.Time { return fixedNow }),
		)
		_, err := c.NewAccountSetLeverageService().
			InstId("BTC-USDT-SWAP").
			Lever("5").
			PosSide("long").
			MgnMode("cross").
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errAccountSetLeverageInvalidPosSideMgnMode {
			t.Fatalf("error = %v, want %v", err, errAccountSetLeverageInvalidPosSideMgnMode)
		}
	})

	t.Run("signed_request_and_body", func(t *testing.T) {
		timestamp := sign.TimestampISO8601Millis(fixedNow)
		wantBody := `{"instId":"BTC-USDT-SWAP","lever":"5","mgnMode":"cross"}`
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodPost, "/api/v5/account/set-leverage", wantBody))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/account/set-leverage"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, ""; got != want {
				t.Fatalf("query = %q, want %q", got, want)
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
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"lever":"5","mgnMode":"cross","instId":"BTC-USDT-SWAP","posSide":""}]}`))
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

		got, err := c.NewAccountSetLeverageService().
			InstId("BTC-USDT-SWAP").
			Lever("5").
			MgnMode("cross").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.InstId != "BTC-USDT-SWAP" || got.Lever != "5" || got.MgnMode != "cross" {
			t.Fatalf("data = %#v", got)
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

		_, err := c.NewAccountSetLeverageService().
			InstId("BTC-USDT-SWAP").
			Lever("5").
			MgnMode("cross").
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errInvalidAccountSetLeverage {
			t.Fatalf("error = %v, want %v", err, errInvalidAccountSetLeverage)
		}
	})

	t.Run("missing_credentials", func(t *testing.T) {
		c := NewClient(
			WithNowFunc(func() time.Time { return fixedNow }),
		)
		_, err := c.NewAccountSetLeverageService().
			InstId("BTC-USDT-SWAP").
			Lever("5").
			MgnMode("cross").
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMissingCredentials {
			t.Fatalf("error = %v, want %v", err, errMissingCredentials)
		}
	})
}
