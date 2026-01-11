package okx

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSupportAnnouncementsService_Do(t *testing.T) {
	t.Run("ok_public", func(t *testing.T) {
		type gotReq struct {
			method         string
			path           string
			query          string
			acceptLanguage string
			accessKey      string
		}
		reqCh := make(chan gotReq, 1)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqCh <- gotReq{
				method:         r.Method,
				path:           r.URL.Path,
				query:          r.URL.RawQuery,
				acceptLanguage: r.Header.Get("Accept-Language"),
				accessKey:      r.Header.Get("OK-ACCESS-KEY"),
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"details":[{"annType":"announcements-new-listings","title":"t1","url":"https://example.com","pTime":"1761620404821","businessPTime":"1761620400000"}],"totalPage":"123"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewSupportAnnouncementsService().
			AnnType("announcements-new-listings").
			Page("2").
			AcceptLanguage("zh-CN").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.TotalPage != "123" || len(got.Details) != 1 || got.Details[0].PTime != 1761620404821 {
			t.Fatalf("got = %#v", got)
		}

		select {
		case r := <-reqCh:
			if got, want := r.method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.path, "/api/v5/support/announcements"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.query, "annType=announcements-new-listings&page=2"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got, want := r.acceptLanguage, "zh-CN"; got != want {
				t.Fatalf("Accept-Language = %q, want %q", got, want)
			}
			if r.accessKey != "" {
				t.Fatalf("expected unsigned request, got OK-ACCESS-KEY=%q", r.accessKey)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("timeout waiting request")
		}
	})

	t.Run("ok_signed", func(t *testing.T) {
		type gotReq struct {
			accessKey string
			sign      string
		}
		reqCh := make(chan gotReq, 1)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqCh <- gotReq{
				accessKey: r.Header.Get("OK-ACCESS-KEY"),
				sign:      r.Header.Get("OK-ACCESS-SIGN"),
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"details":[],"totalPage":"0"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
			WithCredentials(Credentials{
				APIKey:     "k",
				SecretKey:  "s",
				Passphrase: "p",
			}),
		)

		_, err := c.NewSupportAnnouncementsService().Signed(true).Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}

		select {
		case r := <-reqCh:
			if got, want := r.accessKey, "k"; got != want {
				t.Fatalf("OK-ACCESS-KEY = %q, want %q", got, want)
			}
			if r.sign == "" {
				t.Fatalf("expected OK-ACCESS-SIGN")
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

		_, err := c.NewSupportAnnouncementsService().Do(context.Background())
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
