package okx

import (
	"context"
	"errors"
	"net/http"
)

type financeStakingDefiETHPurchaseRequest struct {
	Amt string `json:"amt"`
}

// FinanceStakingDefiETHPurchaseService 申购 ETH 质押（质押 ETH 获取 BETH）。
type FinanceStakingDefiETHPurchaseService struct {
	c   *Client
	amt string
}

// NewFinanceStakingDefiETHPurchaseService 创建 FinanceStakingDefiETHPurchaseService。
func (c *Client) NewFinanceStakingDefiETHPurchaseService() *FinanceStakingDefiETHPurchaseService {
	return &FinanceStakingDefiETHPurchaseService{c: c}
}

// Amt 设置投资数量（必填，字符串）。
func (s *FinanceStakingDefiETHPurchaseService) Amt(amt string) *FinanceStakingDefiETHPurchaseService {
	s.amt = amt
	return s
}

var errFinanceStakingDefiETHPurchaseMissingAmt = errors.New("okx: staking-defi eth purchase requires amt")

// Do 申购 ETH 质押（POST /api/v5/finance/staking-defi/eth/purchase）。
//
// 注意：OKX 返回的 data 为空数组（code=0 代表请求已被成功处理）。
func (s *FinanceStakingDefiETHPurchaseService) Do(ctx context.Context) error {
	if s.amt == "" {
		return errFinanceStakingDefiETHPurchaseMissingAmt
	}

	req := financeStakingDefiETHPurchaseRequest{Amt: s.amt}
	return s.c.do(ctx, http.MethodPost, "/api/v5/finance/staking-defi/eth/purchase", nil, req, true, nil)
}
