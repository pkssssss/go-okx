package okx

import (
	"context"
	"errors"
	"net/http"
)

type tradingBotGridOrderInstantTriggerRequest struct {
	AlgoId   string `json:"algoId"`
	TopUpAmt string `json:"topUpAmt,omitempty"`
}

// TradingBotGridOrderInstantTriggerService 网格策略立即触发。
type TradingBotGridOrderInstantTriggerService struct {
	c *Client
	r tradingBotGridOrderInstantTriggerRequest
}

// NewTradingBotGridOrderInstantTriggerService 创建 TradingBotGridOrderInstantTriggerService。
func (c *Client) NewTradingBotGridOrderInstantTriggerService() *TradingBotGridOrderInstantTriggerService {
	return &TradingBotGridOrderInstantTriggerService{c: c}
}

func (s *TradingBotGridOrderInstantTriggerService) AlgoId(algoId string) *TradingBotGridOrderInstantTriggerService {
	s.r.AlgoId = algoId
	return s
}

func (s *TradingBotGridOrderInstantTriggerService) TopUpAmt(topUpAmt string) *TradingBotGridOrderInstantTriggerService {
	s.r.TopUpAmt = topUpAmt
	return s
}

var (
	errTradingBotGridOrderInstantTriggerMissingAlgoId = errors.New("okx: tradingBot grid order-instant-trigger requires algoId")
	errEmptyTradingBotGridOrderInstantTriggerResponse = errors.New("okx: empty tradingBot grid order-instant-trigger response")
)

// Do 网格策略立即触发（POST /api/v5/tradingBot/grid/order-instant-trigger）。
func (s *TradingBotGridOrderInstantTriggerService) Do(ctx context.Context) (*TradingBotOrderAck, error) {
	if s.r.AlgoId == "" {
		return nil, errTradingBotGridOrderInstantTriggerMissingAlgoId
	}

	var data []TradingBotOrderAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/tradingBot/grid/order-instant-trigger", nil, s.r, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/tradingBot/grid/order-instant-trigger", requestID, errEmptyTradingBotGridOrderInstantTriggerResponse)
	}
	if data[0].SCode != "0" {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodPost,
			RequestPath: "/api/v5/tradingBot/grid/order-instant-trigger",
			RequestID:   requestID,
			Code:        data[0].SCode,
			Message:     data[0].SMsg,
		}
	}
	return &data[0], nil
}
