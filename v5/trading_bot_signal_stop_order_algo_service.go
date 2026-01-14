package okx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const tradingBotSignalMaxStopOrders = 10

// TradingBotSignalStopOrder 表示信号策略停止请求项。
type TradingBotSignalStopOrder struct {
	AlgoId string `json:"algoId"`
}

// TradingBotSignalStopOrderAlgoService 停止信号策略。
type TradingBotSignalStopOrderAlgoService struct {
	c *Client

	orders []TradingBotSignalStopOrder
}

// NewTradingBotSignalStopOrderAlgoService 创建 TradingBotSignalStopOrderAlgoService。
func (c *Client) NewTradingBotSignalStopOrderAlgoService() *TradingBotSignalStopOrderAlgoService {
	return &TradingBotSignalStopOrderAlgoService{c: c}
}

// Orders 设置批量停止列表（最多 10 个）。
func (s *TradingBotSignalStopOrderAlgoService) Orders(orders []TradingBotSignalStopOrder) *TradingBotSignalStopOrderAlgoService {
	s.orders = orders
	return s
}

var (
	errTradingBotSignalStopOrderAlgoMissingOrders = errors.New("okx: tradingBot signal stop-order-algo requires at least one order")
	errTradingBotSignalStopOrderAlgoTooManyOrders = errors.New("okx: tradingBot signal stop-order-algo max 10 orders")
)

// Do 停止信号策略（POST /api/v5/tradingBot/signal/stop-order-algo）。
func (s *TradingBotSignalStopOrderAlgoService) Do(ctx context.Context) ([]TradingBotOrderAck, error) {
	if len(s.orders) == 0 {
		return nil, errTradingBotSignalStopOrderAlgoMissingOrders
	}
	if len(s.orders) > tradingBotSignalMaxStopOrders {
		return nil, errTradingBotSignalStopOrderAlgoTooManyOrders
	}

	req := make([]TradingBotSignalStopOrder, 0, len(s.orders))
	for i, o := range s.orders {
		if o.AlgoId == "" {
			return nil, fmt.Errorf("okx: tradingBot signal stop-order-algo orders[%d] missing algoId", i)
		}
		req = append(req, o)
	}

	var data []TradingBotOrderAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/tradingBot/signal/stop-order-algo", nil, req, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
