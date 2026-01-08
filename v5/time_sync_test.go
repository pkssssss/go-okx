package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_SyncTime(t *testing.T) {
	t0 := time.Unix(1000, 0).UTC()
	t1 := time.Unix(1002, 0).UTC()
	serverTime := time.Unix(999, 0).UTC()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, http.MethodGet; got != want {
			t.Fatalf("method = %q, want %q", got, want)
		}
		if got, want := r.URL.Path, "/api/v5/public/time"; got != want {
			t.Fatalf("path = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"ts":"999000"}]}`))
	}))
	t.Cleanup(srv.Close)

	times := []time.Time{t0, t1}
	var i int
	now := func() time.Time {
		v := times[i]
		i++
		return v
	}

	c := NewClient(
		WithBaseURL(srv.URL),
		WithHTTPClient(srv.Client()),
		WithNowFunc(now),
	)

	res, err := c.SyncTime(context.Background())
	if err != nil {
		t.Fatalf("SyncTime() error = %v", err)
	}

	if got, want := res.ServerTime, serverTime; !got.Equal(want) {
		t.Fatalf("ServerTime = %s, want %s", got.Format(time.RFC3339Nano), want.Format(time.RFC3339Nano))
	}
	if got, want := res.RoundTrip, 2*time.Second; got != want {
		t.Fatalf("RoundTrip = %s, want %s", got, want)
	}
	if got, want := res.Offset, 2*time.Second; got != want {
		t.Fatalf("Offset = %s, want %s", got, want)
	}
	if got, want := c.TimeOffset(), 2*time.Second; got != want {
		t.Fatalf("Client.TimeOffset() = %s, want %s", got, want)
	}
}
