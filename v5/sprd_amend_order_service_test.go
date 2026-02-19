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

func TestSprdAmendOrderService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	t.Run("signed_request_and_body", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/sprd/amend-order"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}

			bodyBytes, _ := io.ReadAll(r.Body)
			if got, want := string(bodyBytes), `{"ordId":"2510789768709120","reqId":"b12344","newSz":"2"}`; got != want {
				t.Fatalf("body = %q, want %q", got, want)
			}

			if got, want := r.Header.Get("OK-ACCESS-KEY"), "mykey"; got != want {
				t.Fatalf("OK-ACCESS-KEY = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"clOrdId":"","ordId":"2510789768709120","reqId":"b12344","sCode":"0","sMsg":""}]}`))
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

		got, err := c.NewSprdAmendOrderService().
			OrdId("2510789768709120").
			ReqId("b12344").
			NewSz("2").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.OrdId != "2510789768709120" || got.ReqId != "b12344" {
			t.Fatalf("ack = %#v", got)
		}
	})
}

func TestSprdAmendOrderService_Do_RejectsBothOrdIDAndClOrdID(t *testing.T) {
	c := NewClient()
	_, err := c.NewSprdAmendOrderService().
		OrdId("2510789768709120").
		ClOrdId("c1").
		NewSz("2").
		Do(context.Background())
	if err != errSprdAmendOrderTooManyId {
		t.Fatalf("error = %v, want %v", err, errSprdAmendOrderTooManyId)
	}
}

func TestSprdAmendOrderService_Do_AckError_IncludesRequestID(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-request-id", "rid-sprd-amend-1")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"clOrdId":"","ordId":"","reqId":"","sCode":"51000","sMsg":"failed"}]}`))
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

	_, err := c.NewSprdAmendOrderService().
		OrdId("2510789768709120").
		NewSz("2").
		Do(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("error = %T, want *APIError", err)
	}
	if got, want := apiErr.RequestID, "rid-sprd-amend-1"; got != want {
		t.Fatalf("RequestID = %q, want %q", got, want)
	}
}

func TestSprdAmendOrderService_Do_InvalidAckResponse(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-request-id", "rid-sprd-amend-invalid")
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

	_, err := c.NewSprdAmendOrderService().
		OrdId("2510789768709120").
		NewSz("2").
		Do(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("error = %T, want *APIError", err)
	}
	if got, want := apiErr.RequestID, "rid-sprd-amend-invalid"; got != want {
		t.Fatalf("RequestID = %q, want %q", got, want)
	}
}

func TestSprdAmendOrderService_Do_MultiAckLengthMismatchFailClose(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-request-id", "rid-sprd-amend-multi")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"clOrdId":"","ordId":"1","reqId":"","sCode":"0","sMsg":""},{"clOrdId":"","ordId":"2","reqId":"","sCode":"0","sMsg":""}]}`))
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

	_, err := c.NewSprdAmendOrderService().
		OrdId("2510789768709120").
		NewSz("2").
		Do(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("error = %T, want *APIError", err)
	}
	if got, want := apiErr.RequestID, "rid-sprd-amend-multi"; got != want {
		t.Fatalf("RequestID = %q, want %q", got, want)
	}
	if got, want := apiErr.Code, "0"; got != want {
		t.Fatalf("Code = %q, want %q", got, want)
	}
	if !strings.Contains(apiErr.Message, "expected 1 ack, got 2") {
		t.Fatalf("Message = %q, want contains %q", apiErr.Message, "expected 1 ack, got 2")
	}
}

func TestSprdAmendOrderService_Do_MultiAckFirstSuccessSecondFailFailClose(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-request-id", "rid-sprd-amend-multi-fail")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"clOrdId":"","ordId":"1","reqId":"","sCode":"0","sMsg":""},{"clOrdId":"","ordId":"2","reqId":"","sCode":"70001","sMsg":"Order does not exist."}]}`))
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

	_, err := c.NewSprdAmendOrderService().
		OrdId("2510789768709120").
		NewSz("2").
		Do(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("error = %T, want *APIError", err)
	}
	if got, want := apiErr.RequestID, "rid-sprd-amend-multi-fail"; got != want {
		t.Fatalf("RequestID = %q, want %q", got, want)
	}
	if got, want := apiErr.Code, "0"; got != want {
		t.Fatalf("Code = %q, want %q", got, want)
	}
	if !strings.Contains(apiErr.Message, "expected 1 ack, got 2") {
		t.Fatalf("Message = %q, want contains %q", apiErr.Message, "expected 1 ack, got 2")
	}
}
