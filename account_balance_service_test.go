package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAccountBalanceService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("signed_request", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/account/balance"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "ccy=BTC"; got != want {
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
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), "OkAjRTXqfIRKxx7SRjIowU96vZPkf4n9X2G+8yduCf4="; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"uTime":"1597026383085","totalEq":"1","adjEq":"1","details":[{"ccy":"BTC","eq":"1","cashBal":"1"}]}]}`))
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

		got, err := c.NewAccountBalanceService().Ccy("BTC").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.TotalEq != "1" {
			t.Fatalf("TotalEq = %q, want %q", got.TotalEq, "1")
		}
		if len(got.Details) != 1 || got.Details[0].Ccy != "BTC" {
			t.Fatalf("Details = %#v, want 1 item BTC", got.Details)
		}
	})

	t.Run("demo_header", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Header.Get("x-simulated-trading"), "1"; got != want {
				t.Fatalf("x-simulated-trading = %q, want %q", got, want)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"uTime":"1597026383085","totalEq":"0","adjEq":"0","details":[]}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
			WithDemoTrading(true),
			WithCredentials(Credentials{
				APIKey:     "mykey",
				SecretKey:  "mysecret",
				Passphrase: "mypass",
			}),
			WithNowFunc(func() time.Time { return fixedNow }),
		)

		_, err := c.NewAccountBalanceService().Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
	})

	t.Run("missing_credentials", func(t *testing.T) {
		c := NewClient(
			WithNowFunc(func() time.Time { return fixedNow }),
		)
		_, err := c.NewAccountBalanceService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMissingCredentials {
			t.Fatalf("error = %v, want %v", err, errMissingCredentials)
		}
	})
}
