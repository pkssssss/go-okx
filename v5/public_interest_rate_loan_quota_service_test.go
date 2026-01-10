package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublicInterestRateLoanQuotaService_Do(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/public/interest-rate-loan-quota"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got, want := r.URL.RawQuery, ""; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"configCcyList":[{"ccy":"USDT","rate":"0.00043728"}],"basic":[{"ccy":"USDT","quota":"500000","rate":"0.00043728"}],"vip":[{"level":"VIP1","loanQuotaCoef":"2","irDiscount":""}],"regular":[{"level":"Lv1","loanQuotaCoef":"1","irDiscount":""}],"config":[{"ccy":"USDT","stgyType":"0","quota":"500000","level":"VIP1"}]}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewPublicInterestRateLoanQuotaService().Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || len(got[0].Basic) != 1 || got[0].Basic[0].Ccy != "USDT" {
			t.Fatalf("data = %#v", got)
		}
		if got[0].Basic[0].Rate == "" || got[0].Basic[0].Quota == "" {
			t.Fatalf("basic = %#v", got[0].Basic[0])
		}
		if len(got[0].VIP) != 1 || got[0].VIP[0].LoanQuotaCoef != "2" {
			t.Fatalf("vip = %#v", got[0].VIP)
		}
		if len(got[0].ConfigCcyList) != 1 || got[0].ConfigCcyList[0].Rate == "" {
			t.Fatalf("configCcyList = %#v", got[0].ConfigCcyList)
		}
	})
}
