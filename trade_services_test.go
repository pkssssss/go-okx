package okx

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPlaceOrderService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	t.Run("signed_request_and_body", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/trade/order"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}

			bodyBytes, _ := io.ReadAll(r.Body)
			if got, want := string(bodyBytes), `{"instId":"BTC-USDT","tdMode":"isolated","side":"buy","ordType":"limit","px":"1","sz":"1"}`; got != want {
				t.Fatalf("body = %q, want %q", got, want)
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
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), "5JgIEfkRBluy4x31t6uitZqzoshK+kWWjq9f597WRqQ="; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"clOrdId":"","ordId":"123","tag":"","ts":"1695190491421","sCode":"0","sMsg":""}]}`))
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

		got, err := c.NewPlaceOrderService().
			InstId("BTC-USDT").
			TdMode("isolated").
			Side("buy").
			OrdType("limit").
			Px("1").
			Sz("1").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.OrdId != "123" {
			t.Fatalf("OrdId = %q, want %q", got.OrdId, "123")
		}
	})

	t.Run("item_error_sCode", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"clOrdId":"","ordId":"","tag":"","ts":"0","sCode":"51000","sMsg":"bad"}]}`))
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

		_, err := c.NewPlaceOrderService().
			InstId("BTC-USDT").
			TdMode("isolated").
			Side("buy").
			OrdType("market").
			Sz("1").
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("expected *APIError, got %T: %v", err, err)
		}
		if apiErr.Code != "51000" {
			t.Fatalf("Code = %q, want %q", apiErr.Code, "51000")
		}
	})
}

func TestCancelOrderService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, http.MethodPost; got != want {
			t.Fatalf("method = %q, want %q", got, want)
		}
		if got, want := r.URL.Path, "/api/v5/trade/cancel-order"; got != want {
			t.Fatalf("path = %q, want %q", got, want)
		}

		bodyBytes, _ := io.ReadAll(r.Body)
		if got, want := string(bodyBytes), `{"instId":"BTC-USDT","ordId":"590908157585625111"}`; got != want {
			t.Fatalf("body = %q, want %q", got, want)
		}

		if got, want := r.Header.Get("OK-ACCESS-TIMESTAMP"), "2020-03-28T12:21:41.274Z"; got != want {
			t.Fatalf("OK-ACCESS-TIMESTAMP = %q, want %q", got, want)
		}
		if got, want := r.Header.Get("OK-ACCESS-SIGN"), "gSIagznstauOP+iwBhQ2CLDadQSvMtS/dEDo2uPaN5w="; got != want {
			t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"clOrdId":"","ordId":"123","ts":"1695190491421","sCode":"0","sMsg":""}]}`))
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

	got, err := c.NewCancelOrderService().
		InstId("BTC-USDT").
		OrdId("590908157585625111").
		Do(context.Background())
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if got.OrdId != "123" {
		t.Fatalf("OrdId = %q, want %q", got.OrdId, "123")
	}
}

func TestAmendOrderService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	t.Run("signed_request_and_body", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/trade/amend-order"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}

			bodyBytes, _ := io.ReadAll(r.Body)
			if got, want := string(bodyBytes), `{"instId":"BTC-USDT","ordId":"590909145319051111","newSz":"2"}`; got != want {
				t.Fatalf("body = %q, want %q", got, want)
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
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), "UlfJcIwsuVXKxzc/FN9S+u15AYDIzP/6Qwz2+K2oYFs="; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"clOrdId":"","ordId":"590909145319051111","reqId":"","ts":"1695190491421","sCode":"0","sMsg":""}]}`))
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

		got, err := c.NewAmendOrderService().
			InstId("BTC-USDT").
			OrdId("590909145319051111").
			NewSz("2").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.OrdId != "590909145319051111" {
			t.Fatalf("OrdId = %q, want %q", got.OrdId, "590909145319051111")
		}
	})
}

func TestOrdersPendingService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, http.MethodGet; got != want {
			t.Fatalf("method = %q, want %q", got, want)
		}
		if got, want := r.URL.Path, "/api/v5/trade/orders-pending"; got != want {
			t.Fatalf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.RawQuery, "instId=BTC-USDT&instType=SPOT&limit=2&ordType=post_only"; got != want {
			t.Fatalf("query = %q, want %q", got, want)
		}

		if got, want := r.Header.Get("OK-ACCESS-TIMESTAMP"), "2020-03-28T12:21:41.274Z"; got != want {
			t.Fatalf("OK-ACCESS-TIMESTAMP = %q, want %q", got, want)
		}
		if got, want := r.Header.Get("OK-ACCESS-SIGN"), "fqMXMl+1KFCoeYIHOlNAhNrEEObODI2bWHg0PzdH41c="; got != want {
			t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instType":"SPOT","instId":"BTC-USDT","ordId":"1","clOrdId":"","tag":"","side":"buy","ordType":"post_only","state":"live","px":"1","sz":"1","avgPx":"0","fillPx":"0","fillSz":"0","accFillSz":"0","uTime":"1597026383085","cTime":"1597026383085"}]}`))
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

	got, err := c.NewOrdersPendingService().
		InstType("SPOT").
		InstId("BTC-USDT").
		OrdType("post_only").
		Limit(2).
		Do(context.Background())
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want %d", len(got), 1)
	}
	if got[0].OrdId != "1" {
		t.Fatalf("OrdId = %q, want %q", got[0].OrdId, "1")
	}
}

func TestOrdersHistoryService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	t.Run("missing_instType", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewOrdersHistoryService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errOrdersHistoryMissingInstType {
			t.Fatalf("error = %v, want %v", err, errOrdersHistoryMissingInstType)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/trade/orders-history"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "instType=SPOT&limit=2&ordType=limit&state=filled"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-TIMESTAMP"), "2020-03-28T12:21:41.274Z"; got != want {
				t.Fatalf("OK-ACCESS-TIMESTAMP = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), "OU6JDIKCbJDDZYHkAyBYYXk/wHdPfGsP9P8bale5XDI="; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instType":"SPOT","instId":"BTC-USDT","ordId":"2","clOrdId":"","tag":"","side":"buy","ordType":"limit","state":"filled","px":"1","sz":"1","avgPx":"1","fillPx":"1","fillSz":"1","accFillSz":"1","uTime":"1597026383085","cTime":"1597026383085"}]}`))
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

		got, err := c.NewOrdersHistoryService().
			InstType("SPOT").
			OrdType("limit").
			State("filled").
			Limit(2).
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("len = %d, want %d", len(got), 1)
		}
		if got[0].OrdId != "2" {
			t.Fatalf("OrdId = %q, want %q", got[0].OrdId, "2")
		}
	})
}

func TestTradeFillsService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, http.MethodGet; got != want {
			t.Fatalf("method = %q, want %q", got, want)
		}
		if got, want := r.URL.Path, "/api/v5/trade/fills"; got != want {
			t.Fatalf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.RawQuery, "instId=BTC-USDT&instType=SPOT&limit=2&ordId=123"; got != want {
			t.Fatalf("query = %q, want %q", got, want)
		}

		if got, want := r.Header.Get("OK-ACCESS-TIMESTAMP"), "2020-03-28T12:21:41.274Z"; got != want {
			t.Fatalf("OK-ACCESS-TIMESTAMP = %q, want %q", got, want)
		}
		if got, want := r.Header.Get("OK-ACCESS-SIGN"), "75LqTfNVapL9ciV+fgyLfx7moOr0/i1/Cl9oZp5HH8Q="; got != want {
			t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instType":"SPOT","instId":"BTC-USDT","ordId":"123","clOrdId":"","tradeId":"1","side":"buy","posSide":"net","fillPx":"1","fillSz":"1","fee":"-0.1","feeCcy":"BTC","fillTime":"1597026383085"}]}`))
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

	got, err := c.NewTradeFillsService().
		InstType("SPOT").
		InstId("BTC-USDT").
		OrdId("123").
		Limit(2).
		Do(context.Background())
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want %d", len(got), 1)
	}
	if got[0].TradeId != "1" {
		t.Fatalf("TradeId = %q, want %q", got[0].TradeId, "1")
	}
	if got[0].FillTime != 1597026383085 {
		t.Fatalf("FillTime = %d, want %d", got[0].FillTime, 1597026383085)
	}
}
