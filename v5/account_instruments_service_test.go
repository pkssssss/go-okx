package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pkssssss/go-okx/v5/internal/sign"
)

func TestAccountInstrumentsService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("missing_inst_type", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewAccountInstrumentsService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errAccountInstrumentsMissingInstType {
			t.Fatalf("error = %v, want %v", err, errAccountInstrumentsMissingInstType)
		}
	})

	t.Run("signed_request", func(t *testing.T) {
		timestamp := sign.TimestampISO8601Millis(fixedNow)
		wantQuery := "instFamily=BTC-USD&instId=BTC-USD-240628-50000-C&instType=OPTION"
		wantPath := "/api/v5/account/instruments?" + wantQuery
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodGet, wantPath, ""))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/account/instruments"; got != want {
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
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instType":"OPTION","instId":"BTC-USD-240628-50000-C","instFamily":"BTC-USD","uly":"BTC-USD","category":"1","baseCcy":"","quoteCcy":"USD","settleCcy":"BTC","tickSz":"1","lotSz":"1","minSz":"1","ctVal":"0.1","ctMult":"1","ctValCcy":"BTC","groupId":"1","tradeQuoteCcyList":["USD"],"state":"live"}]}`))
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

		got, err := c.NewAccountInstrumentsService().
			InstType("OPTION").
			InstFamily("BTC-USD").
			InstId("BTC-USD-240628-50000-C").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("len(data) = %d, want 1", len(got))
		}
		if got[0].InstId != "BTC-USD-240628-50000-C" || got[0].TickSz != "1" || got[0].CtVal != "0.1" || got[0].GroupId != "1" {
			t.Fatalf("data[0] = %#v", got[0])
		}
	})

	t.Run("missing_credentials", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewAccountInstrumentsService().InstType("SPOT").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMissingCredentials {
			t.Fatalf("error = %v, want %v", err, errMissingCredentials)
		}
	})
}
