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

func TestRFQCancelRFQService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	t.Run("signed_request_and_body", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/rfq/cancel-rfq"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}

			bodyBytes, _ := io.ReadAll(r.Body)
			if got, want := string(bodyBytes), `{"rfqId":"22535","clRfqId":"rfq001"}`; got != want {
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
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"rfqId":"22535","clRfqId":"rfq001","sCode":"0","sMsg":""}]}`))
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

		ack, err := c.NewRFQCancelRFQService().RfqId("22535").ClRfqId("rfq001").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if ack == nil || ack.RfqId != "22535" || ack.SCode != "0" {
			t.Fatalf("ack = %#v", ack)
		}
	})

	t.Run("item_scode_error_includes_request_id", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Request-Id", "rid-rfq-cancel-rfq")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"rfqId":"22535","clRfqId":"rfq001","sCode":"51000","sMsg":"rfq failed"}]}`))
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

		_, err := c.NewRFQCancelRFQService().RfqId("22535").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}

		apiErr, ok := err.(*APIError)
		if !ok {
			t.Fatalf("err = %T, want *APIError: %v", err, err)
		}
		if got, want := apiErr.Code, "51000"; got != want {
			t.Fatalf("apiErr.Code = %q, want %q", got, want)
		}
		if got, want := apiErr.RequestPath, "/api/v5/rfq/cancel-rfq"; got != want {
			t.Fatalf("apiErr.RequestPath = %q, want %q", got, want)
		}
		if got, want := apiErr.RequestID, "rid-rfq-cancel-rfq"; got != want {
			t.Fatalf("apiErr.RequestID = %q, want %q", got, want)
		}
	})

	t.Run("invalid_ack_response", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Request-Id", "rid-rfq-cancel-rfq-invalid")
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

		_, err := c.NewRFQCancelRFQService().RfqId("22535").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}

		apiErr, ok := err.(*APIError)
		if !ok {
			t.Fatalf("err = %T, want *APIError: %v", err, err)
		}
		if got, want := apiErr.RequestID, "rid-rfq-cancel-rfq-invalid"; got != want {
			t.Fatalf("apiErr.RequestID = %q, want %q", got, want)
		}
	})

	t.Run("multi_ack_length_mismatch_fail_close", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Request-Id", "rid-rfq-cancel-rfq-multi")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"rfqId":"22535","clRfqId":"rfq001","sCode":"0","sMsg":""},{"rfqId":"22536","clRfqId":"rfq002","sCode":"0","sMsg":""}]}`))
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

		_, err := c.NewRFQCancelRFQService().RfqId("22535").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}

		apiErr, ok := err.(*APIError)
		if !ok {
			t.Fatalf("err = %T, want *APIError: %v", err, err)
		}
		if got, want := apiErr.RequestID, "rid-rfq-cancel-rfq-multi"; got != want {
			t.Fatalf("apiErr.RequestID = %q, want %q", got, want)
		}
		if got, want := apiErr.Code, "0"; got != want {
			t.Fatalf("apiErr.Code = %q, want %q", got, want)
		}
		if !strings.Contains(apiErr.Message, "expected 1 ack, got 2") {
			t.Fatalf("apiErr.Message = %q, want contains %q", apiErr.Message, "expected 1 ack, got 2")
		}
	})

	t.Run("multi_ack_first_success_second_fail_fail_close", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Request-Id", "rid-rfq-cancel-rfq-multi-fail")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"rfqId":"22535","clRfqId":"rfq001","sCode":"0","sMsg":""},{"rfqId":"22536","clRfqId":"rfq002","sCode":"70001","sMsg":"RFQ does not exist."}]}`))
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

		_, err := c.NewRFQCancelRFQService().RfqId("22535").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}

		apiErr, ok := err.(*APIError)
		if !ok {
			t.Fatalf("err = %T, want *APIError: %v", err, err)
		}
		if got, want := apiErr.RequestID, "rid-rfq-cancel-rfq-multi-fail"; got != want {
			t.Fatalf("apiErr.RequestID = %q, want %q", got, want)
		}
		if got, want := apiErr.Code, "0"; got != want {
			t.Fatalf("apiErr.Code = %q, want %q", got, want)
		}
		if !strings.Contains(apiErr.Message, "expected 1 ack, got 2") {
			t.Fatalf("apiErr.Message = %q, want contains %q", apiErr.Message, "expected 1 ack, got 2")
		}
	})

	t.Run("missing_id", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewRFQCancelRFQService().Do(context.Background())
		if !errors.Is(err, errRFQCancelRFQMissingId) {
			t.Fatalf("expected errRFQCancelRFQMissingId, got %T: %v", err, err)
		}
	})
}
