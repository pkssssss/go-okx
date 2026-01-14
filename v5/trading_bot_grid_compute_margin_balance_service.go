package okx

import (
	"context"
	"errors"
	"net/http"
)

type tradingBotGridComputeMarginBalanceRequest struct {
	AlgoId string `json:"algoId"`
	Type   string `json:"type"`
	Amt    string `json:"amt,omitempty"`
}

// TradingBotGridComputeMarginBalanceService 调整保证金计算。
type TradingBotGridComputeMarginBalanceService struct {
	c *Client
	r tradingBotGridComputeMarginBalanceRequest
}

// NewTradingBotGridComputeMarginBalanceService 创建 TradingBotGridComputeMarginBalanceService。
func (c *Client) NewTradingBotGridComputeMarginBalanceService() *TradingBotGridComputeMarginBalanceService {
	return &TradingBotGridComputeMarginBalanceService{c: c}
}

func (s *TradingBotGridComputeMarginBalanceService) AlgoId(algoId string) *TradingBotGridComputeMarginBalanceService {
	s.r.AlgoId = algoId
	return s
}

func (s *TradingBotGridComputeMarginBalanceService) Type(typ string) *TradingBotGridComputeMarginBalanceService {
	s.r.Type = typ
	return s
}

func (s *TradingBotGridComputeMarginBalanceService) Amt(amt string) *TradingBotGridComputeMarginBalanceService {
	s.r.Amt = amt
	return s
}

var (
	errTradingBotGridComputeMarginBalanceMissingRequired = errors.New("okx: tradingBot grid compute-margin-balance requires algoId and type")
	errEmptyTradingBotGridComputeMarginBalanceResponse   = errors.New("okx: empty tradingBot grid compute-margin-balance response")
)

// Do 调整保证金计算（POST /api/v5/tradingBot/grid/compute-margin-balance）。
func (s *TradingBotGridComputeMarginBalanceService) Do(ctx context.Context) (*TradingBotGridComputeMarginBalanceResult, error) {
	if s.r.AlgoId == "" || s.r.Type == "" {
		return nil, errTradingBotGridComputeMarginBalanceMissingRequired
	}

	var data []TradingBotGridComputeMarginBalanceResult
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/tradingBot/grid/compute-margin-balance", nil, s.r, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyTradingBotGridComputeMarginBalanceResponse
	}
	return &data[0], nil
}
