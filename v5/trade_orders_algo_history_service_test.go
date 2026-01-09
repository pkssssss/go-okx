package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pkssssss/go-okx/v5/internal/sign"
)

func TestAlgoOrdersHistoryService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("missing_ord_type", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewAlgoOrdersHistoryService().State("effective").Do(context.Background())
		if err != errAlgoOrdersHistoryMissingOrdType {
			t.Fatalf("error = %v, want %v", err, errAlgoOrdersHistoryMissingOrdType)
		}
	})

	t.Run("missing_state_or_algo_id", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewAlgoOrdersHistoryService().OrdType("conditional").Do(context.Background())
		if err != errAlgoOrdersHistoryMissingStateOrAlgoId {
			t.Fatalf("error = %v, want %v", err, errAlgoOrdersHistoryMissingStateOrAlgoId)
		}
	})

	t.Run("state_and_algo_id_conflict", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewAlgoOrdersHistoryService().OrdType("conditional").State("effective").AlgoId("1").Do(context.Background())
		if err != errAlgoOrdersHistoryStateAndAlgoIdConflict {
			t.Fatalf("error = %v, want %v", err, errAlgoOrdersHistoryStateAndAlgoIdConflict)
		}
	})

	t.Run("signed_request_state", func(t *testing.T) {
		timestamp := sign.TimestampISO8601Millis(fixedNow)
		wantQuery := "after=1&before=2&instId=BTC-USDT&instType=SWAP&limit=10&ordType=conditional&state=effective"
		wantPath := "/api/v5/trade/orders-algo-history?" + wantQuery
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodGet, wantPath, ""))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/trade/orders-algo-history"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, wantQuery; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}

			if got, want := r.Header.Get("OK-ACCESS-TIMESTAMP"), timestamp; got != want {
				t.Fatalf("OK-ACCESS-TIMESTAMP = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), wantSig; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"algoId":"1","instId":"BTC-USDT","instType":"SWAP","ordType":"conditional","side":"buy","tdMode":"cross","sz":"2","state":"effective","cTime":"1724751378980","uTime":"1724751378980"}]}`))
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

		got, err := c.NewAlgoOrdersHistoryService().
			OrdType("conditional").
			State("effective").
			InstType("SWAP").
			InstId("BTC-USDT").
			After("1").
			Before("2").
			Limit(10).
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].AlgoId != "1" || got[0].State != "effective" {
			t.Fatalf("data = %#v", got)
		}
	})

	t.Run("missing_credentials", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewAlgoOrdersHistoryService().OrdType("conditional").State("effective").Do(context.Background())
		if err != errMissingCredentials {
			t.Fatalf("error = %v, want %v", err, errMissingCredentials)
		}
	})
}
