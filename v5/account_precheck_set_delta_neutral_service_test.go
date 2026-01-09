package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pkssssss/go-okx/v5/internal/sign"
)

func TestAccountPrecheckSetDeltaNeutralService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("missing_stgy_type", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewAccountPrecheckSetDeltaNeutralService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errAccountPrecheckSetDeltaNeutralMissingStgyType {
			t.Fatalf("error = %v, want %v", err, errAccountPrecheckSetDeltaNeutralMissingStgyType)
		}
	})

	t.Run("signed_request", func(t *testing.T) {
		timestamp := sign.TimestampISO8601Millis(fixedNow)
		wantPath := "/api/v5/account/precheck-set-delta-neutral?stgyType=1"
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodGet, wantPath, ""))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/account/precheck-set-delta-neutral"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "stgyType=1"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}

			if got, want := r.Header.Get("OK-ACCESS-TIMESTAMP"), timestamp; got != want {
				t.Fatalf("OK-ACCESS-TIMESTAMP = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), wantSig; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"unmatchedInfoCheck":[{"posList":[],"ordList":[],"deltaLever":"","type":"spot_mode"},{"posList":["123","456"],"ordList":["789"],"deltaLever":"0.5","type":"delta_risk"}]}]}`))
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

		got, err := c.NewAccountPrecheckSetDeltaNeutralService().StgyType("1").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got.UnmatchedInfoCheck) != 2 || got.UnmatchedInfoCheck[0].Type != "spot_mode" || got.UnmatchedInfoCheck[1].DeltaLever != "0.5" {
			t.Fatalf("result = %#v", got)
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

		_, err := c.NewAccountPrecheckSetDeltaNeutralService().StgyType("1").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errEmptyAccountPrecheckSetDeltaNeutral {
			t.Fatalf("error = %v, want %v", err, errEmptyAccountPrecheckSetDeltaNeutral)
		}
	})

	t.Run("missing_credentials", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewAccountPrecheckSetDeltaNeutralService().StgyType("1").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMissingCredentials {
			t.Fatalf("error = %v, want %v", err, errMissingCredentials)
		}
	})
}
