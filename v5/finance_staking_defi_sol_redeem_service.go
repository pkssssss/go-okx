package okx

import (
	"context"
	"errors"
	"net/http"
)

type financeStakingDefiSOLRedeemRequest struct {
	Amt string `json:"amt"`
}

// FinanceStakingDefiSOLRedeemService 赎回 SOL 质押（赎回 OKSOL）。
type FinanceStakingDefiSOLRedeemService struct {
	c   *Client
	amt string
}

// NewFinanceStakingDefiSOLRedeemService 创建 FinanceStakingDefiSOLRedeemService。
func (c *Client) NewFinanceStakingDefiSOLRedeemService() *FinanceStakingDefiSOLRedeemService {
	return &FinanceStakingDefiSOLRedeemService{c: c}
}

// Amt 设置赎回数量（必填，字符串）。
func (s *FinanceStakingDefiSOLRedeemService) Amt(amt string) *FinanceStakingDefiSOLRedeemService {
	s.amt = amt
	return s
}

var errFinanceStakingDefiSOLRedeemMissingAmt = errors.New("okx: staking-defi sol redeem requires amt")

// Do 赎回 SOL 质押（POST /api/v5/finance/staking-defi/sol/redeem）。
//
// 注意：OKX 返回的 data 为空数组（code=0 代表请求已被成功处理）。
func (s *FinanceStakingDefiSOLRedeemService) Do(ctx context.Context) error {
	if s.amt == "" {
		return errFinanceStakingDefiSOLRedeemMissingAmt
	}

	req := financeStakingDefiSOLRedeemRequest{Amt: s.amt}
	return s.c.do(ctx, http.MethodPost, "/api/v5/finance/staking-defi/sol/redeem", nil, req, true, nil)
}
