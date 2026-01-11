package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSprdTradesService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, http.MethodGet; got != want {
			t.Fatalf("method = %q, want %q", got, want)
		}
		if got, want := r.URL.Path, "/api/v5/sprd/trades"; got != want {
			t.Fatalf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.RawQuery, "begin=1&beginId=11&end=9&endId=99&limit=10&ordId=123445&sprdId=BTC-USDT_BTC-USDT-SWAP&tradeId=123"; got != want {
			t.Fatalf("query = %q, want %q", got, want)
		}
		if got, want := r.Header.Get("OK-ACCESS-KEY"), "mykey"; got != want {
			t.Fatalf("OK-ACCESS-KEY = %q, want %q", got, want)
		}
		if got, want := r.Header.Get("OK-ACCESS-PASSPHRASE"), "mypass"; got != want {
			t.Fatalf("OK-ACCESS-PASSPHRASE = %q, want %q", got, want)
		}
		if got, want := r.Header.Get("OK-ACCESS-TIMESTAMP"), "2020-03-28T12:21:41.274Z"; got != want {
			t.Fatalf("OK-ACCESS-TIMESTAMP = %q, want %q", got, want)
		}
		if got := r.Header.Get("OK-ACCESS-SIGN"); got == "" {
			t.Fatalf("OK-ACCESS-SIGN empty")
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"sprdId":"BTC-USDT-SWAP_BTC-USDT-200329","tradeId":"123","ordId":"123445","clOrdId":"b16","tag":"","fillPx":"999","fillSz":"3","state":"filled","side":"buy","execType":"M","ts":"1597026383085","legs":[{"instId":"BTC-USDT-SWAP","px":"20000","sz":"3","szCont":"0.03","side":"buy","fillPnl":"","fee":"","feeCcy":"","tradeId":"1232342342"}],"code":"","msg":""}]}`))
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

	trades, err := c.NewSprdTradesService().
		SprdId("BTC-USDT_BTC-USDT-SWAP").
		TradeId("123").
		OrdId("123445").
		BeginId("11").
		EndId("99").
		Begin("1").
		End("9").
		Limit(10).
		Do(context.Background())
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if len(trades) != 1 || trades[0].TradeId != "123" || trades[0].TS != 1597026383085 || len(trades[0].Legs) != 1 {
		t.Fatalf("trades = %#v", trades)
	}
}
