package okx

import (
	"context"
	"errors"
	"net/http"
)

type financeStakingDefiRedeemRequest struct {
	OrdId            string `json:"ordId"`
	ProtocolType     string `json:"protocolType"`
	AllowEarlyRedeem *bool  `json:"allowEarlyRedeem,omitempty"`
}

// FinanceStakingDefiRedeemService 赎回链上赚币项目。
type FinanceStakingDefiRedeemService struct {
	c   *Client
	req financeStakingDefiRedeemRequest
}

// NewFinanceStakingDefiRedeemService 创建 FinanceStakingDefiRedeemService。
func (c *Client) NewFinanceStakingDefiRedeemService() *FinanceStakingDefiRedeemService {
	return &FinanceStakingDefiRedeemService{c: c}
}

// OrdId 设置订单 ID（必填）。
func (s *FinanceStakingDefiRedeemService) OrdId(ordId string) *FinanceStakingDefiRedeemService {
	s.req.OrdId = ordId
	return s
}

// ProtocolType 设置项目类型（必填；defi）。
func (s *FinanceStakingDefiRedeemService) ProtocolType(protocolType string) *FinanceStakingDefiRedeemService {
	s.req.ProtocolType = protocolType
	return s
}

// AllowEarlyRedeem 设置是否提前赎回（可选，默认 false）。
func (s *FinanceStakingDefiRedeemService) AllowEarlyRedeem(allow bool) *FinanceStakingDefiRedeemService {
	s.req.AllowEarlyRedeem = &allow
	return s
}

var (
	errFinanceStakingDefiRedeemMissingRequired = errors.New("okx: staking-defi redeem requires ordId and protocolType")
	errEmptyFinanceStakingDefiRedeemAck        = errors.New("okx: empty staking-defi redeem response")
)

// Do 赎回链上赚币项目（POST /api/v5/finance/staking-defi/redeem）。
func (s *FinanceStakingDefiRedeemService) Do(ctx context.Context) (*FinanceStakingDefiOrderAck, error) {
	if s.req.OrdId == "" || s.req.ProtocolType == "" {
		return nil, errFinanceStakingDefiRedeemMissingRequired
	}

	var data []FinanceStakingDefiOrderAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/finance/staking-defi/redeem", nil, s.req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyFinanceStakingDefiRedeemAck
	}
	return &data[0], nil
}
