package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pkssssss/go-okx/v5/internal/sign"
)

func TestOneClickRepayCurrencyListV2Service_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	timestamp := sign.TimestampISO8601Millis(fixedNow)
	wantPath := "/api/v5/trade/one-click-repay-currency-list-v2"
	wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodGet, wantPath, ""))

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, http.MethodGet; got != want {
			t.Fatalf("method = %q, want %q", got, want)
		}
		if got, want := r.URL.Path, wantPath; got != want {
			t.Fatalf("path = %q, want %q", got, want)
		}
		if got := r.URL.RawQuery; got != "" {
			t.Fatalf("query = %q, want empty", got)
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
		_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"debtData":[{"debtAmt":"100","debtCcy":"USDC"}],"repayData":[{"repayAmt":"1","repayCcy":"BTC"}],"debtCcy":"USDC","repayCcyList":["USDC","BTC"],"ts":"1695190491421"}]}`))
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

	got, err := c.NewOneClickRepayCurrencyListV2Service().Do(context.Background())
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if len(got) != 1 || got[0].DebtCcy != "USDC" || got[0].TS != 1695190491421 {
		t.Fatalf("data = %#v", got)
	}
}
