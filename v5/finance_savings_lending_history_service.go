package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// FinanceSavingsLendingHistoryService 获取活期简单赚币出借明细（最近一个月）。
type FinanceSavingsLendingHistoryService struct {
	c *Client

	ccy    string
	after  string
	before string
	limit  *int
}

// NewFinanceSavingsLendingHistoryService 创建 FinanceSavingsLendingHistoryService。
func (c *Client) NewFinanceSavingsLendingHistoryService() *FinanceSavingsLendingHistoryService {
	return &FinanceSavingsLendingHistoryService{c: c}
}

// Ccy 设置币种过滤（可选）。
func (s *FinanceSavingsLendingHistoryService) Ccy(ccy string) *FinanceSavingsLendingHistoryService {
	s.ccy = ccy
	return s
}

// After 查询在此之前的内容（时间戳毫秒字符串）。
func (s *FinanceSavingsLendingHistoryService) After(after string) *FinanceSavingsLendingHistoryService {
	s.after = after
	return s
}

// Before 查询在此之后的内容（时间戳毫秒字符串）。
func (s *FinanceSavingsLendingHistoryService) Before(before string) *FinanceSavingsLendingHistoryService {
	s.before = before
	return s
}

// Limit 设置返回数量（最大 100，默认 100）。
func (s *FinanceSavingsLendingHistoryService) Limit(limit int) *FinanceSavingsLendingHistoryService {
	s.limit = &limit
	return s
}

// Do 获取活期简单赚币出借明细（GET /api/v5/finance/savings/lending-history）。
func (s *FinanceSavingsLendingHistoryService) Do(ctx context.Context) ([]FinanceSavingsLendingHistory, error) {
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

	var data []FinanceSavingsLendingHistory
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/finance/savings/lending-history", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
