package okx

import (
	"context"
	"errors"
	"net/http"
)

type financeStakingDefiETHCancelRedeemRequest struct {
	OrdId string `json:"ordId"`
}

// FinanceStakingDefiETHCancelRedeemService 撤销 ETH 质押赎回。
type FinanceStakingDefiETHCancelRedeemService struct {
	c     *Client
	ordId string
}

// NewFinanceStakingDefiETHCancelRedeemService 创建 FinanceStakingDefiETHCancelRedeemService。
func (c *Client) NewFinanceStakingDefiETHCancelRedeemService() *FinanceStakingDefiETHCancelRedeemService {
	return &FinanceStakingDefiETHCancelRedeemService{c: c}
}

// OrdId 设置订单 ID（必填）。
func (s *FinanceStakingDefiETHCancelRedeemService) OrdId(ordId string) *FinanceStakingDefiETHCancelRedeemService {
	s.ordId = ordId
	return s
}

var (
	errFinanceStakingDefiETHCancelRedeemMissingOrdId = errors.New("okx: staking-defi eth cancel-redeem requires ordId")
	errEmptyFinanceStakingDefiETHCancelRedeemAck     = errors.New("okx: empty staking-defi eth cancel-redeem response")
)

// Do 撤销 ETH 质押赎回（POST /api/v5/finance/staking-defi/eth/cancel-redeem）。
func (s *FinanceStakingDefiETHCancelRedeemService) Do(ctx context.Context) (*FinanceStakingDefiOrderAck, error) {
	if s.ordId == "" {
		return nil, errFinanceStakingDefiETHCancelRedeemMissingOrdId
	}

	req := financeStakingDefiETHCancelRedeemRequest{OrdId: s.ordId}
	var data []FinanceStakingDefiOrderAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/finance/staking-defi/eth/cancel-redeem", nil, req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyFinanceStakingDefiETHCancelRedeemAck
	}
	return &data[0], nil
}
