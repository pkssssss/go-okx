package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// FinanceFlexibleLoanCollateralAssetsService 获取可抵押资产信息（仅支持资金账户中的资产）。
type FinanceFlexibleLoanCollateralAssetsService struct {
	c   *Client
	ccy string
}

// NewFinanceFlexibleLoanCollateralAssetsService 创建 FinanceFlexibleLoanCollateralAssetsService。
func (c *Client) NewFinanceFlexibleLoanCollateralAssetsService() *FinanceFlexibleLoanCollateralAssetsService {
	return &FinanceFlexibleLoanCollateralAssetsService{c: c}
}

// Ccy 设置币种过滤（可选）。
func (s *FinanceFlexibleLoanCollateralAssetsService) Ccy(ccy string) *FinanceFlexibleLoanCollateralAssetsService {
	s.ccy = ccy
	return s
}

var errEmptyFinanceFlexibleLoanCollateralAssets = errors.New("okx: empty flexible loan collateral assets response")

// Do 获取可抵押资产信息（GET /api/v5/finance/flexible-loan/collateral-assets）。
func (s *FinanceFlexibleLoanCollateralAssetsService) Do(ctx context.Context) (*FinanceFlexibleLoanCollateralAssets, error) {
	var q url.Values
	if s.ccy != "" {
		q = url.Values{}
		q.Set("ccy", s.ccy)
	}

	var data []FinanceFlexibleLoanCollateralAssets
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/finance/flexible-loan/collateral-assets", q, nil, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyFinanceFlexibleLoanCollateralAssets
	}
	return &data[0], nil
}
