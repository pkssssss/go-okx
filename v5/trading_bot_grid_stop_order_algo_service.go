package okx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const tradingBotGridMaxStopOrders = 10

// TradingBotGridStopOrder 表示网格策略停止请求项。
type TradingBotGridStopOrder struct {
	AlgoId      string `json:"algoId"`
	InstId      string `json:"instId"`
	AlgoOrdType string `json:"algoOrdType"`
	StopType    string `json:"stopType"`
}

// TradingBotGridStopOrderAlgoService 网格策略停止。
type TradingBotGridStopOrderAlgoService struct {
	c *Client

	orders []TradingBotGridStopOrder
}

// NewTradingBotGridStopOrderAlgoService 创建 TradingBotGridStopOrderAlgoService。
func (c *Client) NewTradingBotGridStopOrderAlgoService() *TradingBotGridStopOrderAlgoService {
	return &TradingBotGridStopOrderAlgoService{c: c}
}

// Orders 设置批量停止列表（最多 10 个）。
func (s *TradingBotGridStopOrderAlgoService) Orders(orders []TradingBotGridStopOrder) *TradingBotGridStopOrderAlgoService {
	s.orders = orders
	return s
}

var (
	errTradingBotGridStopOrderAlgoMissingOrders = errors.New("okx: tradingBot grid stop-order-algo requires at least one order")
	errTradingBotGridStopOrderAlgoTooManyOrders = errors.New("okx: tradingBot grid stop-order-algo max 10 orders")
)

// Do 网格策略停止（POST /api/v5/tradingBot/grid/stop-order-algo）。
func (s *TradingBotGridStopOrderAlgoService) Do(ctx context.Context) ([]TradingBotOrderAck, error) {
	if len(s.orders) == 0 {
		return nil, errTradingBotGridStopOrderAlgoMissingOrders
	}
	if len(s.orders) > tradingBotGridMaxStopOrders {
		return nil, errTradingBotGridStopOrderAlgoTooManyOrders
	}

	req := make([]TradingBotGridStopOrder, 0, len(s.orders))
	for i, o := range s.orders {
		if o.AlgoId == "" || o.InstId == "" || o.AlgoOrdType == "" || o.StopType == "" {
			return nil, fmt.Errorf("okx: tradingBot grid stop-order-algo orders[%d] missing required fields", i)
		}
		req = append(req, o)
	}

	var data []TradingBotOrderAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/tradingBot/grid/stop-order-algo", nil, req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if err := tradingBotCheckBatchAcks(http.MethodPost, "/api/v5/tradingBot/grid/stop-order-algo", requestID, len(req), data); err != nil {
		return data, err
	}
	return data, nil
}
