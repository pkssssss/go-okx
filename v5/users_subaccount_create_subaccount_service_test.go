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

func TestUsersSubaccountCreateSubaccountService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("missing_sub_acct", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewUsersSubaccountCreateSubaccountService().Type("1").Label("L").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errUsersSubaccountCreateSubaccountMissingSubAcct {
			t.Fatalf("error = %v, want %v", err, errUsersSubaccountCreateSubaccountMissingSubAcct)
		}
	})

	t.Run("missing_type", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewUsersSubaccountCreateSubaccountService().SubAcct("subAccount002").Label("L").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errUsersSubaccountCreateSubaccountMissingType {
			t.Fatalf("error = %v, want %v", err, errUsersSubaccountCreateSubaccountMissingType)
		}
	})

	t.Run("missing_label", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewUsersSubaccountCreateSubaccountService().SubAcct("subAccount002").Type("1").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errUsersSubaccountCreateSubaccountMissingLabel {
			t.Fatalf("error = %v, want %v", err, errUsersSubaccountCreateSubaccountMissingLabel)
		}
	})

	t.Run("signed_request_and_body", func(t *testing.T) {
		timestamp := sign.TimestampISO8601Millis(fixedNow)
		wantBody := `{"subAcct":"subAccount002","type":"1","label":"123456"}`
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodPost, "/api/v5/users/subaccount/create-subaccount", wantBody))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/users/subaccount/create-subaccount"; got != want {
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
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"label":"123456","subAcct":"subAccount002","ts":"1744875304520","uid":"698827017768230914"}]}`))
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

		got, err := c.NewUsersSubaccountCreateSubaccountService().
			SubAcct("subAccount002").
			Type("1").
			Label("123456").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.SubAcct != "subAccount002" || got.Label != "123456" || got.UID != "698827017768230914" || got.TS != "1744875304520" {
			t.Fatalf("result = %#v", got)
		}
	})

	t.Run("missing_credentials", func(t *testing.T) {
		c := NewClient(
			WithNowFunc(func() time.Time { return fixedNow }),
		)
		_, err := c.NewUsersSubaccountCreateSubaccountService().
			SubAcct("subAccount002").
			Type("1").
			Label("123456").
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMissingCredentials {
			t.Fatalf("error = %v, want %v", err, errMissingCredentials)
		}
	})
}
