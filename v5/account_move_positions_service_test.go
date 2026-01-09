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

func TestAccountMovePositionsService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("missing_required", func(t *testing.T) {
		c := NewClient(
			WithNowFunc(func() time.Time { return fixedNow }),
		)
		_, err := c.NewAccountMovePositionsService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errAccountMovePositionsMissingRequired {
			t.Fatalf("error = %v, want %v", err, errAccountMovePositionsMissingRequired)
		}
	})

	t.Run("signed_request_and_body", func(t *testing.T) {
		timestamp := sign.TimestampISO8601Millis(fixedNow)
		wantBody := `{"fromAcct":"0","toAcct":"test","legs":[{"from":{"posId":"2065471111340792832","sz":"1","side":"sell"},"to":{"tdMode":"cross","posSide":"net"}}],"clientId":"test"}`
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodPost, "/api/v5/account/move-positions", wantBody))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/account/move-positions"; got != want {
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
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"clientId":"test","blockTdId":"2065832911119076864","state":"filled","ts":"1734069018526","fromAcct":"0","toAcct":"test","legs":[{"from":{"posId":"2065471111340792832","instId":"BTC-USD-SWAP","px":"100042.7","side":"sell","sz":"1","sCode":"0","sMsg":""},"to":{"instId":"BTC-USD-SWAP","px":"100042.7","side":"buy","sz":"1","tdMode":"cross","posSide":"net","ccy":"","sCode":"0","sMsg":""}}]}]}`))
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

		got, err := c.NewAccountMovePositionsService().
			FromAcct("0").
			ToAcct("test").
			ClientId("test").
			Legs([]AccountMovePositionsLeg{
				{
					From: AccountMovePositionsLegFrom{PosId: "2065471111340792832", Sz: "1", Side: "sell"},
					To:   AccountMovePositionsLegTo{TdMode: "cross", PosSide: "net"},
				},
			}).
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.ClientId != "test" || got.BlockTdId != "2065832911119076864" || got.State != "filled" || got.TS != 1734069018526 {
			t.Fatalf("ack = %#v", got)
		}
		if len(got.Legs) != 1 || got.Legs[0].From.PosId != "2065471111340792832" || got.Legs[0].To.TdMode != "cross" {
			t.Fatalf("legs = %#v", got.Legs)
		}
	})

	t.Run("missing_credentials", func(t *testing.T) {
		c := NewClient(
			WithNowFunc(func() time.Time { return fixedNow }),
		)
		_, err := c.NewAccountMovePositionsService().
			FromAcct("0").
			ToAcct("test").
			ClientId("test").
			Legs([]AccountMovePositionsLeg{
				{
					From: AccountMovePositionsLegFrom{PosId: "2065471111340792832", Sz: "1", Side: "sell"},
					To:   AccountMovePositionsLegTo{TdMode: "cross", PosSide: "net"},
				},
			}).
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMissingCredentials {
			t.Fatalf("error = %v, want %v", err, errMissingCredentials)
		}
	})
}
