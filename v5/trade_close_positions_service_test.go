package okx

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClosePositionsService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	t.Run("signed_request_and_body", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/trade/close-position"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}

			bodyBytes, _ := io.ReadAll(r.Body)
			if got, want := string(bodyBytes), `{"instId":"BTC-USDT-SWAP","mgnMode":"cross"}`; got != want {
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
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), "xqUz/WD81J9qUMnmbuUZpyahaPAOHSG+/KcfQJfc/l4="; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"clOrdId":"","instId":"BTC-USDT-SWAP","posSide":"long","tag":""}]}`))
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

		acks, err := c.NewClosePositionsService().InstId("BTC-USDT-SWAP").MgnMode("cross").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(acks) != 1 {
			t.Fatalf("acks len = %d, want %d", len(acks), 1)
		}
		if acks[0].InstId != "BTC-USDT-SWAP" {
			t.Fatalf("InstId = %q, want %q", acks[0].InstId, "BTC-USDT-SWAP")
		}
	})

	t.Run("invalid_ack_response", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Request-Id", "rid-trade-close-positions-invalid")
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

		_, err := c.NewClosePositionsService().InstId("BTC-USDT-SWAP").MgnMode("cross").Do(context.Background())
		assertInvalidDataAPIError(t, err, errInvalidClosePositionsResponse)
		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("error = %T, want *APIError", err)
		}
		if got, want := apiErr.RequestID, "rid-trade-close-positions-invalid"; got != want {
			t.Fatalf("apiErr.RequestID = %q, want %q", got, want)
		}
	})

	t.Run("validate_missing_required", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewClosePositionsService().InstId("BTC-USDT-SWAP").Do(context.Background())
		if !errors.Is(err, errClosePositionsMissingRequired) {
			t.Fatalf("expected errClosePositionsMissingRequired, got %T: %v", err, err)
		}
	})
}
