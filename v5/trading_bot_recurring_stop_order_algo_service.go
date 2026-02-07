package okx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const tradingBotRecurringMaxStopOrders = 10

// TradingBotRecurringStopOrder 表示定投策略停止请求项。
type TradingBotRecurringStopOrder struct {
	AlgoId string `json:"algoId"`
}

// TradingBotRecurringStopOrderAlgoService 定投策略停止。
type TradingBotRecurringStopOrderAlgoService struct {
	c *Client

	orders []TradingBotRecurringStopOrder
}

// NewTradingBotRecurringStopOrderAlgoService 创建 TradingBotRecurringStopOrderAlgoService。
func (c *Client) NewTradingBotRecurringStopOrderAlgoService() *TradingBotRecurringStopOrderAlgoService {
	return &TradingBotRecurringStopOrderAlgoService{c: c}
}

// Orders 设置批量停止列表（最多 10 个）。
func (s *TradingBotRecurringStopOrderAlgoService) Orders(orders []TradingBotRecurringStopOrder) *TradingBotRecurringStopOrderAlgoService {
	s.orders = orders
	return s
}

var (
	errTradingBotRecurringStopOrderAlgoMissingOrders = errors.New("okx: tradingBot recurring stop-order-algo requires at least one order")
	errTradingBotRecurringStopOrderAlgoTooManyOrders = errors.New("okx: tradingBot recurring stop-order-algo max 10 orders")
)

// Do 定投策略停止（POST /api/v5/tradingBot/recurring/stop-order-algo）。
func (s *TradingBotRecurringStopOrderAlgoService) Do(ctx context.Context) ([]TradingBotOrderAck, error) {
	if len(s.orders) == 0 {
		return nil, errTradingBotRecurringStopOrderAlgoMissingOrders
	}
	if len(s.orders) > tradingBotRecurringMaxStopOrders {
		return nil, errTradingBotRecurringStopOrderAlgoTooManyOrders
	}

	req := make([]TradingBotRecurringStopOrder, 0, len(s.orders))
	for i, o := range s.orders {
		if o.AlgoId == "" {
			return nil, fmt.Errorf("okx: tradingBot recurring stop-order-algo orders[%d] missing algoId", i)
		}
		req = append(req, o)
	}

	var data []TradingBotOrderAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/tradingBot/recurring/stop-order-algo", nil, req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if err := tradingBotCheckBatchAcks(http.MethodPost, "/api/v5/tradingBot/recurring/stop-order-algo", requestID, data); err != nil {
		return data, err
	}
	return data, nil
}
