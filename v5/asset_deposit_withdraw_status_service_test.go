package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAssetDepositWithdrawStatusService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("missing_id", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewAssetDepositWithdrawStatusService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errAssetDepositWithdrawStatusMissingID {
			t.Fatalf("error = %v, want %v", err, errAssetDepositWithdrawStatusMissingID)
		}
	})

	t.Run("multiple_id", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewAssetDepositWithdrawStatusService().WdId("1").TxId("tx").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errAssetDepositWithdrawStatusMultipleID {
			t.Fatalf("error = %v, want %v", err, errAssetDepositWithdrawStatusMultipleID)
		}
	})

	t.Run("deposit_missing_fields", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewAssetDepositWithdrawStatusService().TxId("tx").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errAssetDepositWithdrawStatusDepositMissing {
			t.Fatalf("error = %v, want %v", err, errAssetDepositWithdrawStatusDepositMissing)
		}
	})

	t.Run("signed_request_deposit_by_txId", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/asset/deposit-withdraw-status"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "ccy=USDT&chain=USDT-ERC20&to=addr1&txId=tx123"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}

			if got, want := r.Header.Get("OK-ACCESS-KEY"), "mykey"; got != want {
				t.Fatalf("OK-ACCESS-KEY = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-PASSPHRASE"), "mypass"; got != want {
				t.Fatalf("OK-ACCESS-PASSPHRASE = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-TIMESTAMP"), "2020-12-08T09:08:57.715Z"; got != want {
				t.Fatalf("OK-ACCESS-TIMESTAMP = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), "Oq/H2MbeInx6YLEXwL/2P+ri6JFtD9O0HAXwNlHLl/w="; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"wdId":"","txId":"tx123","state":"phase: ok","estCompleteTime":"01/09/2023, 8:10:48 PM"}]}`))
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

		got, err := c.NewAssetDepositWithdrawStatusService().
			TxId("tx123").
			Ccy("USDT").
			To("addr1").
			Chain("USDT-ERC20").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].TxId != "tx123" {
			t.Fatalf("data = %#v, want 1 item txId=tx123", got)
		}
	})

	t.Run("signed_request_withdraw_by_wdId", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/asset/deposit-withdraw-status"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "wdId=200045249"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}

			if got, want := r.Header.Get("OK-ACCESS-KEY"), "mykey"; got != want {
				t.Fatalf("OK-ACCESS-KEY = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-PASSPHRASE"), "mypass"; got != want {
				t.Fatalf("OK-ACCESS-PASSPHRASE = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-TIMESTAMP"), "2020-12-08T09:08:57.715Z"; got != want {
				t.Fatalf("OK-ACCESS-TIMESTAMP = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), "p65oE9oXDoVkGwR9/052ZgF+wbt+jyHfnpcDQBpPdBU="; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"wdId":"200045249","txId":"","state":"phase: ok","estCompleteTime":"01/09/2023, 8:10:48 PM"}]}`))
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

		got, err := c.NewAssetDepositWithdrawStatusService().
			WdId("200045249").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].WdId != "200045249" {
			t.Fatalf("data = %#v, want 1 item wdId=200045249", got)
		}
	})

	t.Run("missing_credentials", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewAssetDepositWithdrawStatusService().
			WdId("200045249").
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMissingCredentials {
			t.Fatalf("error = %v, want %v", err, errMissingCredentials)
		}
	})
}
