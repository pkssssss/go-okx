package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_ClientStats(t *testing.T) {
	t.Run("retry_then_success", func(t *testing.T) {
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

		_, err := c.NewPublicTimeService().Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if calls != 2 {
			t.Fatalf("calls = %d, want %d", calls, 2)
		}

		st := c.ClientStats()
		if st.RequestTotal != 1 || st.SuccessTotal != 1 || st.FailureTotal != 0 || st.RetryTotal != 1 {
			t.Fatalf("stats = %#v, want request=1 success=1 failure=0 retry=1", st)
		}
		if len(st.ErrorCodeCounts) != 0 {
			t.Fatalf("ErrorCodeCounts = %#v, want empty", st.ErrorCodeCounts)
		}
	})

	t.Run("business_error_distribution", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"50011","msg":"rate limit","data":[]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		_, err := c.NewPublicTimeService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}

		st := c.ClientStats()
		if st.RequestTotal != 1 || st.SuccessTotal != 0 || st.FailureTotal != 1 || st.RetryTotal != 0 {
			t.Fatalf("stats = %#v, want request=1 success=0 failure=1 retry=0", st)
		}
		if st.ErrorCodeCounts["50011"] != 1 {
			t.Fatalf("ErrorCodeCounts = %#v, want 50011=1", st.ErrorCodeCounts)
		}

		st.ErrorCodeCounts["50011"] = 99
		next := c.ClientStats()
		if next.ErrorCodeCounts["50011"] != 1 {
			t.Fatalf("snapshot should be copied, got %#v", next.ErrorCodeCounts)
		}
	})

	t.Run("retry_then_http_failure", func(t *testing.T) {
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
			WithRetry(RetryConfig{MaxRetries: 2}),
		)

		_, err := c.NewPublicTimeService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if calls != 3 {
			t.Fatalf("calls = %d, want %d", calls, 3)
		}

		st := c.ClientStats()
		if st.RequestTotal != 1 || st.SuccessTotal != 0 || st.FailureTotal != 1 || st.RetryTotal != 2 {
			t.Fatalf("stats = %#v, want request=1 success=0 failure=1 retry=2", st)
		}
		if st.ErrorCodeCounts["HTTP_500"] != 1 {
			t.Fatalf("ErrorCodeCounts = %#v, want HTTP_500=1", st.ErrorCodeCounts)
		}
	})
}
