package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// FinanceSavingsLendingRateHistoryService 获取市场借贷历史（公共）。
type FinanceSavingsLendingRateHistoryService struct {
	c *Client

	ccy    string
	after  string
	before string
	limit  *int
}

// NewFinanceSavingsLendingRateHistoryService 创建 FinanceSavingsLendingRateHistoryService。
func (c *Client) NewFinanceSavingsLendingRateHistoryService() *FinanceSavingsLendingRateHistoryService {
	return &FinanceSavingsLendingRateHistoryService{c: c}
}

// Ccy 设置币种过滤（可选）。
func (s *FinanceSavingsLendingRateHistoryService) Ccy(ccy string) *FinanceSavingsLendingRateHistoryService {
	s.ccy = ccy
	return s
}

// After 查询在此之前的内容（时间戳毫秒字符串）。
func (s *FinanceSavingsLendingRateHistoryService) After(after string) *FinanceSavingsLendingRateHistoryService {
	s.after = after
	return s
}

// Before 查询在此之后的内容（时间戳毫秒字符串）。
func (s *FinanceSavingsLendingRateHistoryService) Before(before string) *FinanceSavingsLendingRateHistoryService {
	s.before = before
	return s
}

// Limit 设置返回数量（最大 100，默认 100）。
func (s *FinanceSavingsLendingRateHistoryService) Limit(limit int) *FinanceSavingsLendingRateHistoryService {
	s.limit = &limit
	return s
}

// Do 获取市场借贷历史（GET /api/v5/finance/savings/lending-rate-history）。
func (s *FinanceSavingsLendingRateHistoryService) Do(ctx context.Context) ([]FinanceSavingsLendingRateHistory, error) {
	q := url.Values{}
	if s.ccy != "" {
		q.Set("ccy", s.ccy)
	}
	if s.after != "" {
		q.Set("after", s.after)
	}
	if s.before != "" {
		q.Set("before", s.before)
	}
	if s.limit != nil {
		q.Set("limit", strconv.Itoa(*s.limit))
	}
	if len(q) == 0 {
		q = nil
	}

	var data []FinanceSavingsLendingRateHistory
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/finance/savings/lending-rate-history", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
