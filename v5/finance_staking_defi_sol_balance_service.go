package okx

import (
	"context"
	"errors"
	"net/http"
)

// FinanceStakingDefiSOLBalanceService 获取 SOL 质押余额（OKSOL 实时数据）。
type FinanceStakingDefiSOLBalanceService struct {
	c *Client
}

// NewFinanceStakingDefiSOLBalanceService 创建 FinanceStakingDefiSOLBalanceService。
func (c *Client) NewFinanceStakingDefiSOLBalanceService() *FinanceStakingDefiSOLBalanceService {
	return &FinanceStakingDefiSOLBalanceService{c: c}
}

var errEmptyFinanceStakingDefiSOLBalance = errors.New("okx: empty staking-defi sol balance response")

// Do 获取 SOL 质押余额（GET /api/v5/finance/staking-defi/sol/balance）。
func (s *FinanceStakingDefiSOLBalanceService) Do(ctx context.Context) (*FinanceStakingDefiSOLBalance, error) {
	var data []FinanceStakingDefiSOLBalance
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/finance/staking-defi/sol/balance", nil, nil, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyFinanceStakingDefiSOLBalance
	}
	return &data[0], nil
}
