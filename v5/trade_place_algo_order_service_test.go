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

func TestPlaceAlgoOrderService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("missing_inst_id", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewPlaceAlgoOrderService().Do(context.Background())
		if err != errPlaceAlgoOrderMissingInstId {
			t.Fatalf("error = %v, want %v", err, errPlaceAlgoOrderMissingInstId)
		}
	})

	t.Run("missing_td_mode", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewPlaceAlgoOrderService().InstId("BTC-USDT").Do(context.Background())
		if err != errPlaceAlgoOrderMissingTdMode {
			t.Fatalf("error = %v, want %v", err, errPlaceAlgoOrderMissingTdMode)
		}
	})

	t.Run("missing_side", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewPlaceAlgoOrderService().InstId("BTC-USDT").TdMode("cross").Do(context.Background())
		if err != errPlaceAlgoOrderMissingSide {
			t.Fatalf("error = %v, want %v", err, errPlaceAlgoOrderMissingSide)
		}
	})

	t.Run("missing_ord_type", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewPlaceAlgoOrderService().InstId("BTC-USDT").TdMode("cross").Side("buy").Do(context.Background())
		if err != errPlaceAlgoOrderMissingOrdType {
			t.Fatalf("error = %v, want %v", err, errPlaceAlgoOrderMissingOrdType)
		}
	})

	t.Run("missing_sz_or_close_fraction", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewPlaceAlgoOrderService().
			InstId("BTC-USDT").
			TdMode("cross").
			Side("buy").
			OrdType("conditional").
			Do(context.Background())
		if err != errPlaceAlgoOrderMissingSzOrCloseFraction {
			t.Fatalf("error = %v, want %v", err, errPlaceAlgoOrderMissingSzOrCloseFraction)
		}
	})

	t.Run("sz_and_close_fraction_conflict", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewPlaceAlgoOrderService().
			InstId("BTC-USDT").
			TdMode("cross").
			Side("buy").
			OrdType("conditional").
			Sz("1").
			CloseFraction("1").
			Do(context.Background())
		if err != errPlaceAlgoOrderSzAndCloseFractionConflict {
			t.Fatalf("error = %v, want %v", err, errPlaceAlgoOrderSzAndCloseFractionConflict)
		}
	})

	t.Run("signed_request_and_body", func(t *testing.T) {
		timestamp := sign.TimestampISO8601Millis(fixedNow)
		wantBody := `{"instId":"BTC-USDT","tdMode":"cross","side":"buy","ordType":"conditional","sz":"2","tpTriggerPx":"15","tpOrdPx":"18"}`
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodPost, "/api/v5/trade/order-algo", wantBody))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/trade/order-algo"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, ""; got != want {
				t.Fatalf("query = %q, want %q", got, want)
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
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"algoId":"12345689","clOrdId":"","algoClOrdId":"","sCode":"0","sMsg":"","tag":""}]}`))
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

		got, err := c.NewPlaceAlgoOrderService().
			InstId("BTC-USDT").
			TdMode("cross").
			Side("buy").
			OrdType("conditional").
			Sz("2").
			TpTriggerPx("15").
			TpOrdPx("18").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.AlgoId != "12345689" || got.SCode != "0" {
			t.Fatalf("ack = %#v", got)
		}
	})

	t.Run("missing_credentials", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewPlaceAlgoOrderService().
			InstId("BTC-USDT").
			TdMode("cross").
			Side("buy").
			OrdType("conditional").
			Sz("2").
			TpTriggerPx("15").
			TpOrdPx("18").
			Do(context.Background())
		if err != errMissingCredentials {
			t.Fatalf("error = %v, want %v", err, errMissingCredentials)
		}
	})
}
