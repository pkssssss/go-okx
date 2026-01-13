package okx

import (
	"context"
	"net/http"
	"net/url"
)

// FinanceSavingsBalanceService 获取活期简单赚币余额。
type FinanceSavingsBalanceService struct {
	c   *Client
	ccy string
}

// NewFinanceSavingsBalanceService 创建 FinanceSavingsBalanceService。
func (c *Client) NewFinanceSavingsBalanceService() *FinanceSavingsBalanceService {
	return &FinanceSavingsBalanceService{c: c}
}

// Ccy 设置币种过滤（可选）。
func (s *FinanceSavingsBalanceService) Ccy(ccy string) *FinanceSavingsBalanceService {
	s.ccy = ccy
	return s
}

// Do 获取活期简单赚币余额（GET /api/v5/finance/savings/balance）。
func (s *FinanceSavingsBalanceService) Do(ctx context.Context) ([]FinanceSavingsBalance, error) {
	var q url.Values
	if s.ccy != "" {
		q = url.Values{}
		q.Set("ccy", s.ccy)
	}

	var data []FinanceSavingsBalance
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/finance/savings/balance", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
