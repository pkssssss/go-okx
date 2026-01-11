package okx

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSupportAnnouncementTypesService_Do(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		type gotReq struct {
			method string
			path   string
		}
		reqCh := make(chan gotReq, 1)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqCh <- gotReq{method: r.Method, path: r.URL.Path}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"annType":"announcements-new-listings","annTypeDesc":"New listings"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewSupportAnnouncementTypesService().Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].AnnType != "announcements-new-listings" {
			t.Fatalf("got = %#v", got)
		}

		select {
		case r := <-reqCh:
			if got, want := r.method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.path, "/api/v5/support/announcement-types"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("timeout waiting request")
		}
	})

	t.Run("business_error_code", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"51000","msg":"something wrong","data":[]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		_, err := c.NewSupportAnnouncementTypesService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("expected *APIError via errors.As, got %T: %v", err, err)
		}
		if apiErr.Code != "51000" {
			t.Fatalf("Code = %q, want %q", apiErr.Code, "51000")
		}
	})
}
