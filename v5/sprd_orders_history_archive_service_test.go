package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSprdOrdersHistoryArchiveService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, http.MethodGet; got != want {
			t.Fatalf("method = %q, want %q", got, want)
		}
		if got, want := r.URL.Path, "/api/v5/sprd/orders-history-archive"; got != want {
			t.Fatalf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.RawQuery, "begin=1&beginId=11&end=9&endId=99&instFamily=BTC-USDT&instType=SWAP&limit=10&ordType=limit&sprdId=BTC-USDT_BTC-USDT-SWAP&state=filled"; got != want {
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
		_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"sprdId":"BTC-USDT_BTC-USDT-SWAP","ordId":"o1","clOrdId":"c1","px":"1","sz":"1","ordType":"limit","side":"buy","state":"filled","uTime":"1597026383085","cTime":"1597026383085"}]}`))
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

	orders, err := c.NewSprdOrdersHistoryArchiveService().
		SprdId("BTC-USDT_BTC-USDT-SWAP").
		OrdType("limit").
		State("filled").
		InstType("SWAP").
		InstFamily("BTC-USDT").
		BeginId("11").
		EndId("99").
		Begin("1").
		End("9").
		Limit(10).
		Do(context.Background())
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if len(orders) != 1 || orders[0].OrdId != "o1" || orders[0].State != "filled" {
		t.Fatalf("orders = %#v", orders)
	}
}
