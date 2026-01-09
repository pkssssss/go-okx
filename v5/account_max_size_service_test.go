package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pkssssss/go-okx/v5/internal/sign"
)

func TestAccountMaxSizeService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("missing_required", func(t *testing.T) {
		c := NewClient(
			WithNowFunc(func() time.Time { return fixedNow }),
		)
		_, err := c.NewAccountMaxSizeService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errAccountMaxSizeMissingRequired {
			t.Fatalf("error = %v, want %v", err, errAccountMaxSizeMissingRequired)
		}
	})

	t.Run("signed_request", func(t *testing.T) {
		timestamp := sign.TimestampISO8601Millis(fixedNow)
		wantQuery := "ccy=BTC&instId=BTC-USDT&leverage=5&px=30000&tdMode=isolated&tradeQuoteCcy=USDT"
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodGet, "/api/v5/account/max-size?"+wantQuery, ""))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/account/max-size"; got != want {
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
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"ccy":"BTC","instId":"BTC-USDT","maxBuy":"0.1","maxSell":"0.2"}]}`))
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

		got, err := c.NewAccountMaxSizeService().
			InstId("BTC-USDT").
			TdMode("isolated").
			Ccy("BTC").
			Px("30000").
			Leverage("5").
			TradeQuoteCcy("USDT").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("len(data) = %d, want %d", len(got), 1)
		}
		if got[0].InstId != "BTC-USDT" || got[0].Ccy != "BTC" || got[0].MaxBuy != "0.1" || got[0].MaxSell != "0.2" {
			t.Fatalf("data[0] = %#v", got[0])
		}
	})

	t.Run("empty_list_is_ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[]}`))
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

		got, err := c.NewAccountMaxSizeService().InstId("BTC-USDT").TdMode("cash").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("len(data) = %d, want %d", len(got), 0)
		}
	})

	t.Run("missing_credentials", func(t *testing.T) {
		c := NewClient(
			WithNowFunc(func() time.Time { return fixedNow }),
		)
		_, err := c.NewAccountMaxSizeService().InstId("BTC-USDT").TdMode("cash").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMissingCredentials {
			t.Fatalf("error = %v, want %v", err, errMissingCredentials)
		}
	})
}
