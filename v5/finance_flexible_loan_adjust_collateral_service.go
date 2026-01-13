package okx

import (
	"context"
	"errors"
	"net/http"
)

type financeFlexibleLoanAdjustCollateralRequest struct {
	Type          string `json:"type"`
	CollateralCcy string `json:"collateralCcy"`
	CollateralAmt string `json:"collateralAmt"`
}

// FinanceFlexibleLoanAdjustCollateralService 调整抵押物。
type FinanceFlexibleLoanAdjustCollateralService struct {
	c   *Client
	req financeFlexibleLoanAdjustCollateralRequest
}

// NewFinanceFlexibleLoanAdjustCollateralService 创建 FinanceFlexibleLoanAdjustCollateralService。
func (c *Client) NewFinanceFlexibleLoanAdjustCollateralService() *FinanceFlexibleLoanAdjustCollateralService {
	return &FinanceFlexibleLoanAdjustCollateralService{c: c}
}

// Type 设置操作类型（必填：add/reduce）。
func (s *FinanceFlexibleLoanAdjustCollateralService) Type(typ string) *FinanceFlexibleLoanAdjustCollateralService {
	s.req.Type = typ
	return s
}

// CollateralCcy 设置抵押物币种（必填）。
func (s *FinanceFlexibleLoanAdjustCollateralService) CollateralCcy(collateralCcy string) *FinanceFlexibleLoanAdjustCollateralService {
	s.req.CollateralCcy = collateralCcy
	return s
}

// CollateralAmt 设置抵押物数量（必填，字符串）。
func (s *FinanceFlexibleLoanAdjustCollateralService) CollateralAmt(collateralAmt string) *FinanceFlexibleLoanAdjustCollateralService {
	s.req.CollateralAmt = collateralAmt
	return s
}

var errFinanceFlexibleLoanAdjustCollateralMissingRequired = errors.New("okx: flexible loan adjust collateral requires type, collateralCcy and collateralAmt")

// Do 调整抵押物（POST /api/v5/finance/flexible-loan/adjust-collateral）。
//
// 注意：OKX 返回的 data 为空数组（code=0 代表请求已被接受，不代表处理完成）。
func (s *FinanceFlexibleLoanAdjustCollateralService) Do(ctx context.Context) error {
	if s.req.Type == "" || s.req.CollateralCcy == "" || s.req.CollateralAmt == "" {
		return errFinanceFlexibleLoanAdjustCollateralMissingRequired
	}
	return s.c.do(ctx, http.MethodPost, "/api/v5/finance/flexible-loan/adjust-collateral", nil, s.req, true, nil)
}
