package okx

import (
	"context"
	"errors"
	"net/http"
)

type financeSavingsPurchaseRedemptRequest struct {
	Ccy  string `json:"ccy"`
	Amt  string `json:"amt"`
	Side string `json:"side"`
	Rate string `json:"rate,omitempty"`
}

// FinanceSavingsPurchaseRedemptService 活期简单赚币申购/赎回。
type FinanceSavingsPurchaseRedemptService struct {
	c   *Client
	req financeSavingsPurchaseRedemptRequest
}

// NewFinanceSavingsPurchaseRedemptService 创建 FinanceSavingsPurchaseRedemptService。
func (c *Client) NewFinanceSavingsPurchaseRedemptService() *FinanceSavingsPurchaseRedemptService {
	return &FinanceSavingsPurchaseRedemptService{c: c}
}

// Ccy 设置币种名称（必填）。
func (s *FinanceSavingsPurchaseRedemptService) Ccy(ccy string) *FinanceSavingsPurchaseRedemptService {
	s.req.Ccy = ccy
	return s
}

// Amt 设置申购/赎回数量（必填，字符串）。
func (s *FinanceSavingsPurchaseRedemptService) Amt(amt string) *FinanceSavingsPurchaseRedemptService {
	s.req.Amt = amt
	return s
}

// Side 设置操作类型（必填：purchase/redempt）。
func (s *FinanceSavingsPurchaseRedemptService) Side(side string) *FinanceSavingsPurchaseRedemptService {
	s.req.Side = side
	return s
}

// Rate 设置申购年利率（可选，仅适用于申购；字符串）。
func (s *FinanceSavingsPurchaseRedemptService) Rate(rate string) *FinanceSavingsPurchaseRedemptService {
	s.req.Rate = rate
	return s
}

var (
	errFinanceSavingsPurchaseRedemptMissingRequired = errors.New("okx: savings purchase-redempt requires ccy, amt and side")
	errEmptyFinanceSavingsPurchaseRedempt           = errors.New("okx: empty savings purchase-redempt response")
)

// Do 活期简单赚币申购/赎回（POST /api/v5/finance/savings/purchase-redempt）。
func (s *FinanceSavingsPurchaseRedemptService) Do(ctx context.Context) (*FinanceSavingsPurchaseRedemptAck, error) {
	if s.req.Ccy == "" || s.req.Amt == "" || s.req.Side == "" {
		return nil, errFinanceSavingsPurchaseRedemptMissingRequired
	}

	var data []FinanceSavingsPurchaseRedemptAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/finance/savings/purchase-redempt", nil, s.req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/finance/savings/purchase-redempt", requestID, errEmptyFinanceSavingsPurchaseRedempt)
	}
	return &data[0], nil
}
