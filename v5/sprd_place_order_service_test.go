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
)

func TestSprdPlaceOrderService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	t.Run("signed_request_and_body", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/sprd/order"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}

			bodyBytes, _ := io.ReadAll(r.Body)
			if got, want := string(bodyBytes), `{"sprdId":"BTC-USDT_BTC-USDT-SWAP","side":"buy","ordType":"limit","px":"2.15","sz":"2"}`; got != want {
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
			if got := r.Header.Get("OK-ACCESS-SIGN"); got == "" {
				t.Fatalf("OK-ACCESS-SIGN empty")
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"clOrdId":"","ordId":"312269865356374016","tag":"","sCode":"0","sMsg":""}]}`))
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

		got, err := c.NewSprdPlaceOrderService().
			SprdId("BTC-USDT_BTC-USDT-SWAP").
			Side("buy").
			OrdType("limit").
			Px("2.15").
			Sz("2").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.OrdId != "312269865356374016" {
			t.Fatalf("OrdId = %q, want %q", got.OrdId, "312269865356374016")
		}
	})

	t.Run("missing_price_for_post_only_is_server_validated", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bodyBytes, _ := io.ReadAll(r.Body)
			if got, want := string(bodyBytes), `{"sprdId":"BTC-USDT_BTC-USDT-SWAP","side":"buy","ordType":"post_only","sz":"2"}`; got != want {
				t.Fatalf("body = %q, want %q", got, want)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"51000","msg":"invalid","data":[]}`))
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

		_, err := c.NewSprdPlaceOrderService().
			SprdId("BTC-USDT_BTC-USDT-SWAP").
			Side("buy").
			OrdType("post_only").
			Sz("2").
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("err = %T, want *APIError: %v", err, err)
		}
		if apiErr.Code != "51000" {
			t.Fatalf("apiErr.Code = %q, want %q", apiErr.Code, "51000")
		}
	})

	t.Run("item_error_sCode_includes_request_id", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("x-request-id", "rid-sprd-place-1")
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"clOrdId":"","ordId":"","tag":"","sCode":"51000","sMsg":"bad"}]}`))
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

		_, err := c.NewSprdPlaceOrderService().
			SprdId("BTC-USDT_BTC-USDT-SWAP").
			Side("buy").
			OrdType("market").
			Sz("2").
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		apiErr, ok := err.(*APIError)
		if !ok {
			t.Fatalf("error = %T, want *APIError", err)
		}
		if got, want := apiErr.RequestID, "rid-sprd-place-1"; got != want {
			t.Fatalf("RequestID = %q, want %q", got, want)
		}
	})

	t.Run("invalid_ack_response", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("x-request-id", "rid-sprd-place-invalid")
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

		_, err := c.NewSprdPlaceOrderService().
			SprdId("BTC-USDT_BTC-USDT-SWAP").
			Side("buy").
			OrdType("market").
			Sz("2").
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		apiErr, ok := err.(*APIError)
		if !ok {
			t.Fatalf("error = %T, want *APIError", err)
		}
		if got, want := apiErr.RequestID, "rid-sprd-place-invalid"; got != want {
			t.Fatalf("RequestID = %q, want %q", got, want)
		}
	})

	t.Run("empty_data_response", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("x-request-id", "rid-sprd-place-empty")
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

		_, err := c.NewSprdPlaceOrderService().
			SprdId("BTC-USDT_BTC-USDT-SWAP").
			Side("buy").
			OrdType("market").
			Sz("2").
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		apiErr, ok := err.(*APIError)
		if !ok {
			t.Fatalf("error = %T, want *APIError", err)
		}
		if got, want := apiErr.RequestID, "rid-sprd-place-empty"; got != want {
			t.Fatalf("RequestID = %q, want %q", got, want)
		}
		if got, want := apiErr.Code, "0"; got != want {
			t.Fatalf("Code = %q, want %q", got, want)
		}
		if got, want := apiErr.Message, errEmptySprdPlaceOrderResponse.Error(); got != want {
			t.Fatalf("Message = %q, want %q", got, want)
		}
	})

	t.Run("multi_ack_length_mismatch_fail_close", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("x-request-id", "rid-sprd-place-multi")
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"clOrdId":"b1","ordId":"1","tag":"","sCode":"0","sMsg":""},{"clOrdId":"b2","ordId":"2","tag":"","sCode":"0","sMsg":""}]}`))
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

		_, err := c.NewSprdPlaceOrderService().
			SprdId("BTC-USDT_BTC-USDT-SWAP").
			Side("buy").
			OrdType("market").
			Sz("2").
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		apiErr, ok := err.(*APIError)
		if !ok {
			t.Fatalf("error = %T, want *APIError", err)
		}
		if got, want := apiErr.RequestID, "rid-sprd-place-multi"; got != want {
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
			w.Header().Set("x-request-id", "rid-sprd-place-multi-fail")
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"clOrdId":"b1","ordId":"1","tag":"","sCode":"0","sMsg":""},{"clOrdId":"b2","ordId":"2","tag":"","sCode":"70001","sMsg":"Order does not exist."}]}`))
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

		_, err := c.NewSprdPlaceOrderService().
			SprdId("BTC-USDT_BTC-USDT-SWAP").
			Side("buy").
			OrdType("market").
			Sz("2").
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		apiErr, ok := err.(*APIError)
		if !ok {
			t.Fatalf("error = %T, want *APIError", err)
		}
		if got, want := apiErr.RequestID, "rid-sprd-place-multi-fail"; got != want {
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
