package okx

import (
	"context"
	"errors"
	"net/http"
)

type financeStakingDefiSOLPurchaseRequest struct {
	Amt string `json:"amt"`
}

// FinanceStakingDefiSOLPurchaseService 申购 SOL 质押（质押 SOL 获取 OKSOL）。
type FinanceStakingDefiSOLPurchaseService struct {
	c   *Client
	amt string
}

// NewFinanceStakingDefiSOLPurchaseService 创建 FinanceStakingDefiSOLPurchaseService。
func (c *Client) NewFinanceStakingDefiSOLPurchaseService() *FinanceStakingDefiSOLPurchaseService {
	return &FinanceStakingDefiSOLPurchaseService{c: c}
}

// Amt 设置投资数量（必填，字符串）。
func (s *FinanceStakingDefiSOLPurchaseService) Amt(amt string) *FinanceStakingDefiSOLPurchaseService {
	s.amt = amt
	return s
}

var errFinanceStakingDefiSOLPurchaseMissingAmt = errors.New("okx: staking-defi sol purchase requires amt")

// Do 申购 SOL 质押（POST /api/v5/finance/staking-defi/sol/purchase）。
//
// 注意：OKX 返回的 data 为空数组（code=0 代表请求已被成功处理）。
func (s *FinanceStakingDefiSOLPurchaseService) Do(ctx context.Context) error {
	if s.amt == "" {
		return errFinanceStakingDefiSOLPurchaseMissingAmt
	}

	req := financeStakingDefiSOLPurchaseRequest{Amt: s.amt}
	return s.c.do(ctx, http.MethodPost, "/api/v5/finance/staking-defi/sol/purchase", nil, req, true, nil)
}
