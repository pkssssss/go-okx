package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pkssssss/go-okx/v5/internal/sign"
)

func TestOneClickRepayHistoryV2Service_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	timestamp := sign.TimestampISO8601Millis(fixedNow)
	wantQuery := "after=1&before=2&limit=10"
	wantPath := "/api/v5/trade/one-click-repay-history-v2?" + wantQuery
	wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(timestamp, http.MethodGet, wantPath, ""))

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, http.MethodGet; got != want {
			t.Fatalf("method = %q, want %q", got, want)
		}
		if got, want := r.URL.Path, "/api/v5/trade/one-click-repay-history-v2"; got != want {
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
		_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"debtCcy":"USDC","fillDebtSz":"9.079631989","ordIdInfo":[{"ordId":"1","instId":"USDC-USDT","ordType":"ioc","side":"buy","px":"1.0049","sz":"9.088651","fillPx":"1","fillSz":"9.088651","state":"filled","cTime":"1742194485439"}],"repayCcyList":["USDC","BTC"],"status":"filled","ts":"1742194481852"}]}`))
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

	got, err := c.NewOneClickRepayHistoryV2Service().After("1").Before("2").Limit(10).Do(context.Background())
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if len(got) != 1 || got[0].DebtCcy != "USDC" || got[0].TS != 1742194481852 {
		t.Fatalf("data = %#v", got)
	}
	if len(got[0].OrdIdInfo) != 1 || got[0].OrdIdInfo[0].OrdId != "1" || got[0].OrdIdInfo[0].CTime != 1742194485439 {
		t.Fatalf("order info = %#v", got[0].OrdIdInfo)
	}
}
