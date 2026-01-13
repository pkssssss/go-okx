package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pkssssss/go-okx/v5/internal/sign"
)

func TestAffiliateInviteeDetailService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("missing_uid", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewAffiliateInviteeDetailService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errAffiliateInviteeDetailMissingUID {
			t.Fatalf("error = %v, want %v", err, errAffiliateInviteeDetailMissingUID)
		}
	})

	t.Run("signed_request", func(t *testing.T) {
		ts := sign.TimestampISO8601Millis(fixedNow)
		requestPath := "/api/v5/affiliate/invitee/detail?uid=11111111"
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(ts, http.MethodGet, requestPath, ""))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/affiliate/invitee/detail"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "uid=11111111"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}

			if got, want := r.Header.Get("OK-ACCESS-SIGN"), wantSig; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"accFee":"0","affiliateCode":"HIIIIII","depAmt":"0","firstTradeTime":"","inviteeLevel":"2","inviteeRebateRate":"0.39","joinTime":"1712546713000","kycTime":"","level":"Lv1","region":"越南","totalCommission":"0","volMonth":"0"}]}`))
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

		got, err := c.NewAffiliateInviteeDetailService().UID("11111111").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.AffiliateCode != "HIIIIII" {
			t.Fatalf("AffiliateCode = %q, want %q", got.AffiliateCode, "HIIIIII")
		}
		if got.JoinTime != 1712546713000 {
			t.Fatalf("JoinTime = %d, want %d", got.JoinTime, 1712546713000)
		}
		if got.FirstTradeTime != 0 {
			t.Fatalf("FirstTradeTime = %d, want %d", got.FirstTradeTime, 0)
		}
	})

	t.Run("empty_data", func(t *testing.T) {
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

		_, err := c.NewAffiliateInviteeDetailService().UID("11111111").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errEmptyAffiliateInviteeDetail {
			t.Fatalf("error = %v, want %v", err, errEmptyAffiliateInviteeDetail)
		}
	})
}
