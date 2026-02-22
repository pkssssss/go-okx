package okx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type tradingBotGridCancelCloseOrderRequest struct {
	AlgoId string `json:"algoId"`
	OrdId  string `json:"ordId"`
}

// TradingBotGridCancelCloseOrderService 撤销合约网格平仓单。
type TradingBotGridCancelCloseOrderService struct {
	c *Client
	r tradingBotGridCancelCloseOrderRequest
}

// NewTradingBotGridCancelCloseOrderService 创建 TradingBotGridCancelCloseOrderService。
func (c *Client) NewTradingBotGridCancelCloseOrderService() *TradingBotGridCancelCloseOrderService {
	return &TradingBotGridCancelCloseOrderService{c: c}
}

func (s *TradingBotGridCancelCloseOrderService) AlgoId(algoId string) *TradingBotGridCancelCloseOrderService {
	s.r.AlgoId = algoId
	return s
}

func (s *TradingBotGridCancelCloseOrderService) OrdId(ordId string) *TradingBotGridCancelCloseOrderService {
	s.r.OrdId = ordId
	return s
}

var (
	errTradingBotGridCancelCloseOrderMissingRequired = errors.New("okx: tradingBot grid cancel-close-order requires algoId and ordId")
	errEmptyTradingBotGridCancelCloseOrderResponse   = errors.New("okx: empty tradingBot grid cancel-close-order response")
	errInvalidTradingBotGridCancelCloseOrderResponse = errors.New("okx: invalid tradingBot grid cancel-close-order response")
)

// Do 撤销合约网格平仓单（POST /api/v5/tradingBot/grid/cancel-close-order）。
func (s *TradingBotGridCancelCloseOrderService) Do(ctx context.Context) (*TradingBotGridCloseOrderAck, error) {
	if s.r.AlgoId == "" || s.r.OrdId == "" {
		return nil, errTradingBotGridCancelCloseOrderMissingRequired
	}

	var data []TradingBotGridCloseOrderAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/tradingBot/grid/cancel-close-order", nil, s.r, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/tradingBot/grid/cancel-close-order", requestID, errEmptyTradingBotGridCancelCloseOrderResponse)
	}
	if len(data) != 1 {
		return nil, newInvalidDataAPIError(
			http.MethodPost,
			"/api/v5/tradingBot/grid/cancel-close-order",
			requestID,
			fmt.Errorf("%w: expected 1 ack, got %d", errInvalidTradingBotGridCancelCloseOrderResponse, len(data)),
		)
	}
	if !validTradingBotGridCloseOrderAck(&data[0]) {
		return nil, newInvalidDataAPIError(http.MethodPost, "/api/v5/tradingBot/grid/cancel-close-order", requestID, errInvalidTradingBotGridCancelCloseOrderResponse)
	}
	return &data[0], nil
}
