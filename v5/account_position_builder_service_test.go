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

func TestAccountPositionBuilderService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("simpos_missing_instid", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewAccountPositionBuilderService().
			SimPos([]AccountPositionBuilderSimPos{{Pos: "1", AvgPx: "100"}}).
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if got, want := err.Error(), "okx: position builder simPos[0] missing instId"; got != want {
			t.Fatalf("error = %q, want %q", got, want)
		}
	})

	t.Run("too_many_simpos", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		simPos := make([]AccountPositionBuilderSimPos, accountPositionBuilderMaxSimPos+1)
		_, err := c.NewAccountPositionBuilderService().SimPos(simPos).Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errAccountPositionBuilderTooManySimPos {
			t.Fatalf("error = %v, want %v", err, errAccountPositionBuilderTooManySimPos)
		}
	})

	t.Run("too_many_simasset", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		simAsset := make([]AccountPositionBuilderSimAsset, accountPositionBuilderMaxSimAsset+1)
		_, err := c.NewAccountPositionBuilderService().SimAsset(simAsset).Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errAccountPositionBuilderTooManySimAsset {
			t.Fatalf("error = %v, want %v", err, errAccountPositionBuilderTooManySimAsset)
		}
	})

	t.Run("signed_request_and_body", func(t *testing.T) {
		timestamp := sign.TimestampISO8601Millis(fixedNow)
		wantBody := `{"inclRealPosAndEq":false,"simPos":[{"instId":"BTC-USDT-SWAP","pos":"10","avgPx":"100000"}],"simAsset":[{"ccy":"USDT","amt":"100"}],"greeksType":"CASH"}`
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodPost, "/api/v5/account/position-builder", wantBody))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/account/position-builder"; got != want {
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
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"eq":"-78553.21888979999","totalImr":"9643.45070718144","totalMmr":"2946.5480841048","borrowMmr":"1571.064377796","derivMmr":"1375.4837063088003","marginRatio":"-25.95365779811705","upl":"-78653.1728898","acctLever":"-0.1364949794742562","ts":"1736936801642","assets":[],"positions":[],"riskUnitData":[]} ]}`))
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

		got, err := c.NewAccountPositionBuilderService().
			InclRealPosAndEq(false).
			SimPos([]AccountPositionBuilderSimPos{
				{InstId: "BTC-USDT-SWAP", Pos: "10", AvgPx: "100000"},
			}).
			SimAsset([]AccountPositionBuilderSimAsset{
				{Ccy: "USDT", Amt: "100"},
			}).
			GreeksType("CASH").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.Eq != "-78553.21888979999" || got.TotalImr != "9643.45070718144" || got.TotalMmr != "2946.5480841048" || got.TS != 1736936801642 {
			t.Fatalf("result = %#v", got)
		}
	})

	t.Run("empty_response", func(t *testing.T) {
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

		_, err := c.NewAccountPositionBuilderService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errEmptyAccountPositionBuilder {
			t.Fatalf("error = %v, want %v", err, errEmptyAccountPositionBuilder)
		}
	})

	t.Run("missing_credentials", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewAccountPositionBuilderService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMissingCredentials {
			t.Fatalf("error = %v, want %v", err, errMissingCredentials)
		}
	})
}
