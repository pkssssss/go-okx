package okx

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAssetWithdrawalService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	t.Run("missing_required", func(t *testing.T) {
		c := NewClient(
			WithNowFunc(func() time.Time { return fixedNow }),
		)
		_, err := c.NewAssetWithdrawalService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errAssetWithdrawalMissingRequired {
			t.Fatalf("error = %v, want %v", err, errAssetWithdrawalMissingRequired)
		}
	})

	t.Run("invalid_rcvr_info", func(t *testing.T) {
		c := NewClient(
			WithNowFunc(func() time.Time { return fixedNow }),
		)
		_, err := c.NewAssetWithdrawalService().
			Ccy("BTC").
			Amt("1").
			Dest("4").
			ToAddr("17DKe3kkkkiiiiTvAKKi2vMPbm1Bz3CMKw").
			RcvrInfoJSON(json.RawMessage("{")).
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errAssetWithdrawalInvalidRcvrInfo {
			t.Fatalf("error = %v, want %v", err, errAssetWithdrawalInvalidRcvrInfo)
		}
	})

	t.Run("signed_request_and_body", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/asset/withdrawal"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, ""; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}

			bodyBytes, _ := io.ReadAll(r.Body)
			if got, want := string(bodyBytes), `{"ccy":"BTC","amt":"1","dest":"4","toAddr":"17DKe3kkkkiiiiTvAKKi2vMPbm1Bz3CMKw","chain":"BTC-Bitcoin"}`; got != want {
				t.Fatalf("body = %q, want %q", got, want)
			}

			if got, want := r.Header.Get("OK-ACCESS-KEY"), "mykey"; got != want {
				t.Fatalf("OK-ACCESS-KEY = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-PASSPHRASE"), "mypass"; got != want {
				t.Fatalf("OK-ACCESS-PASSPHRASE = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-TIMESTAMP"), "2020-03-28T12:21:41.274Z"; got != want {
				t.Fatalf("OK-ACCESS-TIMESTAMP = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), "gBO3BmS1cPglcPTHdPjfjtMOoqoqtL74rj8uhfpJJzE="; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"amt":"1","wdId":"67485","ccy":"BTC","clientId":"","chain":"BTC-Bitcoin"}]}`))
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

		got, err := c.NewAssetWithdrawalService().
			Ccy("BTC").
			Amt("1").
			Dest("4").
			ToAddr("17DKe3kkkkiiiiTvAKKi2vMPbm1Bz3CMKw").
			Chain("BTC-Bitcoin").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.WdId != "67485" {
			t.Fatalf("WdId = %q, want %q", got.WdId, "67485")
		}
	})

	t.Run("missing_credentials", func(t *testing.T) {
		c := NewClient(
			WithNowFunc(func() time.Time { return fixedNow }),
		)
		_, err := c.NewAssetWithdrawalService().
			Ccy("BTC").
			Amt("1").
			Dest("4").
			ToAddr("17DKe3kkkkiiiiTvAKKi2vMPbm1Bz3CMKw").
			Chain("BTC-Bitcoin").
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMissingCredentials {
			t.Fatalf("error = %v, want %v", err, errMissingCredentials)
		}
	})
}
