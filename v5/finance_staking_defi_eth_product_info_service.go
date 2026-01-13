package okx

import (
	"context"
	"errors"
	"net/http"
)

// FinanceStakingDefiETHProductInfoService 获取 ETH 质押产品信息。
type FinanceStakingDefiETHProductInfoService struct {
	c *Client
}

// NewFinanceStakingDefiETHProductInfoService 创建 FinanceStakingDefiETHProductInfoService。
func (c *Client) NewFinanceStakingDefiETHProductInfoService() *FinanceStakingDefiETHProductInfoService {
	return &FinanceStakingDefiETHProductInfoService{c: c}
}

var errEmptyFinanceStakingDefiETHProductInfo = errors.New("okx: empty staking-defi eth product-info response")

// Do 获取 ETH 质押产品信息（GET /api/v5/finance/staking-defi/eth/product-info）。
func (s *FinanceStakingDefiETHProductInfoService) Do(ctx context.Context) (*FinanceStakingDefiETHProductInfo, error) {
	var data []FinanceStakingDefiETHProductInfo
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/finance/staking-defi/eth/product-info", nil, nil, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyFinanceStakingDefiETHProductInfo
	}
	return &data[0], nil
}
