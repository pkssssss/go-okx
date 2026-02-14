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

func TestFiatBuySellCurrenciesService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("signed_request", func(t *testing.T) {
		ts := sign.TimestampISO8601Millis(fixedNow)
		requestPath := "/api/v5/fiat/buy-sell/currencies"
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(ts, http.MethodGet, requestPath, ""))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/fiat/buy-sell/currencies"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, ""; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}

			if got, want := r.Header.Get("OK-ACCESS-KEY"), "mykey"; got != want {
				t.Fatalf("OK-ACCESS-KEY = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-PASSPHRASE"), "mypass"; got != want {
				t.Fatalf("OK-ACCESS-PASSPHRASE = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-TIMESTAMP"), ts; got != want {
				t.Fatalf("OK-ACCESS-TIMESTAMP = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-SIGN"), wantSig; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"fiatCcyList":[{"ccy":"USD"},{"ccy":"EUR"}],"cryptoCcyList":[{"ccy":"BTC"}]}]}`))
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

		got, err := c.NewFiatBuySellCurrenciesService().Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got.FiatCcyList) != 2 {
			t.Fatalf("len(fiat) = %d, want %d", len(got.FiatCcyList), 2)
		}
		if got.FiatCcyList[0].Ccy != "USD" {
			t.Fatalf("fiat[0].Ccy = %q, want %q", got.FiatCcyList[0].Ccy, "USD")
		}
	})

	t.Run("empty_data", func(t *testing.T) {
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

		_, err := c.NewFiatBuySellCurrenciesService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errEmptyFiatBuySellCurrenciesResponse {
			t.Fatalf("error = %v, want %v", err, errEmptyFiatBuySellCurrenciesResponse)
		}
	})

	t.Run("missing_credentials", func(t *testing.T) {
		c := NewClient(
			WithNowFunc(func() time.Time { return fixedNow }),
		)
		_, err := c.NewFiatBuySellCurrenciesService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errMissingCredentials {
			t.Fatalf("error = %v, want %v", err, errMissingCredentials)
		}
	})
}

func TestFiatBuySellCurrencyPairService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("missing_required", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewFiatBuySellCurrencyPairService().FromCcy("USD").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errFiatBuySellCurrencyPairMissingRequired {
			t.Fatalf("error = %v, want %v", err, errFiatBuySellCurrencyPairMissingRequired)
		}
	})

	t.Run("signed_request", func(t *testing.T) {
		ts := sign.TimestampISO8601Millis(fixedNow)
		requestPath := "/api/v5/fiat/buy-sell/currency-pair?fromCcy=USD&toCcy=BTC"
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(ts, http.MethodGet, requestPath, ""))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/fiat/buy-sell/currency-pair"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "fromCcy=USD&toCcy=BTC"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}

			if got, want := r.Header.Get("OK-ACCESS-SIGN"), wantSig; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"side":"buy","fromCcy":"USD","toCcy":"BTC","singleTradeMax":"1","singleTradeMin":"0.01","fixedPxRemainingDailyQuota":"","fixedPxDailyLimit":"","paymentMethods":["balance"]}]}`))
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

		got, err := c.NewFiatBuySellCurrencyPairService().FromCcy("USD").ToCcy("BTC").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("len(data) = %d, want %d", len(got), 1)
		}
		if got[0].Side != "buy" || got[0].FromCcy != "USD" || got[0].ToCcy != "BTC" {
			t.Fatalf("data[0] = %#v", got[0])
		}
	})
}

func TestFiatBuySellQuoteService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("missing_required", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewFiatBuySellQuoteService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errFiatBuySellQuoteMissingRequired {
			t.Fatalf("error = %v, want %v", err, errFiatBuySellQuoteMissingRequired)
		}
	})

	t.Run("signed_request_and_body", func(t *testing.T) {
		ts := sign.TimestampISO8601Millis(fixedNow)
		requestPath := "/api/v5/fiat/buy-sell/quote"
		wantBody := `{"side":"buy","fromCcy":"USD","toCcy":"BTC","rfqAmt":"30","rfqCcy":"USD"}`
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(ts, http.MethodPost, requestPath, wantBody))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/fiat/buy-sell/quote"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}

			bodyBytes, _ := io.ReadAll(r.Body)
			if got := string(bodyBytes); got != wantBody {
				t.Fatalf("body = %q, want %q", got, wantBody)
			}

			if got, want := r.Header.Get("OK-ACCESS-SIGN"), wantSig; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"quoteId":"qid","side":"buy","fromCcy":"USD","toCcy":"BTC","rfqAmt":"30","rfqCcy":"USD","quotePx":"2932.4","quoteCcy":"USD","quoteFromAmt":"30","quoteToAmt":"0.01","quoteTime":"1646188510461","ttlMs":"10000"}]}`))
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

		got, err := c.NewFiatBuySellQuoteService().
			Side("buy").
			FromCcy("USD").
			ToCcy("BTC").
			RfqAmt("30").
			RfqCcy("USD").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.QuoteId != "qid" {
			t.Fatalf("QuoteId = %q, want %q", got.QuoteId, "qid")
		}
		if got.QuoteTime != 1646188510461 {
			t.Fatalf("QuoteTime = %d, want %d", got.QuoteTime, 1646188510461)
		}
	})

	t.Run("empty_data", func(t *testing.T) {
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

		_, err := c.NewFiatBuySellQuoteService().
			Side("buy").
			FromCcy("USD").
			ToCcy("BTC").
			RfqAmt("30").
			RfqCcy("USD").
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		assertEmptyDataAPIError(t, err, errEmptyFiatBuySellQuoteResponse)
	})
}

func TestFiatBuySellTradeService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("missing_required", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewFiatBuySellTradeService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errFiatBuySellTradeMissingRequired {
			t.Fatalf("error = %v, want %v", err, errFiatBuySellTradeMissingRequired)
		}
	})

	t.Run("signed_request_and_body", func(t *testing.T) {
		ts := sign.TimestampISO8601Millis(fixedNow)
		requestPath := "/api/v5/fiat/buy-sell/trade"
		wantBody := `{"clOrdId":"123456","side":"buy","fromCcy":"USD","toCcy":"BTC","rfqAmt":"30","rfqCcy":"USD","paymentMethod":"balance","quoteId":"qid"}`
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(ts, http.MethodPost, requestPath, wantBody))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/fiat/buy-sell/trade"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}

			bodyBytes, _ := io.ReadAll(r.Body)
			if got := string(bodyBytes); got != wantBody {
				t.Fatalf("body = %q, want %q", got, wantBody)
			}

			if got, want := r.Header.Get("OK-ACCESS-SIGN"), wantSig; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"ordId":"1234","clOrdId":"123456","quoteId":"qid","state":"completed","side":"buy","fromCcy":"USD","toCcy":"BTC","rfqAmt":"30","rfqCcy":"USD","fillPx":"2932.4","fillQuoteCcy":"USD","fillFromAmt":"30","fillToAmt":"0.01","cTime":"1646188510461"}]}`))
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

		got, err := c.NewFiatBuySellTradeService().
			ClOrdId("123456").
			Side("buy").
			FromCcy("USD").
			ToCcy("BTC").
			RfqAmt("30").
			RfqCcy("USD").
			PaymentMethod("balance").
			QuoteId("qid").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.OrdId != "1234" {
			t.Fatalf("OrdId = %q, want %q", got.OrdId, "1234")
		}
		if got.CTime != 1646188510461 {
			t.Fatalf("CTime = %d, want %d", got.CTime, 1646188510461)
		}
	})
}

func TestFiatBuySellHistoryService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	t.Run("signed_request", func(t *testing.T) {
		ts := sign.TimestampISO8601Millis(fixedNow)
		requestPath := "/api/v5/fiat/buy-sell/history?begin=1&end=2&limit=50&ordId=123&state=completed"
		wantSig := sign.SignHMACSHA256Base64("mysecret", sign.PrehashREST(ts, http.MethodGet, requestPath, ""))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/fiat/buy-sell/history"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "begin=1&end=2&limit=50&ordId=123&state=completed"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}

			if got, want := r.Header.Get("OK-ACCESS-SIGN"), wantSig; got != want {
				t.Fatalf("OK-ACCESS-SIGN = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"ordId":"1234","clOrdId":"123456","quoteId":"qid","state":"completed","side":"buy","fromCcy":"USD","toCcy":"BTC","rfqAmt":"30","rfqCcy":"USD","fillPx":"2932.4","fillQuoteCcy":"USD","fillFromAmt":"30","fillToAmt":"0.01","cTime":"1646188510461","uTime":"1646188510462"}]}`))
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

		got, err := c.NewFiatBuySellHistoryService().
			OrdId("123").
			State("completed").
			Begin("1").
			End("2").
			Limit(50).
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("len(data) = %d, want %d", len(got), 1)
		}
		if got[0].UTime != 1646188510462 {
			t.Fatalf("UTime = %d, want %d", got[0].UTime, 1646188510462)
		}
	})
}
