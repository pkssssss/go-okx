package okx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type tradingBotGridAdjustInvestmentRequest struct {
	AlgoId              string `json:"algoId"`
	Amt                 string `json:"amt"`
	AllowReinvestProfit string `json:"allowReinvestProfit,omitempty"`
}

// TradingBotGridAdjustInvestmentService 网格加仓。
type TradingBotGridAdjustInvestmentService struct {
	c *Client
	r tradingBotGridAdjustInvestmentRequest
}

// NewTradingBotGridAdjustInvestmentService 创建 TradingBotGridAdjustInvestmentService。
func (c *Client) NewTradingBotGridAdjustInvestmentService() *TradingBotGridAdjustInvestmentService {
	return &TradingBotGridAdjustInvestmentService{c: c}
}

func (s *TradingBotGridAdjustInvestmentService) AlgoId(algoId string) *TradingBotGridAdjustInvestmentService {
	s.r.AlgoId = algoId
	return s
}

func (s *TradingBotGridAdjustInvestmentService) Amt(amt string) *TradingBotGridAdjustInvestmentService {
	s.r.Amt = amt
	return s
}

func (s *TradingBotGridAdjustInvestmentService) AllowReinvestProfit(allowReinvestProfit string) *TradingBotGridAdjustInvestmentService {
	s.r.AllowReinvestProfit = allowReinvestProfit
	return s
}

var (
	errTradingBotGridAdjustInvestmentMissingRequired = errors.New("okx: tradingBot grid adjust-investment requires algoId and amt")
	errEmptyTradingBotGridAdjustInvestmentResponse   = errors.New("okx: empty tradingBot grid adjust-investment response")
	errInvalidTradingBotGridAdjustInvestmentResponse = errors.New("okx: invalid tradingBot grid adjust-investment response")
)

// Do 网格加仓（POST /api/v5/tradingBot/grid/adjust-investment）。
func (s *TradingBotGridAdjustInvestmentService) Do(ctx context.Context) (*TradingBotAlgoIdAck, error) {
	if s.r.AlgoId == "" || s.r.Amt == "" {
		return nil, errTradingBotGridAdjustInvestmentMissingRequired
	}

	var data []TradingBotAlgoIdAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/tradingBot/grid/adjust-investment", nil, s.r, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/tradingBot/grid/adjust-investment", requestID, errEmptyTradingBotGridAdjustInvestmentResponse)
	}
	if len(data) != 1 {
		return nil, newInvalidDataAPIError(
			http.MethodPost,
			"/api/v5/tradingBot/grid/adjust-investment",
			requestID,
			fmt.Errorf("%w: expected 1 ack, got %d", errInvalidTradingBotGridAdjustInvestmentResponse, len(data)),
		)
	}
	if !validTradingBotAlgoIDAck(&data[0]) {
		return nil, newInvalidDataAPIError(http.MethodPost, "/api/v5/tradingBot/grid/adjust-investment", requestID, errInvalidTradingBotGridAdjustInvestmentResponse)
	}
	return &data[0], nil
}
