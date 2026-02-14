package rest

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestClient_Do_RedirectBlocked_SignedRequest(t *testing.T) {
	var redirectTargetHits atomic.Int64
	redirectTarget := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		redirectTargetHits.Add(1)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`ok`))
	}))
	t.Cleanup(redirectTarget.Close)

	redirector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, redirectTarget.URL, http.StatusFound)
	}))
	t.Cleanup(redirector.Close)

	c := &Client{BaseURL: redirector.URL}

	header := make(http.Header)
	header.Set("OK-ACCESS-KEY", "k")
	header.Set("OK-ACCESS-PASSPHRASE", "p")
	header.Set("OK-ACCESS-TIMESTAMP", "t")
	header.Set("OK-ACCESS-SIGN", "s")

	_, _, _, err := c.Do(context.Background(), http.MethodGet, "/start", nil, header)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, errRedirectBlockedSignedRequest) {
		t.Fatalf("error = %v, want errors.Is(..., %v)", err, errRedirectBlockedSignedRequest)
	}
	if got := redirectTargetHits.Load(); got != 0 {
		t.Fatalf("redirect target hits = %d, want 0", got)
	}
}

func TestClient_Do_RedirectBlocked_CrossHost(t *testing.T) {
	var redirectTargetHits atomic.Int64
	redirectTarget := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		redirectTargetHits.Add(1)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`ok`))
	}))
	t.Cleanup(redirectTarget.Close)

	redirector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, redirectTarget.URL, http.StatusFound)
	}))
	t.Cleanup(redirector.Close)

	c := &Client{BaseURL: redirector.URL}

	_, _, _, err := c.Do(context.Background(), http.MethodGet, "/start", nil, nil)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, errRedirectBlockedCrossHost) {
		t.Fatalf("error = %v, want errors.Is(..., %v)", err, errRedirectBlockedCrossHost)
	}
	if got := redirectTargetHits.Load(); got != 0 {
		t.Fatalf("redirect target hits = %d, want 0", got)
	}
}

func TestClient_Do_RedirectBlocked_SignedRequest_WithCustomCheckRedirect(t *testing.T) {
	var redirectTargetHits atomic.Int64
	redirectTarget := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		redirectTargetHits.Add(1)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`ok`))
	}))
	t.Cleanup(redirectTarget.Close)

	redirector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, redirectTarget.URL, http.StatusFound)
	}))
	t.Cleanup(redirector.Close)

	var customCheckRedirectCalls atomic.Int64
	customHTTPClient := redirector.Client()
	customHTTPClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		customCheckRedirectCalls.Add(1)
		return nil
	}

	c := &Client{
		BaseURL:    redirector.URL,
		HTTPClient: customHTTPClient,
	}

	header := make(http.Header)
	header.Set("OK-ACCESS-KEY", "k")
	header.Set("OK-ACCESS-PASSPHRASE", "p")
	header.Set("OK-ACCESS-TIMESTAMP", "t")
	header.Set("OK-ACCESS-SIGN", "s")

	_, _, _, err := c.Do(context.Background(), http.MethodGet, "/start", nil, header)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, errRedirectBlockedSignedRequest) {
		t.Fatalf("error = %v, want errors.Is(..., %v)", err, errRedirectBlockedSignedRequest)
	}
	if got := redirectTargetHits.Load(); got != 0 {
		t.Fatalf("redirect target hits = %d, want 0", got)
	}
	if got := customCheckRedirectCalls.Load(); got != 0 {
		t.Fatalf("custom check redirect calls = %d, want 0", got)
	}
}

func TestClient_Do_RedirectBlocked_CrossHost_WithCustomCheckRedirect(t *testing.T) {
	var redirectTargetHits atomic.Int64
	redirectTarget := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		redirectTargetHits.Add(1)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`ok`))
	}))
	t.Cleanup(redirectTarget.Close)

	redirector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, redirectTarget.URL, http.StatusFound)
	}))
	t.Cleanup(redirector.Close)

	var customCheckRedirectCalls atomic.Int64
	customHTTPClient := redirector.Client()
	customHTTPClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		customCheckRedirectCalls.Add(1)
		return nil
	}

	c := &Client{
		BaseURL:    redirector.URL,
		HTTPClient: customHTTPClient,
	}

	_, _, _, err := c.Do(context.Background(), http.MethodGet, "/start", nil, nil)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, errRedirectBlockedCrossHost) {
		t.Fatalf("error = %v, want errors.Is(..., %v)", err, errRedirectBlockedCrossHost)
	}
	if got := redirectTargetHits.Load(); got != 0 {
		t.Fatalf("redirect target hits = %d, want 0", got)
	}
	if got := customCheckRedirectCalls.Load(); got != 0 {
		t.Fatalf("custom check redirect calls = %d, want 0", got)
	}
}

func TestClient_Do_RedirectAllowed_SameHost_WithCustomCheckRedirect(t *testing.T) {
	var customCheckRedirectCalls atomic.Int64

	redirector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/start":
			http.Redirect(w, r, "/ok", http.StatusFound)
		case "/ok":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`ok`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(redirector.Close)

	customHTTPClient := redirector.Client()
	customHTTPClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		customCheckRedirectCalls.Add(1)
		return nil
	}

	c := &Client{
		BaseURL:    redirector.URL,
		HTTPClient: customHTTPClient,
	}

	status, body, _, err := c.Do(context.Background(), http.MethodGet, "/start", nil, nil)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if got, want := status, http.StatusOK; got != want {
		t.Fatalf("status = %d, want %d", got, want)
	}
	if got, want := string(body), "ok"; got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
	if got := customCheckRedirectCalls.Load(); got == 0 {
		t.Fatalf("custom check redirect calls = %d, want > 0", got)
	}
}
