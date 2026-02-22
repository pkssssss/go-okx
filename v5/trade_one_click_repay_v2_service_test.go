package okx

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/pkssssss/go-okx/v5/internal/sign"
)

func TestOneClickRepayV2Service_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	t.Run("signed_request_and_body", func(t *testing.T) {
		timestamp := sign.TimestampISO8601Millis(fixedNow)
		wantBody := `{"debtCcy":"USDC","repayCcyList":["USDC","BTC"]}`
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodPost, "/api/v5/trade/one-click-repay-v2", wantBody))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/trade/one-click-repay-v2"; got != want {
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
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"debtCcy":"USDC","repayCcyList":["USDC","BTC"],"ts":"1695190491421"}]}`))
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

		got, err := c.NewOneClickRepayV2Service().
			DebtCcy("USDC").
			RepayCcyList([]string{"USDC", "BTC"}).
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.DebtCcy != "USDC" || got.TS != 1695190491421 {
			t.Fatalf("data = %#v", got)
		}
	})

	t.Run("invalid_ack_response", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Request-Id", "rid-trade-one-click-repay-v2-invalid")
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{}]}`))
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

		_, err := c.NewOneClickRepayV2Service().
			DebtCcy("USDC").
			RepayCcyList([]string{"USDC", "BTC"}).
			Do(context.Background())
		assertInvalidDataAPIError(t, err, errInvalidOneClickRepayV2Response)
		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("error = %T, want *APIError", err)
		}
		if got, want := apiErr.RequestID, "rid-trade-one-click-repay-v2-invalid"; got != want {
			t.Fatalf("apiErr.RequestID = %q, want %q", got, want)
		}
	})

	t.Run("validate_missing_required", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewOneClickRepayV2Service().DebtCcy("USDC").Do(context.Background())
		if !errors.Is(err, errOneClickRepayV2MissingRequired) {
			t.Fatalf("expected errOneClickRepayV2MissingRequired, got %T: %v", err, err)
		}
	})

	t.Run("multi_ack_length_mismatch_fail_close", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Request-Id", "rid-trade-one-click-repay-v2-multi")
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"debtCcy":"USDC","repayCcyList":["USDC","BTC"],"ts":"1695190491421"},{"debtCcy":"USDC","repayCcyList":["USDC","BTC"],"ts":"1695190491422"}]}`))
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

		_, err := c.NewOneClickRepayV2Service().
			DebtCcy("USDC").
			RepayCcyList([]string{"USDC", "BTC"}).
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		apiErr, ok := err.(*APIError)
		if !ok {
			t.Fatalf("error = %T, want *APIError", err)
		}
		if got, want := apiErr.RequestPath, "/api/v5/trade/one-click-repay-v2"; got != want {
			t.Fatalf("RequestPath = %q, want %q", got, want)
		}
		if got, want := apiErr.RequestID, "rid-trade-one-click-repay-v2-multi"; got != want {
			t.Fatalf("RequestID = %q, want %q", got, want)
		}
		if got, want := apiErr.Code, "0"; got != want {
			t.Fatalf("Code = %q, want %q", got, want)
		}
		if !strings.Contains(apiErr.Message, "expected 1 ack, got 2") {
			t.Fatalf("Message = %q, want contains %q", apiErr.Message, "expected 1 ack, got 2")
		}
	})

	t.Run("multi_ack_first_success_second_fail_fail_close", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Request-Id", "rid-trade-one-click-repay-v2-multi-fail")
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"debtCcy":"USDC","repayCcyList":["USDC","BTC"],"ts":"1695190491421"},{"debtCcy":"USDC","repayCcyList":["USDC"],"ts":"1695190491422"}]}`))
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

		_, err := c.NewOneClickRepayV2Service().
			DebtCcy("USDC").
			RepayCcyList([]string{"USDC", "BTC"}).
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		apiErr, ok := err.(*APIError)
		if !ok {
			t.Fatalf("error = %T, want *APIError", err)
		}
		if got, want := apiErr.RequestPath, "/api/v5/trade/one-click-repay-v2"; got != want {
			t.Fatalf("RequestPath = %q, want %q", got, want)
		}
		if got, want := apiErr.RequestID, "rid-trade-one-click-repay-v2-multi-fail"; got != want {
			t.Fatalf("RequestID = %q, want %q", got, want)
		}
		if got, want := apiErr.Code, "0"; got != want {
			t.Fatalf("Code = %q, want %q", got, want)
		}
		if !strings.Contains(apiErr.Message, "expected 1 ack, got 2") {
			t.Fatalf("Message = %q, want contains %q", apiErr.Message, "expected 1 ack, got 2")
		}
	})
}
