package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublicConvertContractCoinService_Do(t *testing.T) {
	t.Run("missing_required", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewPublicConvertContractCoinService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errPublicConvertContractCoinMissingRequired {
			t.Fatalf("error = %v, want %v", err, errPublicConvertContractCoinMissingRequired)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/public/convert-contract-coin"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			wantQuery := "instId=BTC-USD-SWAP&opType=open&px=35000&sz=0.888&type=1&unit=coin"
			if got, want := r.URL.RawQuery, wantQuery; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instId":"BTC-USD-SWAP","px":"35000","sz":"311","type":"1","unit":"coin"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewPublicConvertContractCoinService().
			InstId("BTC-USD-SWAP").
			Type("1").
			Sz("0.888").
			Px("35000").
			Unit("coin").
			OpType("open").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].InstId != "BTC-USD-SWAP" || got[0].Sz != "311" || got[0].Type != "1" || got[0].Unit != "coin" {
			t.Fatalf("data = %#v", got)
		}
	})
}
