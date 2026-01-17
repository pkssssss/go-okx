package okx

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pkssssss/go-okx/v5/internal/sign"
)

func TestAmendAlgoOrderService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("missing_inst_id", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewAmendAlgoOrderService().Do(context.Background())
		if err != errAmendAlgoOrderMissingInstId {
			t.Fatalf("error = %v, want %v", err, errAmendAlgoOrderMissingInstId)
		}
	})

	t.Run("missing_id", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewAmendAlgoOrderService().InstId("BTC-USDT").NewSz("2").Do(context.Background())
		if err != errAmendAlgoOrderMissingId {
			t.Fatalf("error = %v, want %v", err, errAmendAlgoOrderMissingId)
		}
	})

	t.Run("missing_change", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewAmendAlgoOrderService().InstId("BTC-USDT").AlgoId("1").Do(context.Background())
		if err != errAmendAlgoOrderMissingAnyChange {
			t.Fatalf("error = %v, want %v", err, errAmendAlgoOrderMissingAnyChange)
		}
	})

	t.Run("signed_request_and_body", func(t *testing.T) {
		timestamp := sign.TimestampISO8601Millis(fixedNow)
		wantBody := `{"instId":"BTC-USDT","algoId":"2510789768709120","newSz":"2"}`
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodPost, "/api/v5/trade/amend-algos", wantBody))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/trade/amend-algos"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}

			bodyBytes, _ := io.ReadAll(r.Body)
			if got, want := string(bodyBytes), wantBody; got != want {
				t.Fatalf("body = %q, want %q", got, want)
			}

			if got, want := r.Header.Get("OK-ACCESS-TIMESTAMP"), timestamp; got != want {
				t.Fatalf("OK-ACCESS-TIMESTAMP = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), wantSig; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"algoClOrdId":"algo_01","algoId":"2510789768709120","reqId":"po103ux","sCode":"0","sMsg":""}]}`))
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

		got, err := c.NewAmendAlgoOrderService().
			InstId("BTC-USDT").
			AlgoId("2510789768709120").
			NewSz("2").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.AlgoId != "2510789768709120" || got.ReqId != "po103ux" || got.SCode != "0" {
			t.Fatalf("ack = %#v", got)
		}
	})

	t.Run("ack_error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("x-request-id", "rid-amend-algo-1")
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"algoClOrdId":"algo_01","algoId":"1","reqId":"po","sCode":"51000","sMsg":"failed"}]}`))
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

		_, err := c.NewAmendAlgoOrderService().InstId("BTC-USDT").AlgoId("1").NewSz("2").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		apiErr, ok := err.(*APIError)
		if !ok {
			t.Fatalf("error = %T, want *APIError", err)
		}
		if got, want := apiErr.RequestID, "rid-amend-algo-1"; got != want {
			t.Fatalf("RequestID = %q, want %q", got, want)
		}
	})

	t.Run("missing_credentials", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewAmendAlgoOrderService().InstId("BTC-USDT").AlgoId("1").NewSz("2").Do(context.Background())
		if err != errMissingCredentials {
			t.Fatalf("error = %v, want %v", err, errMissingCredentials)
		}
	})
}
