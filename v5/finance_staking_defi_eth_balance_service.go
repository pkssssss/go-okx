package okx

import (
	"context"
	"errors"
	"net/http"
)

// FinanceStakingDefiETHBalanceService 获取 ETH 质押余额（BETH 快照）。
type FinanceStakingDefiETHBalanceService struct {
	c *Client
}

// NewFinanceStakingDefiETHBalanceService 创建 FinanceStakingDefiETHBalanceService。
func (c *Client) NewFinanceStakingDefiETHBalanceService() *FinanceStakingDefiETHBalanceService {
	return &FinanceStakingDefiETHBalanceService{c: c}
}

var errEmptyFinanceStakingDefiETHBalance = errors.New("okx: empty staking-defi eth balance response")

// Do 获取 ETH 质押余额（GET /api/v5/finance/staking-defi/eth/balance）。
func (s *FinanceStakingDefiETHBalanceService) Do(ctx context.Context) (*FinanceStakingDefiETHBalance, error) {
	var data []FinanceStakingDefiETHBalance
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/finance/staking-defi/eth/balance", nil, nil, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyFinanceStakingDefiETHBalance
	}
	return &data[0], nil
}
