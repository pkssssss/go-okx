package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// FinanceFlexibleLoanLoanHistoryService 获取借贷历史。
type FinanceFlexibleLoanLoanHistoryService struct {
	c *Client

	typ    string
	after  string
	before string
	limit  *int
}

// NewFinanceFlexibleLoanLoanHistoryService 创建 FinanceFlexibleLoanLoanHistoryService。
func (c *Client) NewFinanceFlexibleLoanLoanHistoryService() *FinanceFlexibleLoanLoanHistoryService {
	return &FinanceFlexibleLoanLoanHistoryService{c: c}
}

// Type 设置操作类型（可选，见 OKX 文档）。
func (s *FinanceFlexibleLoanLoanHistoryService) Type(typ string) *FinanceFlexibleLoanLoanHistoryService {
	s.typ = typ
	return s
}

// After 请求此 ID 之前（更旧的数据）的分页内容（传 refId，不包含）。
func (s *FinanceFlexibleLoanLoanHistoryService) After(after string) *FinanceFlexibleLoanLoanHistoryService {
	s.after = after
	return s
}

// Before 请求此 ID 之后（更新的数据）的分页内容（传 refId，不包含）。
func (s *FinanceFlexibleLoanLoanHistoryService) Before(before string) *FinanceFlexibleLoanLoanHistoryService {
	s.before = before
	return s
}

// Limit 设置返回数量（最大 100，默认 100）。
func (s *FinanceFlexibleLoanLoanHistoryService) Limit(limit int) *FinanceFlexibleLoanLoanHistoryService {
	s.limit = &limit
	return s
}

// Do 获取借贷历史（GET /api/v5/finance/flexible-loan/loan-history）。
func (s *FinanceFlexibleLoanLoanHistoryService) Do(ctx context.Context) ([]FinanceFlexibleLoanLoanHistory, error) {
	q := url.Values{}
	if s.typ != "" {
		q.Set("type", s.typ)
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

	var data []FinanceFlexibleLoanLoanHistory
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/finance/flexible-loan/loan-history", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
