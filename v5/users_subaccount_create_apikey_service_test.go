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

func TestUsersSubaccountCreateAPIKeyService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("missing_sub_acct", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewUsersSubaccountCreateAPIKeyService().Label("v5").Passphrase("Abcdef1!").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errUsersSubaccountCreateAPIKeyMissingSubAcct {
			t.Fatalf("error = %v, want %v", err, errUsersSubaccountCreateAPIKeyMissingSubAcct)
		}
	})

	t.Run("missing_label", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewUsersSubaccountCreateAPIKeyService().SubAcct("test-1").Passphrase("Abcdef1!").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errUsersSubaccountCreateAPIKeyMissingLabel {
			t.Fatalf("error = %v, want %v", err, errUsersSubaccountCreateAPIKeyMissingLabel)
		}
	})

	t.Run("missing_passphrase", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewUsersSubaccountCreateAPIKeyService().SubAcct("test-1").Label("v5").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errUsersSubaccountCreateAPIKeyMissingPassphrase {
			t.Fatalf("error = %v, want %v", err, errUsersSubaccountCreateAPIKeyMissingPassphrase)
		}
	})

	t.Run("signed_request_and_body", func(t *testing.T) {
		timestamp := sign.TimestampISO8601Millis(fixedNow)
		wantBody := `{"subAcct":"test-1","label":"v5","passphrase":"Abcdef1!","perm":"trade","ip":"1.1.1.1,2.2.2.2"}`
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodPost, "/api/v5/users/subaccount/apikey", wantBody))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/users/subaccount/apikey"; got != want {
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
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"subAcct":"test-1","label":"v5","apiKey":"k","secretKey":"s","passphrase":"p","perm":"read_only,trade","ip":"1.1.1.1,2.2.2.2","ts":"1597026383085"}]}`))
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

		got, err := c.NewUsersSubaccountCreateAPIKeyService().
			SubAcct("test-1").
			Label("v5").
			Passphrase("Abcdef1!").
			Perm("trade").
			IP("1.1.1.1,2.2.2.2").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.SubAcct != "test-1" || got.Label != "v5" || got.APIKey != "k" || got.SecretKey.Value() != "s" || got.Passphrase.Value() != "p" || got.IP != "1.1.1.1,2.2.2.2" {
			t.Fatalf("result = %#v", got)
		}
	})

	t.Run("missing_credentials", func(t *testing.T) {
		c := NewClient(
			WithNowFunc(func() time.Time { return fixedNow }),
		)
		_, err := c.NewUsersSubaccountCreateAPIKeyService().
			SubAcct("test-1").
			Label("v5").
			Passphrase("Abcdef1!").
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMissingCredentials {
			t.Fatalf("error = %v, want %v", err, errMissingCredentials)
		}
	})
}
