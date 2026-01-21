package okx

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMarketBooksSBEService_Do(t *testing.T) {
	t.Run("missing_instIdCode", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewMarketBooksSBEService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMarketBooksSBEMissingInstIdCode {
			t.Fatalf("error = %v, want %v", err, errMarketBooksSBEMissingInstIdCode)
		}
	})

	t.Run("ok", func(t *testing.T) {
		wantBody := []byte{0x01, 0x02, 0x03}
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/market/books-sbe"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "instIdCode=12345&source=0"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("Accept"); got == "" {
				t.Fatalf("missing header Accept")
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/sbe")
			_, _ = w.Write(wantBody)
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewMarketBooksSBEService().InstIdCode(12345).Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if string(got) != string(wantBody) {
			t.Fatalf("body = %v, want %v", got, wantBody)
		}
	})

	t.Run("json_error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"51000","msg":"Parameter instIdCode error","data":[]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		_, err := c.NewMarketBooksSBEService().InstIdCode(0).Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("expected *APIError, got %T: %v", err, err)
		}
		if apiErr.Code != "51000" {
			t.Fatalf("Code = %q, want %q", apiErr.Code, "51000")
		}
	})

	t.Run("respects_request_gate", func(t *testing.T) {
		started := make(chan struct{})
		unblock := make(chan struct{})

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			select {
			case <-started:
			default:
				close(started)
			}

			<-unblock
			w.Header().Set("Content-Type", "application/sbe")
			_, _ = w.Write([]byte{0x01})
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
			WithRequestGate(RequestGateConfig{MaxConcurrent: 1}),
		)

		done := make(chan error, 1)
		go func() {
			_, err := c.NewMarketBooksSBEService().InstIdCode(12345).Do(context.Background())
			done <- err
		}()

		select {
		case <-started:
		case <-time.After(2 * time.Second):
			t.Fatalf("timeout waiting first request")
		}

		ctx2, cancel2 := context.WithTimeout(context.Background(), 100*time.Millisecond)
		t.Cleanup(cancel2)

		_, err := c.NewMarketBooksSBEService().InstIdCode(12345).Do(ctx2)
		var stErr *RequestStateError
		if err == nil || !errors.As(err, &stErr) || stErr.Stage != RequestStageGate || stErr.Dispatched {
			t.Fatalf("error = %v, want gate RequestStateError dispatched=false", err)
		}

		close(unblock)
		if err := <-done; err != nil {
			t.Fatalf("first request error = %v", err)
		}
	})
}
