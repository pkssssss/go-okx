package okx

import (
	"context"
	"errors"
	"net/http"
)

type financeStakingDefiETHRedeemRequest struct {
	Amt string `json:"amt"`
}

// FinanceStakingDefiETHRedeemService 赎回 ETH 质押（赎回 BETH）。
type FinanceStakingDefiETHRedeemService struct {
	c   *Client
	amt string
}

// NewFinanceStakingDefiETHRedeemService 创建 FinanceStakingDefiETHRedeemService。
func (c *Client) NewFinanceStakingDefiETHRedeemService() *FinanceStakingDefiETHRedeemService {
	return &FinanceStakingDefiETHRedeemService{c: c}
}

// Amt 设置赎回数量（必填，字符串）。
func (s *FinanceStakingDefiETHRedeemService) Amt(amt string) *FinanceStakingDefiETHRedeemService {
	s.amt = amt
	return s
}

var errFinanceStakingDefiETHRedeemMissingAmt = errors.New("okx: staking-defi eth redeem requires amt")

// Do 赎回 ETH 质押（POST /api/v5/finance/staking-defi/eth/redeem）。
//
// 注意：OKX 返回的 data 为空数组（code=0 代表请求已被成功处理）。
func (s *FinanceStakingDefiETHRedeemService) Do(ctx context.Context) error {
	if s.amt == "" {
		return errFinanceStakingDefiETHRedeemMissingAmt
	}

	req := financeStakingDefiETHRedeemRequest{Amt: s.amt}
	return s.c.do(ctx, http.MethodPost, "/api/v5/finance/staking-defi/eth/redeem", nil, req, true, nil)
}
