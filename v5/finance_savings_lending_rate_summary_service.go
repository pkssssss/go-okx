package okx

import (
	"context"
	"net/http"
	"net/url"
)

// FinanceSavingsLendingRateSummaryService 获取市场借贷信息（公共）。
type FinanceSavingsLendingRateSummaryService struct {
	c   *Client
	ccy string
}

// NewFinanceSavingsLendingRateSummaryService 创建 FinanceSavingsLendingRateSummaryService。
func (c *Client) NewFinanceSavingsLendingRateSummaryService() *FinanceSavingsLendingRateSummaryService {
	return &FinanceSavingsLendingRateSummaryService{c: c}
}

// Ccy 设置币种过滤（可选）。
func (s *FinanceSavingsLendingRateSummaryService) Ccy(ccy string) *FinanceSavingsLendingRateSummaryService {
	s.ccy = ccy
	return s
}

// Do 获取市场借贷信息（GET /api/v5/finance/savings/lending-rate-summary）。
func (s *FinanceSavingsLendingRateSummaryService) Do(ctx context.Context) ([]FinanceSavingsLendingRateSummary, error) {
	var q url.Values
	if s.ccy != "" {
		q = url.Values{}
		q.Set("ccy", s.ccy)
	}

	var data []FinanceSavingsLendingRateSummary
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/finance/savings/lending-rate-summary", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
