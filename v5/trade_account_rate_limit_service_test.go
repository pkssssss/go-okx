package okx

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTradeAccountRateLimitService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, http.MethodGet; got != want {
			t.Fatalf("method = %q, want %q", got, want)
		}
		if got, want := r.URL.Path, "/api/v5/trade/account-rate-limit"; got != want {
			t.Fatalf("path = %q, want %q", got, want)
		}

		bodyBytes, _ := io.ReadAll(r.Body)
		if got, want := string(bodyBytes), ``; got != want {
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
		if got, want := r.Header.Get("OK-ACCESS-SIGN"), "DJTOH4QY48mBC5BBeXmWeEy8GOVFNAdCU2xRsHJhE1I="; got != want {
			t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"accRateLimit":"60","fillRatio":"0.01","mainFillRatio":"0.02","nextAccRateLimit":"60","ts":"1695190491421"}]}`))
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

	got, err := c.NewTradeAccountRateLimitService().Do(context.Background())
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if got.AccRateLimit != "60" {
		t.Fatalf("AccRateLimit = %q, want %q", got.AccRateLimit, "60")
	}
}
