package okx

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestRequestGate_MaxConcurrentBlocks(t *testing.T) {
	started := make(chan struct{})
	unblock := make(chan struct{})

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/api/v5/public/time"; got != want {
			t.Fatalf("path = %q, want %q", got, want)
		}

		select {
		case <-started:
		default:
			close(started)
		}

		<-unblock
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"ts":"1"}]}`))
	}))
	t.Cleanup(srv.Close)

	c := NewClient(
		WithBaseURL(srv.URL),
		WithHTTPClient(srv.Client()),
		WithRequestGate(RequestGateConfig{MaxConcurrent: 1}),
	)

	done := make(chan error, 1)
	go func() {
		_, err := c.NewPublicTimeService().Do(context.Background())
		done <- err
	}()

	<-started

	ctx2, cancel2 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	t.Cleanup(cancel2)

	_, err := c.NewPublicTimeService().Do(ctx2)
	if err == nil || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("error = %v, want context deadline exceeded", err)
	}
	var stErr *RequestStateError
	if !errors.As(err, &stErr) {
		t.Fatalf("error = %T, want *RequestStateError", err)
	}
	if stErr.Stage != RequestStageGate || stErr.Dispatched {
		t.Fatalf("RequestStateError = %#v, want stage=gate dispatched=false", stErr)
	}

	close(unblock)

	if err := <-done; err != nil {
		t.Fatalf("first request error = %v", err)
	}
}

func TestTradeAccountRateLimitService_Do_UpdatesGate(t *testing.T) {
	var orderRequests atomic.Int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v5/trade/account-rate-limit":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"accRateLimit":"1","fillRatio":"0","mainFillRatio":"0","nextAccRateLimit":"1","ts":"1"}]}`))
		case "/api/v5/trade/order":
			orderRequests.Add(1)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)

	c := NewClient(
		WithBaseURL(srv.URL),
		WithHTTPClient(srv.Client()),
		WithCredentials(Credentials{APIKey: "k", SecretKey: "s", Passphrase: "p"}),
		WithRequestGate(RequestGateConfig{MaxConcurrent: 10}),
	)

	if _, err := c.NewTradeAccountRateLimitService().Do(context.Background()); err != nil {
		t.Fatalf("TradeAccountRateLimitService.Do() error = %v", err)
	}

	ctx1, cancel1 := context.WithTimeout(context.Background(), time.Second)
	t.Cleanup(cancel1)
	if err := c.do(ctx1, http.MethodPost, "/api/v5/trade/order", nil, nil, true, nil); err != nil {
		t.Fatalf("first trade/order error = %v", err)
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	t.Cleanup(cancel2)
	err := c.do(ctx2, http.MethodPost, "/api/v5/trade/order", nil, nil, true, nil)
	if err == nil || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("second trade/order error = %v, want context deadline exceeded", err)
	}
	var stErr *RequestStateError
	if !errors.As(err, &stErr) {
		t.Fatalf("error = %T, want *RequestStateError", err)
	}
	if stErr.Stage != RequestStageGate || stErr.Dispatched {
		t.Fatalf("RequestStateError = %#v, want stage=gate dispatched=false", stErr)
	}

	if got, want := orderRequests.Load(), int32(1); got != want {
		t.Fatalf("orderRequests = %d, want %d", got, want)
	}
}

func TestTradeAccountRateLimitService_Do_AccRateLimitZero_FailSafe(t *testing.T) {
	var orderRequests atomic.Int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v5/trade/account-rate-limit":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"accRateLimit":"0","fillRatio":"0","mainFillRatio":"0","nextAccRateLimit":"0","ts":"1"}]}`))
		case "/api/v5/trade/order":
			orderRequests.Add(1)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)

	c := NewClient(
		WithBaseURL(srv.URL),
		WithHTTPClient(srv.Client()),
		WithCredentials(Credentials{APIKey: "k", SecretKey: "s", Passphrase: "p"}),
		WithRequestGate(RequestGateConfig{MaxConcurrent: 10}),
	)

	if _, err := c.NewTradeAccountRateLimitService().Do(context.Background()); err != nil {
		t.Fatalf("TradeAccountRateLimitService.Do() error = %v", err)
	}

	ctx1, cancel1 := context.WithTimeout(context.Background(), time.Second)
	t.Cleanup(cancel1)
	if err := c.do(ctx1, http.MethodPost, "/api/v5/trade/order", nil, nil, true, nil); err != nil {
		t.Fatalf("first trade/order error = %v", err)
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	t.Cleanup(cancel2)
	err := c.do(ctx2, http.MethodPost, "/api/v5/trade/order", nil, nil, true, nil)
	if err == nil || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("second trade/order error = %v, want context deadline exceeded", err)
	}
	var stErr *RequestStateError
	if !errors.As(err, &stErr) || stErr.Stage != RequestStageGate || stErr.Dispatched {
		t.Fatalf("error = %v, want gate RequestStateError dispatched=false", err)
	}

	if got, want := orderRequests.Load(), int32(1); got != want {
		t.Fatalf("orderRequests = %d, want %d", got, want)
	}
}

func TestTradeAccountRateLimitService_Do_UpdatesGate_WSCoversTradeOps(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/api/v5/trade/account-rate-limit"; got != want {
			t.Fatalf("path = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"accRateLimit":"1","fillRatio":"0","mainFillRatio":"0","nextAccRateLimit":"1","ts":"1"}]}`))
	}))
	t.Cleanup(srv.Close)

	c := NewClient(
		WithBaseURL(srv.URL),
		WithHTTPClient(srv.Client()),
		WithCredentials(Credentials{APIKey: "k", SecretKey: "s", Passphrase: "p"}),
		WithRequestGate(RequestGateConfig{MaxConcurrent: 10}),
	)

	if _, err := c.NewTradeAccountRateLimitService().Do(context.Background()); err != nil {
		t.Fatalf("TradeAccountRateLimitService.Do() error = %v", err)
	}

	g := c.gate
	g.mu.RLock()
	_, okOrder := g.routeLimiter[routeKey{Method: requestGateMethodWS, Endpoint: wsOpGateKey(wsOpOrder)}]
	_, okCancel := g.routeLimiter[routeKey{Method: requestGateMethodWS, Endpoint: wsOpGateKey(wsOpCancelOrder)}]
	_, okAmend := g.routeLimiter[routeKey{Method: requestGateMethodWS, Endpoint: wsOpGateKey(wsOpAmendOrder)}]
	g.mu.RUnlock()

	if !okOrder || !okCancel || !okAmend {
		t.Fatalf("missing ws op limiters: order=%v cancel=%v amend=%v", okOrder, okCancel, okAmend)
	}
}

func TestRequestStateError_HTTPStage_DispatchedTrue(t *testing.T) {
	started := make(chan struct{})

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/api/v5/public/time"; got != want {
			t.Fatalf("path = %q, want %q", got, want)
		}
		select {
		case <-started:
		default:
			close(started)
		}
		time.Sleep(300 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"ts":"1"}]}`))
	}))
	t.Cleanup(srv.Close)

	c := NewClient(
		WithBaseURL(srv.URL),
		WithHTTPClient(srv.Client()),
		WithRequestGate(RequestGateConfig{MaxConcurrent: 10}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	t.Cleanup(cancel)

	done := make(chan error, 1)
	go func() {
		_, err := c.NewPublicTimeService().Do(ctx)
		done <- err
	}()

	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting server handler")
	}

	err := <-done
	if err == nil || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("error = %v, want context deadline exceeded", err)
	}
	var stErr *RequestStateError
	if !errors.As(err, &stErr) {
		t.Fatalf("error = %T, want *RequestStateError", err)
	}
	if stErr.Stage != RequestStageHTTP || !stErr.Dispatched {
		t.Fatalf("RequestStateError = %#v, want stage=http dispatched=true", stErr)
	}
}
