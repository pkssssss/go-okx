package okx

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSprdCancelOrderService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	t.Run("signed_request_and_body_ordId", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/sprd/cancel-order"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}

			bodyBytes, _ := io.ReadAll(r.Body)
			if got, want := string(bodyBytes), `{"ordId":"2510789768709120"}`; got != want {
				t.Fatalf("body = %q, want %q", got, want)
			}

			if got, want := r.Header.Get("OK-ACCESS-KEY"), "mykey"; got != want {
				t.Fatalf("OK-ACCESS-KEY = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"clOrdId":"b1","ordId":"2510789768709120","sCode":"0","sMsg":""}]}`))
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

		got, err := c.NewSprdCancelOrderService().
			OrdId("2510789768709120").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.OrdId != "2510789768709120" {
			t.Fatalf("OrdId = %q, want %q", got.OrdId, "2510789768709120")
		}
	})
}

func TestSprdCancelOrderService_Do_RejectsBothOrdIDAndClOrdID(t *testing.T) {
	c := NewClient()
	_, err := c.NewSprdCancelOrderService().
		OrdId("2510789768709120").
		ClOrdId("c1").
		Do(context.Background())
	if err != errSprdCancelOrderTooManyId {
		t.Fatalf("error = %v, want %v", err, errSprdCancelOrderTooManyId)
	}
}

func TestSprdCancelOrderService_Do_AckError_IncludesRequestID(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-request-id", "rid-sprd-cancel-1")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"clOrdId":"","ordId":"","sCode":"51000","sMsg":"failed"}]}`))
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

	_, err := c.NewSprdCancelOrderService().
		OrdId("2510789768709120").
		Do(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("error = %T, want *APIError", err)
	}
	if got, want := apiErr.RequestID, "rid-sprd-cancel-1"; got != want {
		t.Fatalf("RequestID = %q, want %q", got, want)
	}
}

func TestSprdCancelOrderService_Do_InvalidAckResponse(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-request-id", "rid-sprd-cancel-invalid")
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

	_, err := c.NewSprdCancelOrderService().
		OrdId("2510789768709120").
		Do(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("error = %T, want *APIError", err)
	}
	if got, want := apiErr.RequestID, "rid-sprd-cancel-invalid"; got != want {
		t.Fatalf("RequestID = %q, want %q", got, want)
	}
}

func TestSprdCancelOrderService_Do_MultiAckLengthMismatchFailClose(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-request-id", "rid-sprd-cancel-multi")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"clOrdId":"b1","ordId":"1","sCode":"0","sMsg":""},{"clOrdId":"b2","ordId":"2","sCode":"0","sMsg":""}]}`))
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

	_, err := c.NewSprdCancelOrderService().
		OrdId("2510789768709120").
		Do(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("error = %T, want *APIError", err)
	}
	if got, want := apiErr.RequestID, "rid-sprd-cancel-multi"; got != want {
		t.Fatalf("RequestID = %q, want %q", got, want)
	}
	if got, want := apiErr.Code, "0"; got != want {
		t.Fatalf("Code = %q, want %q", got, want)
	}
	if !strings.Contains(apiErr.Message, "expected 1 ack, got 2") {
		t.Fatalf("Message = %q, want contains %q", apiErr.Message, "expected 1 ack, got 2")
	}
}

func TestSprdCancelOrderService_Do_MultiAckFirstSuccessSecondFailFailClose(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-request-id", "rid-sprd-cancel-multi-fail")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"clOrdId":"b1","ordId":"1","sCode":"0","sMsg":""},{"clOrdId":"b2","ordId":"2","sCode":"70001","sMsg":"Order does not exist."}]}`))
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

	_, err := c.NewSprdCancelOrderService().
		OrdId("2510789768709120").
		Do(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("error = %T, want *APIError", err)
	}
	if got, want := apiErr.RequestID, "rid-sprd-cancel-multi-fail"; got != want {
		t.Fatalf("RequestID = %q, want %q", got, want)
	}
	if got, want := apiErr.Code, "0"; got != want {
		t.Fatalf("Code = %q, want %q", got, want)
	}
	if !strings.Contains(apiErr.Message, "expected 1 ack, got 2") {
		t.Fatalf("Message = %q, want contains %q", apiErr.Message, "expected 1 ack, got 2")
	}
}
