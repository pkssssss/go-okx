package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pkssssss/go-okx/v5/internal/sign"
)

func TestAccountSwitchPrecheckService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("missing_acct_lv", func(t *testing.T) {
		c := NewClient(
			WithNowFunc(func() time.Time { return fixedNow }),
		)
		_, err := c.NewAccountSwitchPrecheckService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errAccountSwitchPrecheckMissingAcctLv {
			t.Fatalf("error = %v, want %v", err, errAccountSwitchPrecheckMissingAcctLv)
		}
	})

	t.Run("signed_request", func(t *testing.T) {
		timestamp := sign.TimestampISO8601Millis(fixedNow)
		wantQuery := "acctLv=3"
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodGet, "/api/v5/account/set-account-switch-precheck?"+wantQuery, ""))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/account/set-account-switch-precheck"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, wantQuery; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}

			if got, want := r.Header.Get("OK-ACCESS-TIMESTAMP"), timestamp; got != want {
				t.Fatalf("OK-ACCESS-TIMESTAMP = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), wantSig; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"sCode":"0","curAcctLv":"4","acctLv":"3","riskOffsetType":"","unmatchedInfoCheck":[],"posList":[],"posTierCheck":[],"mgnBf":null,"mgnAft":null}]}`))
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

		got, err := c.NewAccountSwitchPrecheckService().AcctLv("3").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.SCode != "0" || got.CurAcctLv != "4" || got.AcctLv != "3" {
			t.Fatalf("result = %#v", got)
		}
		if got.MgnBf != nil || got.MgnAft != nil {
			t.Fatalf("mgnBf/mgnAft = %#v/%#v, want nil", got.MgnBf, got.MgnAft)
		}
	})

	t.Run("missing_credentials", func(t *testing.T) {
		c := NewClient(
			WithNowFunc(func() time.Time { return fixedNow }),
		)
		_, err := c.NewAccountSwitchPrecheckService().AcctLv("3").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMissingCredentials {
			t.Fatalf("error = %v, want %v", err, errMissingCredentials)
		}
	})
}
