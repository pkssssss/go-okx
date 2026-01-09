package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientRetry_GET_ServerError(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("Content-Type", "application/json")
		if calls == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[]}`))
			return
		}
		_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"ts":"1"}]}`))
	}))
	t.Cleanup(srv.Close)

	c := NewClient(
		WithBaseURL(srv.URL),
		WithHTTPClient(srv.Client()),
		WithRetry(RetryConfig{MaxRetries: 1}),
	)

	st, err := c.NewPublicTimeService().Do(context.Background())
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if st.TS != 1 {
		t.Fatalf("TS = %d, want %d", st.TS, 1)
	}
	if calls != 2 {
		t.Fatalf("calls = %d, want %d", calls, 2)
	}
}

func TestClientRetry_GET_RateLimit(t *testing.T) {
	t.Run("not_enabled_by_default", func(t *testing.T) {
		calls := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls++
			w.Header().Set("Content-Type", "application/json")
			if calls == 1 {
				_, _ = w.Write([]byte(`{"code":"50011","msg":"rate limit","data":[]}`))
				return
			}
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"ts":"1"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
			WithRetry(RetryConfig{MaxRetries: 1}),
		)

		_, err := c.NewPublicTimeService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if calls != 1 {
			t.Fatalf("calls = %d, want %d", calls, 1)
		}
	})

	t.Run("enabled", func(t *testing.T) {
		calls := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls++
			w.Header().Set("Content-Type", "application/json")
			if calls == 1 {
				_, _ = w.Write([]byte(`{"code":"50011","msg":"rate limit","data":[]}`))
				return
			}
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"ts":"1"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
			WithRetry(RetryConfig{MaxRetries: 1, RetryOnRateLimit: true}),
		)

		st, err := c.NewPublicTimeService().Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if st.TS != 1 {
			t.Fatalf("TS = %d, want %d", st.TS, 1)
		}
		if calls != 2 {
			t.Fatalf("calls = %d, want %d", calls, 2)
		}
	})
}

func TestClientRetry_DoesNotRetry_NonGET(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[]}`))
	}))
	t.Cleanup(srv.Close)

	c := NewClient(
		WithBaseURL(srv.URL),
		WithHTTPClient(srv.Client()),
		WithRetry(RetryConfig{MaxRetries: 1}),
	)

	err := c.doWithHeaders(context.Background(), http.MethodPost, "/api/v5/test", nil, map[string]string{"k": "v"}, false, nil, nil)
	if err == nil {
		t.Fatalf("expected error")
	}
	if calls != 1 {
		t.Fatalf("calls = %d, want %d", calls, 1)
	}
}
