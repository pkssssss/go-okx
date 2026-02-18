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

func TestEasyConvertService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	t.Run("signed_request_and_body", func(t *testing.T) {
		timestamp := sign.TimestampISO8601Millis(fixedNow)
		wantBody := `{"fromCcy":["ADA","USDC"],"toCcy":"OKB","source":"2"}`
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodPost, "/api/v5/trade/easy-convert", wantBody))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/trade/easy-convert"; got != want {
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
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"fillFromSz":"6.58","fillToSz":"0.17","fromCcy":"ADA","status":"running","toCcy":"OKB","uTime":"1661419684687"}]}`))
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

		got, err := c.NewEasyConvertService().
			FromCcy([]string{"ADA", "USDC"}).
			ToCcy("OKB").
			Source("2").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].FromCcy != "ADA" || got[0].UTime != 1661419684687 {
			t.Fatalf("data = %#v", got)
		}
	})

	t.Run("invalid_ack_response", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Request-Id", "rid-trade-easy-convert-invalid")
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

		_, err := c.NewEasyConvertService().
			FromCcy([]string{"ADA", "USDC"}).
			ToCcy("OKB").
			Do(context.Background())
		assertInvalidDataAPIError(t, err, errInvalidEasyConvertResponse)
		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("error = %T, want *APIError", err)
		}
		if got, want := apiErr.RequestID, "rid-trade-easy-convert-invalid"; got != want {
			t.Fatalf("apiErr.RequestID = %q, want %q", got, want)
		}
	})

	t.Run("validate_missing_required", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewEasyConvertService().ToCcy("OKB").Do(context.Background())
		if !errors.Is(err, errEasyConvertMissingRequired) {
			t.Fatalf("expected errEasyConvertMissingRequired, got %T: %v", err, err)
		}
	})

	t.Run("validate_too_many_from_ccy", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewEasyConvertService().
			FromCcy([]string{"1", "2", "3", "4", "5", "6"}).
			ToCcy("OKB").
			Do(context.Background())
		if !errors.Is(err, errEasyConvertTooManyFromCcy) {
			t.Fatalf("expected errEasyConvertTooManyFromCcy, got %T: %v", err, err)
		}
	})

	t.Run("validate_same_currency", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewEasyConvertService().
			FromCcy([]string{"OKB"}).
			ToCcy("OKB").
			Do(context.Background())
		if !errors.Is(err, errEasyConvertSameCurrency) {
			t.Fatalf("expected errEasyConvertSameCurrency, got %T: %v", err, err)
		}
	})
}
