package okx

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pkssssss/go-okx/v5/internal/sign"
)

func TestOneClickRepayService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	t.Run("signed_request_and_body", func(t *testing.T) {
		timestamp := sign.TimestampISO8601Millis(fixedNow)
		wantBody := `{"debtCcy":["ETH","BTC"],"repayCcy":"USDT"}`
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodPost, "/api/v5/trade/one-click-repay", wantBody))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/trade/one-click-repay"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}

			bodyBytes, _ := io.ReadAll(r.Body)
			if got, want := string(bodyBytes), wantBody; got != want {
				t.Fatalf("body = %q, want %q", got, want)
			}

			if got, want := r.Header.Get("OK-ACCESS-KEY"), "mykey"; got != want {
				t.Fatalf("OK-ACCESS-KEY = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-PASSPHRASE"), "mypass"; got != want {
				t.Fatalf("OK-ACCESS-PASSPHRASE = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-TIMESTAMP"), timestamp; got != want {
				t.Fatalf("OK-ACCESS-TIMESTAMP = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), wantSig; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"debtCcy":"ETH","fillDebtSz":"0.0102","fillRepaySz":"30","repayCcy":"USDT","status":"filled","uTime":"1646188520338"}]}`))
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

		got, err := c.NewOneClickRepayService().
			DebtCcy([]string{"ETH", "BTC"}).
			RepayCcy("USDT").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].DebtCcy != "ETH" || got[0].UTime != 1646188520338 {
			t.Fatalf("data = %#v", got)
		}
	})

	t.Run("validate_missing_required", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewOneClickRepayService().RepayCcy("USDT").Do(context.Background())
		if !errors.Is(err, errOneClickRepayMissingRequired) {
			t.Fatalf("expected errOneClickRepayMissingRequired, got %T: %v", err, err)
		}
	})

	t.Run("validate_too_many_debt_ccy", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewOneClickRepayService().
			DebtCcy([]string{"1", "2", "3", "4", "5", "6"}).
			RepayCcy("USDT").
			Do(context.Background())
		if !errors.Is(err, errOneClickRepayTooManyDebtCcy) {
			t.Fatalf("expected errOneClickRepayTooManyDebtCcy, got %T: %v", err, err)
		}
	})

	t.Run("validate_same_currency", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewOneClickRepayService().
			DebtCcy([]string{"USDT"}).
			RepayCcy("USDT").
			Do(context.Background())
		if !errors.Is(err, errOneClickRepaySameCurrency) {
			t.Fatalf("expected errOneClickRepaySameCurrency, got %T: %v", err, err)
		}
	})
}
