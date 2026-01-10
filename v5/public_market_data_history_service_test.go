package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublicMarketDataHistoryService_Do(t *testing.T) {
	t.Run("missing_required", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewPublicMarketDataHistoryService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errPublicMarketDataHistoryMissingRequired {
			t.Fatalf("error = %v, want %v", err, errPublicMarketDataHistoryMissingRequired)
		}
	})

	t.Run("missing_instIdList_for_spot", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewPublicMarketDataHistoryService().
			Module("1").
			InstType("SPOT").
			DateAggrType("daily").
			Begin("1").
			End("2").
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errPublicMarketDataHistoryMissingInstIdList {
			t.Fatalf("error = %v, want %v", err, errPublicMarketDataHistoryMissingInstIdList)
		}
	})

	t.Run("missing_instFamilyList_for_non_spot", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewPublicMarketDataHistoryService().
			Module("1").
			InstType("SWAP").
			DateAggrType("daily").
			Begin("1").
			End("2").
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errPublicMarketDataHistoryMissingInstFamilyList {
			t.Fatalf("error = %v, want %v", err, errPublicMarketDataHistoryMissingInstFamilyList)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/public/market-data-history"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			wantQuery := "begin=1&dateAggrType=daily&end=2&instFamilyList=BTC-USDT&instType=SWAP&module=1"
			if got, want := r.URL.RawQuery, wantQuery; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"dateAggrType":"daily","details":[{"dateRangeEnd":"1756656000000","dateRangeStart":"1756569600000","groupDetails":[{"dateTs":"1756656000000","filename":"BTC-USDT-SWAP-trades-2025-09-01.zip","sizeMB":"10.82","url":"https://static.okx.com/cdn/okex/traderecords/trades/daily/20250901/BTC-USDT-SWAP-trades-2025-09-01.zip"}],"groupSizeMB":"10.82","instFamily":"BTC-USDT","instId":"","instType":"SWAP"}],"totalSizeMB":"10.82","ts":"1756882260390"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewPublicMarketDataHistoryService().
			Module("1").
			InstType("SWAP").
			InstFamilyList("BTC-USDT").
			DateAggrType("daily").
			Begin("1").
			End("2").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].DateAggrType != "daily" || got[0].TotalSizeMB != "10.82" || got[0].TS != 1756882260390 {
			t.Fatalf("data = %#v", got)
		}
		if len(got[0].Details) != 1 || got[0].Details[0].InstFamily != "BTC-USDT" || got[0].Details[0].DateRangeStart != 1756569600000 {
			t.Fatalf("details = %#v", got[0])
		}
		if len(got[0].Details[0].GroupDetails) != 1 || got[0].Details[0].GroupDetails[0].Filename == "" || got[0].Details[0].GroupDetails[0].URL == "" {
			t.Fatalf("files = %#v", got[0].Details[0].GroupDetails)
		}
	})
}
