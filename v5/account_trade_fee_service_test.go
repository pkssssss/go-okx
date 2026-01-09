package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pkssssss/go-okx/v5/internal/sign"
)

func TestAccountTradeFeeService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("missing_inst_type", func(t *testing.T) {
		c := NewClient(
			WithNowFunc(func() time.Time { return fixedNow }),
		)
		_, err := c.NewAccountTradeFeeService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errAccountTradeFeeMissingInstType {
			t.Fatalf("error = %v, want %v", err, errAccountTradeFeeMissingInstType)
		}
	})

	t.Run("group_id_conflict", func(t *testing.T) {
		c := NewClient(
			WithNowFunc(func() time.Time { return fixedNow }),
		)
		_, err := c.NewAccountTradeFeeService().
			InstType("SPOT").
			InstId("BTC-USDT").
			GroupId("1").
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errAccountTradeFeeGroupIdConflict {
			t.Fatalf("error = %v, want %v", err, errAccountTradeFeeGroupIdConflict)
		}
	})

	t.Run("signed_request", func(t *testing.T) {
		timestamp := sign.TimestampISO8601Millis(fixedNow)
		wantQuery := "instId=BTC-USDT&instType=SPOT"
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodGet, "/api/v5/account/trade-fee?"+wantQuery, ""))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/account/trade-fee"; got != want {
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
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"level":"Lv1","instType":"SPOT","feeGroup":[{"groupId":"1","maker":"-0.0008","taker":"-0.001"}],"delivery":"","exercise":"","maker":"-0.0008","taker":"-0.001","makerU":"","takerU":"","makerUSDC":"","takerUSDC":"","ruleType":"normal","category":"1","fiat":[],"ts":"1763979985847"}]}`))
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

		got, err := c.NewAccountTradeFeeService().
			InstType("SPOT").
			InstId("BTC-USDT").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("len(data) = %d, want %d", len(got), 1)
		}
		if got[0].InstType != "SPOT" || got[0].Level != "Lv1" || got[0].Maker != "-0.0008" || got[0].TS != 1763979985847 {
			t.Fatalf("data[0] = %#v", got[0])
		}
		if len(got[0].FeeGroup) != 1 || got[0].FeeGroup[0].GroupId != "1" {
			t.Fatalf("feeGroup = %#v", got[0].FeeGroup)
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

		got, err := c.NewAccountTradeFeeService().InstType("SPOT").Do(context.Background())
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
		_, err := c.NewAccountTradeFeeService().InstType("SPOT").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMissingCredentials {
			t.Fatalf("error = %v, want %v", err, errMissingCredentials)
		}
	})
}
