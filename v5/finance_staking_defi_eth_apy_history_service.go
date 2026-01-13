package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// FinanceStakingDefiETHAPYHistoryService 获取 ETH 质押历史收益率（公共）。
type FinanceStakingDefiETHAPYHistoryService struct {
	c    *Client
	days string
}

// NewFinanceStakingDefiETHAPYHistoryService 创建 FinanceStakingDefiETHAPYHistoryService。
func (c *Client) NewFinanceStakingDefiETHAPYHistoryService() *FinanceStakingDefiETHAPYHistoryService {
	return &FinanceStakingDefiETHAPYHistoryService{c: c}
}

// Days 设置查询最近多少天内的数据（必填，不超过 365 天）。
func (s *FinanceStakingDefiETHAPYHistoryService) Days(days string) *FinanceStakingDefiETHAPYHistoryService {
	s.days = days
	return s
}

var errFinanceStakingDefiETHAPYHistoryMissingDays = errors.New("okx: staking-defi eth apy-history requires days")

// Do 获取 ETH 质押历史收益率（GET /api/v5/finance/staking-defi/eth/apy-history）。
func (s *FinanceStakingDefiETHAPYHistoryService) Do(ctx context.Context) ([]FinanceStakingDefiAPYHistory, error) {
	if s.days == "" {
		return nil, errFinanceStakingDefiETHAPYHistoryMissingDays
	}

	q := url.Values{}
	q.Set("days", s.days)

	var data []FinanceStakingDefiAPYHistory
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/finance/staking-defi/eth/apy-history", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
