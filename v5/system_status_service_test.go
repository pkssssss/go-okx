package okx

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSystemStatusService_Do(t *testing.T) {
	t.Run("ok_with_state", func(t *testing.T) {
		type gotReq struct {
			method string
			path   string
			query  string
		}
		reqCh := make(chan gotReq, 1)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqCh <- gotReq{method: r.Method, path: r.URL.Path, query: r.URL.RawQuery}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"begin":"1672823400000","end":"1672823520000","href":"","preOpenBegin":"","scheDesc":"","serviceType":"8","state":"completed","maintType":"1","env":"1","system":"unified","title":"Trading account system upgrade"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewSystemStatusService().State("canceled").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].Title == "" || got[0].Begin != 1672823400000 {
			t.Fatalf("got = %#v", got)
		}

		select {
		case r := <-reqCh:
			if got, want := r.method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.path, "/api/v5/system/status"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.query, "state=canceled"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
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

		_, err := c.NewSystemStatusService().Do(context.Background())
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
