package okx

import (
	"context"
	"net/http"
	"net/url"
)

// FinanceStakingDefiOffersService 查看链上赚币项目。
type FinanceStakingDefiOffersService struct {
	c *Client

	productId    string
	protocolType string
	ccy          string
}

// NewFinanceStakingDefiOffersService 创建 FinanceStakingDefiOffersService。
func (c *Client) NewFinanceStakingDefiOffersService() *FinanceStakingDefiOffersService {
	return &FinanceStakingDefiOffersService{c: c}
}

// ProductId 设置项目 ID（可选）。
func (s *FinanceStakingDefiOffersService) ProductId(productId string) *FinanceStakingDefiOffersService {
	s.productId = productId
	return s
}

// ProtocolType 设置项目类型（可选；defi）。
func (s *FinanceStakingDefiOffersService) ProtocolType(protocolType string) *FinanceStakingDefiOffersService {
	s.protocolType = protocolType
	return s
}

// Ccy 设置投资币种（可选）。
func (s *FinanceStakingDefiOffersService) Ccy(ccy string) *FinanceStakingDefiOffersService {
	s.ccy = ccy
	return s
}

// Do 查看链上赚币项目（GET /api/v5/finance/staking-defi/offers）。
func (s *FinanceStakingDefiOffersService) Do(ctx context.Context) ([]FinanceStakingDefiOffer, error) {
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
	if len(q) == 0 {
		q = nil
	}

	var data []FinanceStakingDefiOffer
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/finance/staking-defi/offers", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
