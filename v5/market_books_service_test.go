package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMarketBooksService_Do_EmptyResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, http.MethodGet; got != want {
			t.Fatalf("method = %q, want %q", got, want)
		}
		if got, want := r.URL.Path, "/api/v5/market/books"; got != want {
			t.Fatalf("path = %q, want %q", got, want)
		}
		if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
			t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[]}`))
	}))
	t.Cleanup(srv.Close)

	c := NewClient(
		WithBaseURL(srv.URL),
		WithHTTPClient(srv.Client()),
	)

	_, err := c.NewMarketBooksService().InstId("BTC-USDT").Do(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
	if err != errEmptyMarketBooksResponse {
		t.Fatalf("error = %v, want %v", err, errEmptyMarketBooksResponse)
	}
}
