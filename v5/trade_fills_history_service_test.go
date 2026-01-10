package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pkssssss/go-okx/v5/internal/sign"
)

func TestTradeFillsHistoryService_Do_WithAfterBeforeAndInstFamily(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("missing_inst_type", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewTradeFillsHistoryService().Do(context.Background())
		if err != errTradeFillsHistoryMissingInstType {
			t.Fatalf("error = %v, want %v", err, errTradeFillsHistoryMissingInstType)
		}
	})

	t.Run("signed_request", func(t *testing.T) {
		timestamp := sign.TimestampISO8601Millis(fixedNow)
		wantQuery := "after=1&before=2&begin=111&end=222&instFamily=BTC-USD&instId=BTC-USDT&instType=SPOT&limit=50&ordId=123&subType=1"
		wantPath := "/api/v5/trade/fills-history?" + wantQuery
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodGet, wantPath, ""))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/trade/fills-history"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, wantQuery; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}

			if got, want := r.Header.Get("OK-ACCESS-KEY"), "mykey"; got != want {
				t.Fatalf("OK-ACCESS-KEY = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-PASSPHRASE"), "mypass"; got != want {
				t.Fatalf("OK-ACCESS-PASSPHRASE = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-TIMESTAMP"), timestamp; got != want {
				t.Fatalf("OK-ACCESS-TIMESTAMP = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), wantSig; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instType":"SPOT","instId":"BTC-USDT","tradeId":"t1","ordId":"123","side":"buy","fillPx":"1","fillSz":"1","fee":"0.1","feeCcy":"USDT","ts":"1695190491421","fillTime":"1695190491421"}]}`))
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

		got, err := c.NewTradeFillsHistoryService().
			InstType("SPOT").
			InstFamily("BTC-USD").
			InstId("BTC-USDT").
			OrdId("123").
			SubType("1").
			After("1").
			Before("2").
			Begin("111").
			End("222").
			Limit(50).
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].OrdId != "123" || got[0].TS != 1695190491421 {
			t.Fatalf("data = %#v", got)
		}
	})

	t.Run("missing_credentials", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewTradeFillsHistoryService().InstType("SPOT").Do(context.Background())
		if err != errMissingCredentials {
			t.Fatalf("error = %v, want %v", err, errMissingCredentials)
		}
	})
}
