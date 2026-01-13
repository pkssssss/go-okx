package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// FinanceFlexibleLoanInterestAccruedService 获取计息记录。
type FinanceFlexibleLoanInterestAccruedService struct {
	c *Client

	ccy    string
	after  string
	before string
	limit  *int
}

// NewFinanceFlexibleLoanInterestAccruedService 创建 FinanceFlexibleLoanInterestAccruedService。
func (c *Client) NewFinanceFlexibleLoanInterestAccruedService() *FinanceFlexibleLoanInterestAccruedService {
	return &FinanceFlexibleLoanInterestAccruedService{c: c}
}

// Ccy 设置借贷币种过滤（可选）。
func (s *FinanceFlexibleLoanInterestAccruedService) Ccy(ccy string) *FinanceFlexibleLoanInterestAccruedService {
	s.ccy = ccy
	return s
}

// After 请求此 ID 之前（更旧的数据）的分页内容（传 refId，不包含）。
func (s *FinanceFlexibleLoanInterestAccruedService) After(after string) *FinanceFlexibleLoanInterestAccruedService {
	s.after = after
	return s
}

// Before 请求此 ID 之后（更新的数据）的分页内容（传 refId，不包含）。
func (s *FinanceFlexibleLoanInterestAccruedService) Before(before string) *FinanceFlexibleLoanInterestAccruedService {
	s.before = before
	return s
}

// Limit 设置返回数量（最大 100，默认 100）。
func (s *FinanceFlexibleLoanInterestAccruedService) Limit(limit int) *FinanceFlexibleLoanInterestAccruedService {
	s.limit = &limit
	return s
}

// Do 获取计息记录（GET /api/v5/finance/flexible-loan/interest-accrued）。
func (s *FinanceFlexibleLoanInterestAccruedService) Do(ctx context.Context) ([]FinanceFlexibleLoanInterestAccrued, error) {
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

	var data []FinanceFlexibleLoanInterestAccrued
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/finance/flexible-loan/interest-accrued", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
