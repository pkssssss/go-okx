package okx

import (
	"context"
	"net/http"
	"net/url"
)

// FinanceStakingDefiOrdersActiveService 查看链上赚币活跃订单。
type FinanceStakingDefiOrdersActiveService struct {
	c *Client

	productId    string
	protocolType string
	ccy          string
	state        string
}

// NewFinanceStakingDefiOrdersActiveService 创建 FinanceStakingDefiOrdersActiveService。
func (c *Client) NewFinanceStakingDefiOrdersActiveService() *FinanceStakingDefiOrdersActiveService {
	return &FinanceStakingDefiOrdersActiveService{c: c}
}

// ProductId 设置项目 ID（可选）。
func (s *FinanceStakingDefiOrdersActiveService) ProductId(productId string) *FinanceStakingDefiOrdersActiveService {
	s.productId = productId
	return s
}

// ProtocolType 设置项目类型（可选；defi）。
func (s *FinanceStakingDefiOrdersActiveService) ProtocolType(protocolType string) *FinanceStakingDefiOrdersActiveService {
	s.protocolType = protocolType
	return s
}

// Ccy 设置投资币种（可选）。
func (s *FinanceStakingDefiOrdersActiveService) Ccy(ccy string) *FinanceStakingDefiOrdersActiveService {
	s.ccy = ccy
	return s
}

// State 设置订单状态（可选，见 OKX 文档）。
func (s *FinanceStakingDefiOrdersActiveService) State(state string) *FinanceStakingDefiOrdersActiveService {
	s.state = state
	return s
}

// Do 查看链上赚币活跃订单（GET /api/v5/finance/staking-defi/orders-active）。
func (s *FinanceStakingDefiOrdersActiveService) Do(ctx context.Context) ([]FinanceStakingDefiOrder, error) {
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
	if s.state != "" {
		q.Set("state", s.state)
	}
	if len(q) == 0 {
		q = nil
	}

	var data []FinanceStakingDefiOrder
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/finance/staking-defi/orders-active", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
