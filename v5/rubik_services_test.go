package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRubikSupportCoinService_Do(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/rubik/stat/trading-data/support-coin"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, ""; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":{"contract":["BTC"],"option":["BTC"],"spot":["BTC"]}}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewRubikSupportCoinService().Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got.Contract) != 1 || got.Contract[0] != "BTC" {
			t.Fatalf("contract = %#v", got.Contract)
		}
	})
}

func TestRubikOpenInterestHistoryService_Do(t *testing.T) {
	t.Run("missing_instId", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewRubikOpenInterestHistoryService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errRubikOpenInterestHistoryMissingInstId {
			t.Fatalf("error = %v, want %v", err, errRubikOpenInterestHistoryMissingInstId)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.URL.Path, "/api/v5/rubik/stat/contracts/open-interest-history"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "begin=1&end=2&instId=BTC-USDT-SWAP&limit=50&period=1D"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[["1701417600000","731377.57500501","111","8888888"]]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewRubikOpenInterestHistoryService().
			InstId("BTC-USDT-SWAP").
			Period("1D").
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
		if got[0].TS != 1701417600000 {
			t.Fatalf("TS = %d, want %d", got[0].TS, 1701417600000)
		}
		if got[0].OIUsd != "8888888" {
			t.Fatalf("OIUsd = %q, want %q", got[0].OIUsd, "8888888")
		}
	})
}

func TestRubikTakerVolumeService_Do(t *testing.T) {
	t.Run("missing_required", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewRubikTakerVolumeService().Ccy("BTC").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errRubikTakerVolumeMissingRequired {
			t.Fatalf("error = %v, want %v", err, errRubikTakerVolumeMissingRequired)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.URL.Path, "/api/v5/rubik/stat/taker-volume"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "begin=1&ccy=BTC&end=2&instType=SPOT&period=1D"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[["1630425600000","7596.2651","7149.4855"]]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewRubikTakerVolumeService().
			Ccy("BTC").
			InstType("SPOT").
			Begin("1").
			End("2").
			Period("1D").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].TS != 1630425600000 {
			t.Fatalf("data = %#v", got)
		}
		if got[0].SellVol != "7596.2651" {
			t.Fatalf("SellVol = %q, want %q", got[0].SellVol, "7596.2651")
		}
	})
}

func TestRubikTakerVolumeContractService_Do(t *testing.T) {
	t.Run("missing_instId", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewRubikTakerVolumeContractService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errRubikTakerVolumeContractMissingInstId {
			t.Fatalf("error = %v, want %v", err, errRubikTakerVolumeContractMissingInstId)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.URL.Path, "/api/v5/rubik/stat/taker-volume-contract"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "begin=1&end=2&instId=BTC-USDT-SWAP&limit=50&period=1D&unit=2"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[["1701417600000","200","380"]]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewRubikTakerVolumeContractService().
			InstId("BTC-USDT-SWAP").
			Period("1D").
			Unit("2").
			Begin("1").
			End("2").
			Limit(50).
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].BuyVol != "380" {
			t.Fatalf("data = %#v", got)
		}
	})
}

func TestRubikMarginLoanRatioService_Do(t *testing.T) {
	t.Run("missing_ccy", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewRubikMarginLoanRatioService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errRubikMarginLoanRatioMissingCcy {
			t.Fatalf("error = %v, want %v", err, errRubikMarginLoanRatioMissingCcy)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.URL.Path, "/api/v5/rubik/stat/margin/loan-ratio"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "begin=1&ccy=BTC&end=2&period=1D"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[["1630492800000","0.4614"]]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewRubikMarginLoanRatioService().
			Ccy("BTC").
			Begin("1").
			End("2").
			Period("1D").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].Ratio != "0.4614" {
			t.Fatalf("data = %#v", got)
		}
	})
}

func TestRubikLongShortAccountRatioContractTopTraderService_Do(t *testing.T) {
	t.Run("missing_instId", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewRubikLongShortAccountRatioContractTopTraderService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errRubikLongShortAccountRatioContractTopTraderMissingInstId {
			t.Fatalf("error = %v, want %v", err, errRubikLongShortAccountRatioContractTopTraderMissingInstId)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.URL.Path, "/api/v5/rubik/stat/contracts/long-short-account-ratio-contract-top-trader"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "begin=1&end=2&instId=BTC-USDT-SWAP&limit=50&period=1D"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[["1701417600000","1.1739"]]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewRubikLongShortAccountRatioContractTopTraderService().
			InstId("BTC-USDT-SWAP").
			Period("1D").
			Begin("1").
			End("2").
			Limit(50).
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].Ratio != "1.1739" {
			t.Fatalf("data = %#v", got)
		}
	})
}

func TestRubikLongShortPositionRatioContractTopTraderService_Do(t *testing.T) {
	t.Run("missing_instId", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewRubikLongShortPositionRatioContractTopTraderService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errRubikLongShortPositionRatioContractTopTraderMissingInstId {
			t.Fatalf("error = %v, want %v", err, errRubikLongShortPositionRatioContractTopTraderMissingInstId)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.URL.Path, "/api/v5/rubik/stat/contracts/long-short-position-ratio-contract-top-trader"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "begin=1&end=2&instId=BTC-USDT-SWAP&limit=50&period=1D"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[["1701417600000","0.1236"]]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewRubikLongShortPositionRatioContractTopTraderService().
			InstId("BTC-USDT-SWAP").
			Period("1D").
			Begin("1").
			End("2").
			Limit(50).
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].Ratio != "0.1236" {
			t.Fatalf("data = %#v", got)
		}
	})
}

func TestRubikLongShortAccountRatioContractService_Do(t *testing.T) {
	t.Run("missing_instId", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewRubikLongShortAccountRatioContractService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errRubikLongShortAccountRatioContractMissingInstId {
			t.Fatalf("error = %v, want %v", err, errRubikLongShortAccountRatioContractMissingInstId)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.URL.Path, "/api/v5/rubik/stat/contracts/long-short-account-ratio-contract"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "begin=1&end=2&instId=BTC-USDT-SWAP&limit=50&period=1D"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[["1701417600000","1.25"]]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewRubikLongShortAccountRatioContractService().
			InstId("BTC-USDT-SWAP").
			Period("1D").
			Begin("1").
			End("2").
			Limit(50).
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].Ratio != "1.25" {
			t.Fatalf("data = %#v", got)
		}
	})
}

func TestRubikLongShortAccountRatioService_Do(t *testing.T) {
	t.Run("missing_ccy", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewRubikLongShortAccountRatioService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errRubikLongShortAccountRatioMissingCcy {
			t.Fatalf("error = %v, want %v", err, errRubikLongShortAccountRatioMissingCcy)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.URL.Path, "/api/v5/rubik/stat/contracts/long-short-account-ratio"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "begin=1&ccy=BTC&end=2&period=1D"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[["1630502100000","1.25"]]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewRubikLongShortAccountRatioService().
			Ccy("BTC").
			Begin("1").
			End("2").
			Period("1D").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].Ratio != "1.25" {
			t.Fatalf("data = %#v", got)
		}
	})
}

func TestRubikContractsOpenInterestVolumeService_Do(t *testing.T) {
	t.Run("missing_ccy", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewRubikContractsOpenInterestVolumeService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errRubikContractsOpenInterestVolumeMissingCcy {
			t.Fatalf("error = %v, want %v", err, errRubikContractsOpenInterestVolumeMissingCcy)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.URL.Path, "/api/v5/rubik/stat/contracts/open-interest-volume"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "begin=1&ccy=BTC&end=2&period=1D"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[["1630502400000","1713028741.6898","39800873.554"]]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewRubikContractsOpenInterestVolumeService().
			Ccy("BTC").
			Begin("1").
			End("2").
			Period("1D").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].OI != "1713028741.6898" {
			t.Fatalf("data = %#v", got)
		}
	})
}

func TestRubikOptionOpenInterestVolumeService_Do(t *testing.T) {
	t.Run("missing_ccy", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewRubikOptionOpenInterestVolumeService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errRubikOptionOpenInterestVolumeMissingCcy {
			t.Fatalf("error = %v, want %v", err, errRubikOptionOpenInterestVolumeMissingCcy)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.URL.Path, "/api/v5/rubik/stat/option/open-interest-volume"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "ccy=BTC&period=8H"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[["1630368000000","3458.1000","78.8000"]]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewRubikOptionOpenInterestVolumeService().Ccy("BTC").Period("8H").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].TS != 1630368000000 {
			t.Fatalf("data = %#v", got)
		}
	})
}

func TestRubikOptionOpenInterestVolumeRatioService_Do(t *testing.T) {
	t.Run("missing_ccy", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewRubikOptionOpenInterestVolumeRatioService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errRubikOptionOpenInterestVolumeRatioMissingCcy {
			t.Fatalf("error = %v, want %v", err, errRubikOptionOpenInterestVolumeRatioMissingCcy)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.URL.Path, "/api/v5/rubik/stat/option/open-interest-volume-ratio"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "ccy=BTC&period=8H"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[["1630512000000","2.7261","2.3447"]]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewRubikOptionOpenInterestVolumeRatioService().Ccy("BTC").Period("8H").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].VolRatio != "2.3447" {
			t.Fatalf("data = %#v", got)
		}
	})
}

func TestRubikOptionOpenInterestVolumeExpiryService_Do(t *testing.T) {
	t.Run("missing_ccy", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewRubikOptionOpenInterestVolumeExpiryService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errRubikOptionOpenInterestVolumeExpiryMissingCcy {
			t.Fatalf("error = %v, want %v", err, errRubikOptionOpenInterestVolumeExpiryMissingCcy)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.URL.Path, "/api/v5/rubik/stat/option/open-interest-volume-expiry"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "ccy=BTC&period=8H"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[["1630540800000","20210902","6.4","18.4","0.7","0.4"]]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewRubikOptionOpenInterestVolumeExpiryService().Ccy("BTC").Period("8H").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].ExpTime != "20210902" {
			t.Fatalf("data = %#v", got)
		}
	})
}

func TestRubikOptionOpenInterestVolumeStrikeService_Do(t *testing.T) {
	t.Run("missing_required", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewRubikOptionOpenInterestVolumeStrikeService().Ccy("BTC").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errRubikOptionOpenInterestVolumeStrikeMissingRequired {
			t.Fatalf("error = %v, want %v", err, errRubikOptionOpenInterestVolumeStrikeMissingRequired)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.URL.Path, "/api/v5/rubik/stat/option/open-interest-volume-strike"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "ccy=BTC&expTime=20210901&period=8H"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[["1630540800000","10000","0","0.5","0","0"]]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewRubikOptionOpenInterestVolumeStrikeService().Ccy("BTC").ExpTime("20210901").Period("8H").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].Strike != "10000" {
			t.Fatalf("data = %#v", got)
		}
	})
}

func TestRubikOptionTakerBlockVolumeService_Do(t *testing.T) {
	t.Run("missing_ccy", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewRubikOptionTakerBlockVolumeService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errRubikOptionTakerBlockVolumeMissingCcy {
			t.Fatalf("error = %v, want %v", err, errRubikOptionTakerBlockVolumeMissingCcy)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.URL.Path, "/api/v5/rubik/stat/option/taker-block-volume"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, "ccy=BTC&period=8H"; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":["1630512000000","8.55","67.3","16.05","16.3","126.4","40.7"]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewRubikOptionTakerBlockVolumeService().Ccy("BTC").Period("8H").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if got.TS != 1630512000000 || got.CallBuyVol != "8.55" || got.PutBlockVol != "40.7" {
			t.Fatalf("data = %#v", got)
		}
	})
}
