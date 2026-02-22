package okx

import (
	"context"
	"errors"
	"fmt"
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
	errInvalidTradingBotGridComputeMarginBalanceResponse = errors.New("okx: invalid tradingBot grid compute-margin-balance response")
)

// Do 调整保证金计算（POST /api/v5/tradingBot/grid/compute-margin-balance）。
func (s *TradingBotGridComputeMarginBalanceService) Do(ctx context.Context) (*TradingBotGridComputeMarginBalanceResult, error) {
	if s.r.AlgoId == "" || s.r.Type == "" {
		return nil, errTradingBotGridComputeMarginBalanceMissingRequired
	}

	var data []TradingBotGridComputeMarginBalanceResult
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/tradingBot/grid/compute-margin-balance", nil, s.r, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/tradingBot/grid/compute-margin-balance", requestID, errEmptyTradingBotGridComputeMarginBalanceResponse)
	}
	if len(data) != 1 {
		return nil, newInvalidDataAPIError(
			http.MethodPost,
			"/api/v5/tradingBot/grid/compute-margin-balance",
			requestID,
			fmt.Errorf("%w: expected 1 ack, got %d", errInvalidTradingBotGridComputeMarginBalanceResponse, len(data)),
		)
	}
	return &data[0], nil
}
