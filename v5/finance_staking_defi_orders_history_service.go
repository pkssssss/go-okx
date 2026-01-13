package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// FinanceStakingDefiOrdersHistoryService 查看链上赚币历史订单。
type FinanceStakingDefiOrdersHistoryService struct {
	c *Client

	productId    string
	protocolType string
	ccy          string

	after  string
	before string
	limit  *int
}

// NewFinanceStakingDefiOrdersHistoryService 创建 FinanceStakingDefiOrdersHistoryService。
func (c *Client) NewFinanceStakingDefiOrdersHistoryService() *FinanceStakingDefiOrdersHistoryService {
	return &FinanceStakingDefiOrdersHistoryService{c: c}
}

// ProductId 设置项目 ID（可选）。
func (s *FinanceStakingDefiOrdersHistoryService) ProductId(productId string) *FinanceStakingDefiOrdersHistoryService {
	s.productId = productId
	return s
}

// ProtocolType 设置项目类型（可选；defi）。
func (s *FinanceStakingDefiOrdersHistoryService) ProtocolType(protocolType string) *FinanceStakingDefiOrdersHistoryService {
	s.protocolType = protocolType
	return s
}

// Ccy 设置投资币种（可选）。
func (s *FinanceStakingDefiOrdersHistoryService) Ccy(ccy string) *FinanceStakingDefiOrdersHistoryService {
	s.ccy = ccy
	return s
}

// After 请求此 ID 之前（更旧的数据）的分页内容（传 ordId）。
func (s *FinanceStakingDefiOrdersHistoryService) After(after string) *FinanceStakingDefiOrdersHistoryService {
	s.after = after
	return s
}

// Before 请求此 ID 之后（更新的数据）的分页内容（传 ordId）。
func (s *FinanceStakingDefiOrdersHistoryService) Before(before string) *FinanceStakingDefiOrdersHistoryService {
	s.before = before
	return s
}

// Limit 设置返回数量（最大 100，默认 100）。
func (s *FinanceStakingDefiOrdersHistoryService) Limit(limit int) *FinanceStakingDefiOrdersHistoryService {
	s.limit = &limit
	return s
}

// Do 查看链上赚币历史订单（GET /api/v5/finance/staking-defi/orders-history）。
func (s *FinanceStakingDefiOrdersHistoryService) Do(ctx context.Context) ([]FinanceStakingDefiOrder, error) {
	q := url.Values{}
	if s.productId != "" {
		q.Set("productId", s.productId)
	}
	if s.protocolType != "" {
		q.Set("protocolType", s.protocolType)
	}
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

	var data []FinanceStakingDefiOrder
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/finance/staking-defi/orders-history", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
