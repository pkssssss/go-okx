package okx

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/pkssssss/go-okx/v5/internal/sign"
)

func TestCancelAlgoOrdersService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("missing_orders", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewCancelAlgoOrdersService().Do(context.Background())
		if err != errCancelAlgoOrdersMissingOrders {
			t.Fatalf("error = %v, want %v", err, errCancelAlgoOrdersMissingOrders)
		}
	})

	t.Run("too_many_orders", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		orders := make([]CancelAlgoOrder, tradeAlgoMaxCancelOrders+1)
		for i := range orders {
			orders[i] = CancelAlgoOrder{InstId: "BTC-USDT", AlgoId: "1"}
		}
		_, err := c.NewCancelAlgoOrdersService().Orders(orders).Do(context.Background())
		if err != errCancelAlgoOrdersTooManyOrders {
			t.Fatalf("error = %v, want %v", err, errCancelAlgoOrdersTooManyOrders)
		}
	})

	t.Run("signed_request_and_body", func(t *testing.T) {
		timestamp := sign.TimestampISO8601Millis(fixedNow)
		wantBody := `[{"instId":"BTC-USDT","algoId":"590919993110396111"}]`
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodPost, "/api/v5/trade/cancel-algos", wantBody))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if handleTradeAccountRateLimitMock(w, r) {
				return
			}
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/trade/cancel-algos"; got != want {
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
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"algoClOrdId":"","algoId":"590919993110396111","clOrdId":"","sCode":"0","sMsg":"","tag":""}]}`))
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

		got, err := c.NewCancelAlgoOrdersService().Orders([]CancelAlgoOrder{
			{InstId: "BTC-USDT", AlgoId: "590919993110396111", AlgoClOrdId: "ignored"},
		}).Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].AlgoId != "590919993110396111" {
			t.Fatalf("acks = %#v", got)
		}
	})

	t.Run("partial_failure", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if handleTradeAccountRateLimitMock(w, r) {
				return
			}
			w.Header().Set("x-request-id", "rid-algo-1")
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"algoClOrdId":"","algoId":"1","clOrdId":"","sCode":"51000","sMsg":"failed","tag":""}]}`))
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

		acks, err := c.NewCancelAlgoOrdersService().Orders([]CancelAlgoOrder{{InstId: "BTC-USDT", AlgoId: "1"}}).Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		var batchErr *TradeAlgoBatchError
		if !errors.As(err, &batchErr) {
			t.Fatalf("error = %T, want *TradeAlgoBatchError", err)
		}
		if got, want := batchErr.RequestID, "rid-algo-1"; got != want {
			t.Fatalf("RequestID = %q, want %q", got, want)
		}
		if !strings.Contains(err.Error(), "requestId=rid-algo-1") {
			t.Fatalf("err.Error() = %q", err.Error())
		}
		if len(acks) != 1 || acks[0].SCode != "51000" {
			t.Fatalf("acks = %#v", acks)
		}
	})

	t.Run("missing_credentials", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewCancelAlgoOrdersService().Orders([]CancelAlgoOrder{{InstId: "BTC-USDT", AlgoId: "1"}}).Do(context.Background())
		if err != errMissingCredentials {
			t.Fatalf("error = %v, want %v", err, errMissingCredentials)
		}
	})
}
