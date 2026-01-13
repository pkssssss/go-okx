package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// FinanceFlexibleLoanMaxCollateralRedeemAmountService 获取抵押物最大可赎回数量。
type FinanceFlexibleLoanMaxCollateralRedeemAmountService struct {
	c   *Client
	ccy string
}

// NewFinanceFlexibleLoanMaxCollateralRedeemAmountService 创建 FinanceFlexibleLoanMaxCollateralRedeemAmountService。
func (c *Client) NewFinanceFlexibleLoanMaxCollateralRedeemAmountService() *FinanceFlexibleLoanMaxCollateralRedeemAmountService {
	return &FinanceFlexibleLoanMaxCollateralRedeemAmountService{c: c}
}

// Ccy 设置抵押物币种（必填）。
func (s *FinanceFlexibleLoanMaxCollateralRedeemAmountService) Ccy(ccy string) *FinanceFlexibleLoanMaxCollateralRedeemAmountService {
	s.ccy = ccy
	return s
}

var (
	errFinanceFlexibleLoanMaxCollateralRedeemAmountMissingCcy = errors.New("okx: flexible loan max collateral redeem amount requires ccy")
	errEmptyFinanceFlexibleLoanMaxCollateralRedeemAmount      = errors.New("okx: empty flexible loan max collateral redeem amount response")
)

// Do 获取抵押物最大可赎回数量（GET /api/v5/finance/flexible-loan/max-collateral-redeem-amount）。
func (s *FinanceFlexibleLoanMaxCollateralRedeemAmountService) Do(ctx context.Context) (*FinanceFlexibleLoanMaxCollateralRedeemAmount, error) {
	if s.ccy == "" {
		return nil, errFinanceFlexibleLoanMaxCollateralRedeemAmountMissingCcy
	}

	q := url.Values{}
	q.Set("ccy", s.ccy)

	var data []FinanceFlexibleLoanMaxCollateralRedeemAmount
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/finance/flexible-loan/max-collateral-redeem-amount", q, nil, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyFinanceFlexibleLoanMaxCollateralRedeemAmount
	}
	return &data[0], nil
}
