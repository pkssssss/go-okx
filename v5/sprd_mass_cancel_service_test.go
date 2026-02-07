package okx

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSprdMassCancelService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, http.MethodPost; got != want {
			t.Fatalf("method = %q, want %q", got, want)
		}
		if got, want := r.URL.Path, "/api/v5/sprd/mass-cancel"; got != want {
			t.Fatalf("path = %q, want %q", got, want)
		}

		bodyBytes, _ := io.ReadAll(r.Body)
		if got, want := string(bodyBytes), `{"sprdId":"BTC-USDT_BTC-USDT-SWAP"}`; got != want {
			t.Fatalf("body = %q, want %q", got, want)
		}

		if got, want := r.Header.Get("OK-ACCESS-KEY"), "mykey"; got != want {
			t.Fatalf("OK-ACCESS-KEY = %q, want %q", got, want)
		}
		if got, want := r.Header.Get("OK-ACCESS-PASSPHRASE"), "mypass"; got != want {
			t.Fatalf("OK-ACCESS-PASSPHRASE = %q, want %q", got, want)
		}
		if got, want := r.Header.Get("OK-ACCESS-TIMESTAMP"), "2020-03-28T12:21:41.274Z"; got != want {
			t.Fatalf("OK-ACCESS-TIMESTAMP = %q, want %q", got, want)
		}
		if got := r.Header.Get("OK-ACCESS-SIGN"); got == "" {
			t.Fatalf("OK-ACCESS-SIGN empty")
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"result":true}]}`))
	}))
	t.Cleanup(srv.Close)

	c := NewClient(
		WithBaseURL(srv.URL),
		WithHTTPClient(srv.Client()),
		WithCredentials(Credentials{
			APIKey:     "mykey",
			SecretKey:  "mysecret",
			Passphrase: "mypass",
		}),
		WithNowFunc(func() time.Time { return fixedNow }),
	)

	ack, err := c.NewSprdMassCancelService().SprdId("BTC-USDT_BTC-USDT-SWAP").Do(context.Background())
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if ack == nil || !ack.Result {
		t.Fatalf("ack = %#v", ack)
	}

	t.Run("result_false_fail_close", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("x-request-id", "rid-sprd-mass-cancel")
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"result":false}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
			WithCredentials(Credentials{
				APIKey:     "mykey",
				SecretKey:  "mysecret",
				Passphrase: "mypass",
			}),
			WithNowFunc(func() time.Time { return fixedNow }),
		)

		_, err := c.NewSprdMassCancelService().SprdId("BTC-USDT_BTC-USDT-SWAP").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		apiErr, ok := err.(*APIError)
		if !ok {
			t.Fatalf("err = %T, want *APIError: %v", err, err)
		}
		if got, want := apiErr.RequestID, "rid-sprd-mass-cancel"; got != want {
			t.Fatalf("apiErr.RequestID = %q, want %q", got, want)
		}
		if got, want := apiErr.RequestPath, "/api/v5/sprd/mass-cancel"; got != want {
			t.Fatalf("apiErr.RequestPath = %q, want %q", got, want)
		}
	})
}
