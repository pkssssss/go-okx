package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// FinanceStakingDefiSOLAPYHistoryService 获取 SOL 质押历史收益率（公共）。
type FinanceStakingDefiSOLAPYHistoryService struct {
	c    *Client
	days string
}

// NewFinanceStakingDefiSOLAPYHistoryService 创建 FinanceStakingDefiSOLAPYHistoryService。
func (c *Client) NewFinanceStakingDefiSOLAPYHistoryService() *FinanceStakingDefiSOLAPYHistoryService {
	return &FinanceStakingDefiSOLAPYHistoryService{c: c}
}

// Days 设置查询最近多少天内的数据（必填，不超过 365 天）。
func (s *FinanceStakingDefiSOLAPYHistoryService) Days(days string) *FinanceStakingDefiSOLAPYHistoryService {
	s.days = days
	return s
}

var errFinanceStakingDefiSOLAPYHistoryMissingDays = errors.New("okx: staking-defi sol apy-history requires days")

// Do 获取 SOL 质押历史收益率（GET /api/v5/finance/staking-defi/sol/apy-history）。
func (s *FinanceStakingDefiSOLAPYHistoryService) Do(ctx context.Context) ([]FinanceStakingDefiAPYHistory, error) {
	if s.days == "" {
		return nil, errFinanceStakingDefiSOLAPYHistoryMissingDays
	}

	q := url.Values{}
	q.Set("days", s.days)

	var data []FinanceStakingDefiAPYHistory
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/finance/staking-defi/sol/apy-history", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
