package okx

import "net/http"

func handleTradeAccountRateLimitMock(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodGet || r.URL.Path != "/api/v5/trade/account-rate-limit" {
		return false
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"accRateLimit":"1","fillRatio":"0","mainFillRatio":"0","nextAccRateLimit":"1","ts":"1"}]}`))
	return true
}
