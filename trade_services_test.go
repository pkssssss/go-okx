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

	t.Run("signed_request_and_body_with_options_and_expTime", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/trade/order"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("expTime"), "1597026383085"; got != want {
				t.Fatalf("expTime = %q, want %q", got, want)
			}

			bodyBytes, _ := io.ReadAll(r.Body)
			if got, want := string(bodyBytes), `{"instId":"BTC-USDT","tdMode":"isolated","ccy":"USDT","clOrdId":"b15","tag":"t1","side":"buy","posSide":"long","ordType":"post_only","px":"1","sz":"1","reduceOnly":true,"tgtCcy":"quote_ccy","banAmend":true,"pxAmendType":"1","tradeQuoteCcy":"USDT","stpMode":"cancel_maker"}`; got != want {
				t.Fatalf("body = %q, want %q", got, want)
			}

			if got, want := r.Header.Get("OK-ACCESS-TIMESTAMP"), "2020-03-28T12:21:41.274Z"; got != want {
				t.Fatalf("OK-ACCESS-TIMESTAMP = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), "ppsaF9o0uI3K4Jpfo7CSfPA4fHoxUPpI27UcBBj5llw="; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"clOrdId":"b15","ordId":"123","tag":"","ts":"1695190491421","sCode":"0","sMsg":""}]}`))
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
			Ccy("USDT").
			ClOrdId("b15").
			Tag("t1").
			Side("buy").
			PosSide("long").
			OrdType("post_only").
			Px("1").
			Sz("1").
			ReduceOnly(true).
			TgtCcy("quote_ccy").
			BanAmend(true).
			PxAmendType("1").
			TradeQuoteCcy("USDT").
			StpMode("cancel_maker").
			ExpTime("1597026383085").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.OrdId != "123" {
			t.Fatalf("OrdId = %q, want %q", got.OrdId, "123")
		}
	})

	t.Run("validate_missing_price_for_post_only", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewPlaceOrderService().
			InstId("BTC-USDT").
			TdMode("isolated").
			Side("buy").
			OrdType("post_only").
			Sz("1").
			Do(context.Background())
		if !errors.Is(err, errPlaceOrderMissingPx) {
			t.Fatalf("expected errPlaceOrderMissingPx, got %T: %v", err, err)
		}
	})

	t.Run("validate_too_many_price_fields", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewPlaceOrderService().
			InstId("BTC-USDT").
			TdMode("isolated").
			Side("buy").
			OrdType("limit").
			Px("1").
			PxUsd("100").
			Sz("1").
			Do(context.Background())
		if !errors.Is(err, errPlaceOrderTooManyPx) {
			t.Fatalf("expected errPlaceOrderTooManyPx, got %T: %v", err, err)
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

func TestBatchPlaceOrdersService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	t.Run("signed_request_and_body", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/trade/batch-orders"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}

			bodyBytes, _ := io.ReadAll(r.Body)
			if got, want := string(bodyBytes), `[{"instId":"BTC-USDT","tdMode":"cash","clOrdId":"b15","side":"buy","ordType":"limit","px":"2.15","sz":"2"},{"instId":"BTC-USDT","tdMode":"cash","clOrdId":"b16","side":"buy","ordType":"limit","px":"2.15","sz":"2"}]`; got != want {
				t.Fatalf("body = %q, want %q", got, want)
			}

			if got, want := r.Header.Get("OK-ACCESS-TIMESTAMP"), "2020-03-28T12:21:41.274Z"; got != want {
				t.Fatalf("OK-ACCESS-TIMESTAMP = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), "San19o20Me1SBBYmsu72mIfBzHimUtsEKtknsPDWZNg="; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"clOrdId":"b15","ordId":"1","tag":"","ts":"1695190491421","sCode":"0","sMsg":""},{"clOrdId":"b16","ordId":"2","tag":"","ts":"1695190491422","sCode":"0","sMsg":""}]}`))
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

		acks, err := c.NewBatchPlaceOrdersService().Orders([]BatchPlaceOrder{
			{InstId: "BTC-USDT", TdMode: "cash", ClOrdId: "b15", Side: "buy", OrdType: "limit", Px: "2.15", Sz: "2"},
			{InstId: "BTC-USDT", TdMode: "cash", ClOrdId: "b16", Side: "buy", OrdType: "limit", Px: "2.15", Sz: "2"},
		}).Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(acks) != 2 || acks[0].OrdId != "1" || acks[1].OrdId != "2" {
			t.Fatalf("acks = %#v, want 2 items ordId 1/2", acks)
		}
	})

	t.Run("partial_failure_returns_TradeBatchError", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"clOrdId":"b15","ordId":"","tag":"","ts":"0","sCode":"51000","sMsg":"bad"},{"clOrdId":"b16","ordId":"2","tag":"","ts":"0","sCode":"0","sMsg":""}]}`))
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

		acks, err := c.NewBatchPlaceOrdersService().Orders([]BatchPlaceOrder{
			{InstId: "BTC-USDT", TdMode: "cash", ClOrdId: "b15", Side: "buy", OrdType: "limit", Px: "2.15", Sz: "2"},
			{InstId: "BTC-USDT", TdMode: "cash", ClOrdId: "b16", Side: "buy", OrdType: "limit", Px: "2.15", Sz: "2"},
		}).Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		var batchErr *TradeBatchError
		if !errors.As(err, &batchErr) {
			t.Fatalf("expected *TradeBatchError, got %T: %v", err, err)
		}
		if len(acks) != 2 {
			t.Fatalf("acks = %#v, want 2 items", acks)
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

func TestBatchCancelOrdersService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, http.MethodPost; got != want {
			t.Fatalf("method = %q, want %q", got, want)
		}
		if got, want := r.URL.Path, "/api/v5/trade/cancel-batch-orders"; got != want {
			t.Fatalf("path = %q, want %q", got, want)
		}

		bodyBytes, _ := io.ReadAll(r.Body)
		if got, want := string(bodyBytes), `[{"instId":"BTC-USDT","ordId":"590908157585625111"},{"instId":"BTC-USDT","ordId":"590908544950571222"}]`; got != want {
			t.Fatalf("body = %q, want %q", got, want)
		}

		if got, want := r.Header.Get("OK-ACCESS-SIGN"), "o4Z4f9pFTnbDKIcYqE3f9gxVjb0JCLzLGheFqFbNmyY="; got != want {
			t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"clOrdId":"","ordId":"1","tag":"","ts":"1695190491421","sCode":"0","sMsg":""},{"clOrdId":"","ordId":"2","tag":"","ts":"1695190491422","sCode":"0","sMsg":""}]}`))
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

	acks, err := c.NewBatchCancelOrdersService().Orders([]BatchCancelOrder{
		{InstId: "BTC-USDT", OrdId: "590908157585625111", ClOrdId: "ignored"},
		{InstId: "BTC-USDT", OrdId: "590908544950571222"},
	}).Do(context.Background())
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if len(acks) != 2 || acks[0].OrdId != "1" {
		t.Fatalf("acks = %#v, want 2 items", acks)
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

	t.Run("signed_request_and_body_with_options_and_expTime", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/trade/amend-order"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("expTime"), "1597026383085"; got != want {
				t.Fatalf("expTime = %q, want %q", got, want)
			}

			bodyBytes, _ := io.ReadAll(r.Body)
			if got, want := string(bodyBytes), `{"instId":"BTC-USDT","cxlOnFail":true,"ordId":"590909145319051111","reqId":"r1","newPxUsd":"1000","pxAmendType":"1"}`; got != want {
				t.Fatalf("body = %q, want %q", got, want)
			}

			if got, want := r.Header.Get("OK-ACCESS-TIMESTAMP"), "2020-03-28T12:21:41.274Z"; got != want {
				t.Fatalf("OK-ACCESS-TIMESTAMP = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), "qOS+hKQ248zdhZNLoJcvBWtpGQueQ7PM8jmDWGL2Ymk="; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"clOrdId":"","ordId":"590909145319051111","reqId":"r1","ts":"1695190491421","sCode":"0","sMsg":""}]}`))
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
			ReqId("r1").
			CxlOnFail(true).
			NewPxUsd("1000").
			PxAmendType("1").
			ExpTime("1597026383085").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.OrdId != "590909145319051111" {
			t.Fatalf("OrdId = %q, want %q", got.OrdId, "590909145319051111")
		}
	})

	t.Run("validate_too_many_price_fields", func(t *testing.T) {
		c := NewClient(WithNowFunc(func() time.Time { return fixedNow }))
		_, err := c.NewAmendOrderService().
			InstId("BTC-USDT").
			OrdId("590909145319051111").
			NewPx("1").
			NewPxUsd("1000").
			Do(context.Background())
		if !errors.Is(err, errAmendOrderTooManyPx) {
			t.Fatalf("expected errAmendOrderTooManyPx, got %T: %v", err, err)
		}
	})
}

func TestBatchAmendOrdersService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, http.MethodPost; got != want {
			t.Fatalf("method = %q, want %q", got, want)
		}
		if got, want := r.URL.Path, "/api/v5/trade/amend-batch-orders"; got != want {
			t.Fatalf("path = %q, want %q", got, want)
		}

		bodyBytes, _ := io.ReadAll(r.Body)
		if got, want := string(bodyBytes), `[{"instId":"BTC-USDT","ordId":"590909308792049444","newSz":"2"},{"instId":"BTC-USDT","ordId":"590909308792049555","newSz":"2"}]`; got != want {
			t.Fatalf("body = %q, want %q", got, want)
		}
		if got, want := r.Header.Get("OK-ACCESS-SIGN"), "D5QJnE66gv/NFrGIxjif9KAnskmiBN72DwYOBXHsw9M="; got != want {
			t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"clOrdId":"","ordId":"590909308792049444","reqId":"","ts":"1695190491421","sCode":"0","sMsg":""},{"clOrdId":"","ordId":"590909308792049555","reqId":"","ts":"1695190491422","sCode":"0","sMsg":""}]}`))
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

	acks, err := c.NewBatchAmendOrdersService().Orders([]BatchAmendOrder{
		{InstId: "BTC-USDT", OrdId: "590909308792049444", ClOrdId: "ignored", NewSz: "2"},
		{InstId: "BTC-USDT", OrdId: "590909308792049555", NewSz: "2"},
	}).Do(context.Background())
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if len(acks) != 2 || acks[0].OrdId != "590909308792049444" {
		t.Fatalf("acks = %#v, want 2 items", acks)
	}
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
		_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instType":"SPOT","instId":"BTC-USDT","ordId":"1","clOrdId":"","tag":"","side":"buy","posSide":"net","tdMode":"cash","ordType":"post_only","state":"live","ccy":"USDT","tgtCcy":"quote_ccy","tradeQuoteCcy":"USDT","reduceOnly":"false","px":"1","pxUsd":"","pxVol":"","pxType":"","sz":"1","avgPx":"0","fillPx":"0","fillSz":"0","accFillSz":"0","tradeId":"","fillTime":"","pnl":"0","fee":"0","feeCcy":"","rebate":"0","rebateCcy":"","stpMode":"cancel_maker","cancelSource":"","cancelSourceReason":"","uTime":"1597026383085","cTime":"1597026383085"}]}`))
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
	if got, want := got[0].TdMode, "cash"; got != want {
		t.Fatalf("TdMode = %q, want %q", got, want)
	}
	if got, want := got[0].PosSide, "net"; got != want {
		t.Fatalf("PosSide = %q, want %q", got, want)
	}
	if got, want := got[0].TradeQuoteCcy, "USDT"; got != want {
		t.Fatalf("TradeQuoteCcy = %q, want %q", got, want)
	}
	if got, want := got[0].ReduceOnly, "false"; got != want {
		t.Fatalf("ReduceOnly = %q, want %q", got, want)
	}
	if got, want := got[0].FillTime, ""; got != want {
		t.Fatalf("FillTime = %q, want %q", got, want)
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
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instType":"SPOT","instId":"BTC-USDT","ordId":"2","clOrdId":"","tag":"","side":"buy","posSide":"net","tdMode":"cash","ordType":"limit","state":"filled","ccy":"USDT","tgtCcy":"base_ccy","tradeQuoteCcy":"USDT","reduceOnly":"false","px":"1","pxUsd":"","pxVol":"","pxType":"","sz":"1","avgPx":"1","fillPx":"1","fillSz":"1","accFillSz":"1","tradeId":"1","fillTime":"1597026383085","pnl":"0","fee":"-0.1","feeCcy":"USDT","rebate":"0","rebateCcy":"","stpMode":"cancel_maker","cancelSource":"","cancelSourceReason":"","uTime":"1597026383085","cTime":"1597026383085"}]}`))
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
		if got, want := got[0].FillTime, "1597026383085"; got != want {
			t.Fatalf("FillTime = %q, want %q", got, want)
		}
	})
}

func TestOrdersHistoryArchiveService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	t.Run("missing_instType", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewOrdersHistoryArchiveService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errOrdersHistoryArchiveMissingInstType {
			t.Fatalf("error = %v, want %v", err, errOrdersHistoryArchiveMissingInstType)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/trade/orders-history-archive"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "begin=1597020000000&category=twap&end=1597026383085&instId=BTC-USDT&instType=SPOT&limit=2&ordType=limit&state=filled"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-TIMESTAMP"), "2020-03-28T12:21:41.274Z"; got != want {
				t.Fatalf("OK-ACCESS-TIMESTAMP = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), "Bcy73sJzeem1LLwEp+u7WcjUY+EyCljibIVh6W1TrOc="; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instType":"SPOT","instId":"BTC-USDT","ordId":"2","clOrdId":"","tag":"","side":"buy","posSide":"net","tdMode":"cash","ordType":"limit","state":"filled","ccy":"USDT","tgtCcy":"base_ccy","tradeQuoteCcy":"USDT","reduceOnly":"false","px":"1","pxUsd":"","pxVol":"","pxType":"","sz":"1","avgPx":"1","fillPx":"1","fillSz":"1","accFillSz":"1","tradeId":"1","fillTime":"1597026383085","pnl":"0","fee":"-0.1","feeCcy":"USDT","rebate":"0","rebateCcy":"","stpMode":"cancel_maker","cancelSource":"","cancelSourceReason":"","uTime":"1597026383085","cTime":"1597026383085"}]}`))
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

		got, err := c.NewOrdersHistoryArchiveService().
			InstType("SPOT").
			InstId("BTC-USDT").
			OrdType("limit").
			State("filled").
			Category("twap").
			Begin("1597020000000").
			End("1597026383085").
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

func TestGetOrderService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	t.Run("missing_instId", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewGetOrderService().OrdId("1").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errGetOrderMissingInstId {
			t.Fatalf("error = %v, want %v", err, errGetOrderMissingInstId)
		}
	})

	t.Run("missing_id", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewGetOrderService().InstId("BTC-USDT").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errGetOrderMissingId {
			t.Fatalf("error = %v, want %v", err, errGetOrderMissingId)
		}
	})

	t.Run("ordId_takes_precedence_over_clOrdId", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/trade/order"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "instId=BTC-USDT&ordId=590909145319051111"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}

			if got, want := r.Header.Get("OK-ACCESS-TIMESTAMP"), "2020-03-28T12:21:41.274Z"; got != want {
				t.Fatalf("OK-ACCESS-TIMESTAMP = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), "qsOWDbfs3eRqc0t1ITnPgIPcgV8fzz9lNwR2pAU1ccY="; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instType":"SPOT","instId":"BTC-USDT","ordId":"590909145319051111","clOrdId":"c1","tag":"","side":"buy","posSide":"net","tdMode":"cash","ordType":"limit","state":"live","ccy":"USDT","tgtCcy":"quote_ccy","tradeQuoteCcy":"USDT","reduceOnly":"true","px":"1","pxUsd":"","pxVol":"","pxType":"","sz":"1","avgPx":"0","fillPx":"0","fillSz":"0","accFillSz":"0","tradeId":"","fillTime":"1597026383085","pnl":"0","fee":"0","feeCcy":"","rebate":"0","rebateCcy":"","stpMode":"cancel_maker","cancelSource":"","cancelSourceReason":"","uTime":"1597026383085","cTime":"1597026383085"}]}`))
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

		got, err := c.NewGetOrderService().
			InstId("BTC-USDT").
			OrdId("590909145319051111").
			ClOrdId("ignored").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.OrdId != "590909145319051111" {
			t.Fatalf("OrdId = %q, want %q", got.OrdId, "590909145319051111")
		}
		if got, want := got.TdMode, "cash"; got != want {
			t.Fatalf("TdMode = %q, want %q", got, want)
		}
		if got, want := got.PosSide, "net"; got != want {
			t.Fatalf("PosSide = %q, want %q", got, want)
		}
		if got, want := got.ReduceOnly, "true"; got != want {
			t.Fatalf("ReduceOnly = %q, want %q", got, want)
		}
		if got, want := got.FillTime, "1597026383085"; got != want {
			t.Fatalf("FillTime = %q, want %q", got, want)
		}
	})

	t.Run("clOrdId_only", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.URL.RawQuery, "clOrdId=c1&instId=BTC-USDT"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instType":"SPOT","instId":"BTC-USDT","ordId":"1","clOrdId":"c1","tag":"","side":"buy","ordType":"limit","state":"live","px":"1","sz":"1","avgPx":"0","fillPx":"0","fillSz":"0","accFillSz":"0","uTime":"1597026383085","cTime":"1597026383085"}]}`))
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

		got, err := c.NewGetOrderService().
			InstId("BTC-USDT").
			ClOrdId("c1").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.ClOrdId != "c1" {
			t.Fatalf("ClOrdId = %q, want %q", got.ClOrdId, "c1")
		}
	})
}

func TestTradeFillsService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	t.Run("basic", func(t *testing.T) {
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
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instType":"SPOT","instId":"BTC-USDT","tradeQuoteCcy":"USDT","tradeId":"1","ordId":"123","clOrdId":"","billId":"b1","subType":"1","tag":"t1","fillPx":"1","fillSz":"1","fillIdxPx":"1.1","fillPnl":"0","fillPxVol":"","fillPxUsd":"","fillMarkVol":"","fillFwdPx":"","fillMarkPx":"","side":"buy","posSide":"net","execType":"T","feeCcy":"BTC","fee":"-0.1","ts":"1597026383085","fillTime":"1597026383085"}]}`))
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
		if got[0].TradeQuoteCcy != "USDT" {
			t.Fatalf("TradeQuoteCcy = %q, want %q", got[0].TradeQuoteCcy, "USDT")
		}
		if got[0].ExecType != "T" {
			t.Fatalf("ExecType = %q, want %q", got[0].ExecType, "T")
		}
		if got[0].TS != 1597026383085 {
			t.Fatalf("TS = %d, want %d", got[0].TS, 1597026383085)
		}
		if got[0].FillTime != 1597026383085 {
			t.Fatalf("FillTime = %d, want %d", got[0].FillTime, 1597026383085)
		}
	})

	t.Run("with_subType_begin_end", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/trade/fills"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "begin=1597020000000&end=1597026383085&instId=BTC-USDT&instType=SPOT&limit=2&ordId=123&subType=1"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), "WunVeXRdF3CgIFvT8VEmYA/V/qg2qunMuOn08PuDglM="; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

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

		_, err := c.NewTradeFillsService().
			InstType("SPOT").
			InstId("BTC-USDT").
			OrdId("123").
			SubType("1").
			Begin("1597020000000").
			End("1597026383085").
			Limit(2).
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
	})
}

func TestTradeFillsHistoryService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	t.Run("missing_instType", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewTradeFillsHistoryService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errTradeFillsHistoryMissingInstType {
			t.Fatalf("error = %v, want %v", err, errTradeFillsHistoryMissingInstType)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/trade/fills-history"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "begin=1597020000000&end=1597026383085&instId=BTC-USDT&instType=SPOT&limit=2&ordId=123&subType=1"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), "w74/lG5sW2Nj6W6ySaydC6lVUvz2mlQ0F7AyvXGKZXU="; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instType":"SPOT","instId":"BTC-USDT","tradeQuoteCcy":"USDT","tradeId":"1","ordId":"123","clOrdId":"","billId":"b1","subType":"1","tag":"t1","fillPx":"1","fillSz":"1","fillIdxPx":"1.1","fillPnl":"0","fillPxVol":"","fillPxUsd":"","fillMarkVol":"","fillFwdPx":"","fillMarkPx":"","side":"buy","posSide":"net","execType":"T","feeCcy":"BTC","fee":"-0.1","ts":"1597026383085","fillTime":"1597026383085"}]}`))
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

		items, err := c.NewTradeFillsHistoryService().
			InstType("SPOT").
			InstId("BTC-USDT").
			OrdId("123").
			SubType("1").
			Begin("1597020000000").
			End("1597026383085").
			Limit(2).
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(items) != 1 {
			t.Fatalf("len = %d, want %d", len(items), 1)
		}
		if items[0].TradeId != "1" {
			t.Fatalf("TradeId = %q, want %q", items[0].TradeId, "1")
		}
		if items[0].TS != 1597026383085 {
			t.Fatalf("TS = %d, want %d", items[0].TS, 1597026383085)
		}
		if items[0].FillTime != 1597026383085 {
			t.Fatalf("FillTime = %d, want %d", items[0].FillTime, 1597026383085)
		}
	})
}
