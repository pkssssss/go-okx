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

func TestOrderPrecheckService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("missing_inst_id", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewOrderPrecheckService().Do(context.Background())
		if err != errOrderPrecheckMissingInstId {
			t.Fatalf("error = %v, want %v", err, errOrderPrecheckMissingInstId)
		}
	})

	t.Run("missing_px_for_limit", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewOrderPrecheckService().InstId("BTC-USDT").TdMode("cash").Side("buy").OrdType("limit").Sz("2").Do(context.Background())
		if err != errOrderPrecheckMissingPx {
			t.Fatalf("error = %v, want %v", err, errOrderPrecheckMissingPx)
		}
	})

	t.Run("signed_request_and_body", func(t *testing.T) {
		timestamp := sign.TimestampISO8601Millis(fixedNow)
		wantBody := `{"instId":"BTC-USDT","tdMode":"cash","side":"buy","ordType":"limit","sz":"2","px":"2.15"}`
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodPost, "/api/v5/trade/order-precheck", wantBody))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/trade/order-precheck"; got != want {
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
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"adjEq":"41.94347460746277","adjEqChg":"-226.05616481626","imr":"0","imrChg":"57.74709688430927","mmr":"0","mmrChg":"0","mgnRatio":"0","mgnRatioChg":"0"}]}`))
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

		got, err := c.NewOrderPrecheckService().
			InstId("BTC-USDT").
			TdMode("cash").
			Side("buy").
			OrdType("limit").
			Sz("2").
			Px("2.15").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].AdjEq == "" || got[0].AdjEqChg == "" {
			t.Fatalf("data = %#v", got)
		}
	})

	t.Run("missing_credentials", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewOrderPrecheckService().InstId("BTC-USDT").TdMode("cash").Side("buy").OrdType("market").Sz("2").Do(context.Background())
		if err != errMissingCredentials {
			t.Fatalf("error = %v, want %v", err, errMissingCredentials)
		}
	})
}
