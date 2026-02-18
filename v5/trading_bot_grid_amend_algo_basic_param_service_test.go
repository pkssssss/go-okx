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

func TestTradingBotGridAmendAlgoBasicParamService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("missing_required", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewTradingBotGridAmendAlgoBasicParamService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errTradingBotGridAmendAlgoBasicParamMissingRequired {
			t.Fatalf("error = %v, want %v", err, errTradingBotGridAmendAlgoBasicParamMissingRequired)
		}
	})

	t.Run("signed_request_and_body_envelope", func(t *testing.T) {
		timestamp := sign.TimestampISO8601Millis(fixedNow)
		wantBody := `{"algoId":"123","minPx":"1","maxPx":"10","gridNum":"50"}`
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodPost, "/api/v5/tradingBot/grid/amend-algo-basic-param", wantBody))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/tradingBot/grid/amend-algo-basic-param"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}

			bodyBytes, _ := io.ReadAll(r.Body)
			if got, want := string(bodyBytes), wantBody; got != want {
				t.Fatalf("body = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-TIMESTAMP"), timestamp; got != want {
				t.Fatalf("OK-ACCESS-TIMESTAMP = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), wantSig; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":{"algoId":"123","requiredTopupAmount":"0"}}`))
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

		got, err := c.NewTradingBotGridAmendAlgoBasicParamService().
			AlgoId("123").
			MinPx("1").
			MaxPx("10").
			GridNum("50").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.AlgoId != "123" || got.RequiredTopupAmount != "0" {
			t.Fatalf("result = %#v", got)
		}
	})

	t.Run("signed_request_and_body_raw_array_response", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[{"algoId":"123","requiredTopupAmount":"1"}]`))
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

		got, err := c.NewTradingBotGridAmendAlgoBasicParamService().
			AlgoId("123").
			MinPx("1").
			MaxPx("10").
			GridNum("50").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.AlgoId != "123" || got.RequiredTopupAmount != "1" {
			t.Fatalf("result = %#v", got)
		}
	})

	t.Run("empty_data_response_returns_api_error_with_request_id", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Request-Id", "rid-grid-amend-basic-empty")
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

		_, err := c.NewTradingBotGridAmendAlgoBasicParamService().
			AlgoId("123").
			MinPx("1").
			MaxPx("10").
			GridNum("50").
			Do(context.Background())
		assertEmptyDataAPIError(t, err, errEmptyTradingBotGridAmendAlgoBasicParamResponse)

		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("error = %T, want *APIError", err)
		}
		if got, want := apiErr.RequestID, "rid-grid-amend-basic-empty"; got != want {
			t.Fatalf("apiErr.RequestID = %q, want %q", got, want)
		}
	})

	t.Run("invalid_envelope_object_missing_algo_id", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":{}}`))
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

		_, err := c.NewTradingBotGridAmendAlgoBasicParamService().
			AlgoId("123").
			MinPx("1").
			MaxPx("10").
			GridNum("50").
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		assertInvalidDataAPIError(t, err, errInvalidTradingBotGridAmendAlgoBasicParamResponse)
	})

	t.Run("invalid_raw_status_error_response", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"status":"error"}`))
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

		_, err := c.NewTradingBotGridAmendAlgoBasicParamService().
			AlgoId("123").
			MinPx("1").
			MaxPx("10").
			GridNum("50").
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		assertInvalidDataAPIError(t, err, errInvalidTradingBotGridAmendAlgoBasicParamResponse)
	})

	t.Run("missing_credentials", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewTradingBotGridAmendAlgoBasicParamService().
			AlgoId("123").
			MinPx("1").
			MaxPx("10").
			GridNum("50").
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMissingCredentials {
			t.Fatalf("error = %v, want %v", err, errMissingCredentials)
		}
	})
}
