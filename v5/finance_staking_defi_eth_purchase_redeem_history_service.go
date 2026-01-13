package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// FinanceStakingDefiETHPurchaseRedeemHistoryService 获取 ETH 质押申购/赎回记录。
type FinanceStakingDefiETHPurchaseRedeemHistoryService struct {
	c *Client

	typ    string
	status string
	after  string
	before string
	limit  *int
}

// NewFinanceStakingDefiETHPurchaseRedeemHistoryService 创建 FinanceStakingDefiETHPurchaseRedeemHistoryService。
func (c *Client) NewFinanceStakingDefiETHPurchaseRedeemHistoryService() *FinanceStakingDefiETHPurchaseRedeemHistoryService {
	return &FinanceStakingDefiETHPurchaseRedeemHistoryService{c: c}
}

// Type 设置类型（可选：purchase/redeem）。
func (s *FinanceStakingDefiETHPurchaseRedeemHistoryService) Type(typ string) *FinanceStakingDefiETHPurchaseRedeemHistoryService {
	s.typ = typ
	return s
}

// Status 设置状态（可选，见 OKX 文档）。
func (s *FinanceStakingDefiETHPurchaseRedeemHistoryService) Status(status string) *FinanceStakingDefiETHPurchaseRedeemHistoryService {
	s.status = status
	return s
}

// After 请求此 requestTime 之前（更旧的数据）的分页内容（时间戳毫秒字符串）。
func (s *FinanceStakingDefiETHPurchaseRedeemHistoryService) After(after string) *FinanceStakingDefiETHPurchaseRedeemHistoryService {
	s.after = after
	return s
}

// Before 请求此 requestTime 之后（更新的数据）的分页内容（时间戳毫秒字符串）。
func (s *FinanceStakingDefiETHPurchaseRedeemHistoryService) Before(before string) *FinanceStakingDefiETHPurchaseRedeemHistoryService {
	s.before = before
	return s
}

// Limit 设置返回数量（最大 100，默认 100）。
func (s *FinanceStakingDefiETHPurchaseRedeemHistoryService) Limit(limit int) *FinanceStakingDefiETHPurchaseRedeemHistoryService {
	s.limit = &limit
	return s
}

// Do 获取 ETH 质押申购/赎回记录（GET /api/v5/finance/staking-defi/eth/purchase-redeem-history）。
func (s *FinanceStakingDefiETHPurchaseRedeemHistoryService) Do(ctx context.Context) ([]FinanceStakingDefiPurchaseRedeemHistory, error) {
	q := url.Values{}
	if s.typ != "" {
		q.Set("type", s.typ)
	}
	if s.status != "" {
		q.Set("status", s.status)
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

	var data []FinanceStakingDefiPurchaseRedeemHistory
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/finance/staking-defi/eth/purchase-redeem-history", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
