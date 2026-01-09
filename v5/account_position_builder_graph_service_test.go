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

func TestAccountPositionBuilderGraphService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("missing_required", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewAccountPositionBuilderGraphService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errAccountPositionBuilderGraphMissingRequired {
			t.Fatalf("error = %v, want %v", err, errAccountPositionBuilderGraphMissingRequired)
		}
	})

	t.Run("signed_request_and_body", func(t *testing.T) {
		timestamp := sign.TimestampISO8601Millis(fixedNow)
		wantBody := `{"inclRealPosAndEq":false,"simPos":[{"instId":"BTC-USDT-SWAP","pos":"-10","avgPx":"100000"},{"instId":"LTC-USDT-SWAP","pos":"10","avgPx":"8000"}],"simAsset":[{"ccy":"USDT","amt":"100"}],"greeksType":"CASH","type":"mmr","mmrConfig":{"acctLv":"3","lever":"1"}}`
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodPost, "/api/v5/account/position-builder-graph", wantBody))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/account/position-builder-graph"; got != want {
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
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"type":"mmr","mmrData":[{"shockFactor":"-0.94","mmr":"1415.0254039225917","mmrRatio":"-47.45603627655477"},{"shockFactor":"-0.93","mmr":"1417.732491243024","mmrRatio":"-47.436684685735386"}]}]}`))
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

		got, err := c.NewAccountPositionBuilderGraphService().
			InclRealPosAndEq(false).
			SimPos([]AccountPositionBuilderSimPos{
				{InstId: "BTC-USDT-SWAP", Pos: "-10", AvgPx: "100000"},
				{InstId: "LTC-USDT-SWAP", Pos: "10", AvgPx: "8000"},
			}).
			SimAsset([]AccountPositionBuilderSimAsset{
				{Ccy: "USDT", Amt: "100"},
			}).
			GreeksType("CASH").
			Type("mmr").
			MmrConfig(AccountPositionBuilderGraphMmrConfig{AcctLv: "3", Lever: "1"}).
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("len(data) = %d, want 1", len(got))
		}
		if got[0].Type != "mmr" || len(got[0].MmrData) != 2 || got[0].MmrData[0].ShockFactor != "-0.94" {
			t.Fatalf("data = %#v", got)
		}
	})

	t.Run("empty_data_ok", func(t *testing.T) {
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

		got, err := c.NewAccountPositionBuilderGraphService().
			Type("mmr").
			MmrConfig(AccountPositionBuilderGraphMmrConfig{}).
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("len(data) = %d, want 0", len(got))
		}
	})

	t.Run("missing_credentials", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewAccountPositionBuilderGraphService().
			Type("mmr").
			MmrConfig(AccountPositionBuilderGraphMmrConfig{}).
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMissingCredentials {
			t.Fatalf("error = %v, want %v", err, errMissingCredentials)
		}
	})
}
