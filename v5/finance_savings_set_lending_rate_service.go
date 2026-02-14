package okx

import (
	"context"
	"errors"
	"net/http"
)

type financeSavingsSetLendingRateRequest struct {
	Ccy  string `json:"ccy"`
	Rate string `json:"rate"`
}

// FinanceSavingsSetLendingRateService 设置活期简单赚币借贷利率。
type FinanceSavingsSetLendingRateService struct {
	c   *Client
	req financeSavingsSetLendingRateRequest
}

// NewFinanceSavingsSetLendingRateService 创建 FinanceSavingsSetLendingRateService。
func (c *Client) NewFinanceSavingsSetLendingRateService() *FinanceSavingsSetLendingRateService {
	return &FinanceSavingsSetLendingRateService{c: c}
}

// Ccy 设置币种名称（必填）。
func (s *FinanceSavingsSetLendingRateService) Ccy(ccy string) *FinanceSavingsSetLendingRateService {
	s.req.Ccy = ccy
	return s
}

// Rate 设置贷出年利率（必填，字符串）。
func (s *FinanceSavingsSetLendingRateService) Rate(rate string) *FinanceSavingsSetLendingRateService {
	s.req.Rate = rate
	return s
}

var (
	errFinanceSavingsSetLendingRateMissingRequired = errors.New("okx: savings set lending rate requires ccy and rate")
	errEmptyFinanceSavingsSetLendingRate           = errors.New("okx: empty savings set lending rate response")
)

// Do 设置活期简单赚币借贷利率（POST /api/v5/finance/savings/set-lending-rate）。
func (s *FinanceSavingsSetLendingRateService) Do(ctx context.Context) (*FinanceSavingsSetLendingRateAck, error) {
	if s.req.Ccy == "" || s.req.Rate == "" {
		return nil, errFinanceSavingsSetLendingRateMissingRequired
	}

	var data []FinanceSavingsSetLendingRateAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/finance/savings/set-lending-rate", nil, s.req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/finance/savings/set-lending-rate", requestID, errEmptyFinanceSavingsSetLendingRate)
	}
	return &data[0], nil
}
