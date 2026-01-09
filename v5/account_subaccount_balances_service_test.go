package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pkssssss/go-okx/v5/internal/sign"
)

func TestAccountSubaccountBalancesService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("missing_sub_acct", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewAccountSubaccountBalancesService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errAccountSubaccountBalancesMissingSubAcct {
			t.Fatalf("error = %v, want %v", err, errAccountSubaccountBalancesMissingSubAcct)
		}
	})

	t.Run("signed_request", func(t *testing.T) {
		timestamp := sign.TimestampISO8601Millis(fixedNow)
		wantPath := "/api/v5/account/subaccount/balances?subAcct=test1"
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodGet, wantPath, ""))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/account/subaccount/balances"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "subAcct=test1"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}

			if got, want := r.Header.Get("OK-ACCESS-TIMESTAMP"), timestamp; got != want {
				t.Fatalf("OK-ACCESS-TIMESTAMP = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), wantSig; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"uTime":"1739332269934","totalEq":"101.46752000000001","adjEq":"101.46752000000001","availEq":"","details":[{"ccy":"USDT","eq":"101.5","eqUsd":"101.46752000000001","cashBal":"101.5","availBal":"101.5","availEq":"101.5","frozenBal":"0","liab":"0"}]}]}`))
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

		got, err := c.NewAccountSubaccountBalancesService().SubAcct("test1").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.UTime != 1739332269934 || got.TotalEq != "101.46752000000001" || got.AdjEq != "101.46752000000001" || len(got.Details) != 1 {
			t.Fatalf("balance = %#v", got)
		}
		if got.Details[0].Ccy != "USDT" || got.Details[0].Eq != "101.5" || got.Details[0].AvailBal != "101.5" {
			t.Fatalf("detail[0] = %#v", got.Details[0])
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

		_, err := c.NewAccountSubaccountBalancesService().SubAcct("test1").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errEmptyAccountSubaccountBalances {
			t.Fatalf("error = %v, want %v", err, errEmptyAccountSubaccountBalances)
		}
	})

	t.Run("missing_credentials", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewAccountSubaccountBalancesService().SubAcct("test1").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMissingCredentials {
			t.Fatalf("error = %v, want %v", err, errMissingCredentials)
		}
	})
}
