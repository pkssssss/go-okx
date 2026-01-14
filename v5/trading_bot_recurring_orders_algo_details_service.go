package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// TradingBotRecurringOrdersAlgoDetailsService 获取定投策略委托订单详情。
type TradingBotRecurringOrdersAlgoDetailsService struct {
	c *Client

	algoId string
}

// NewTradingBotRecurringOrdersAlgoDetailsService 创建 TradingBotRecurringOrdersAlgoDetailsService。
func (c *Client) NewTradingBotRecurringOrdersAlgoDetailsService() *TradingBotRecurringOrdersAlgoDetailsService {
	return &TradingBotRecurringOrdersAlgoDetailsService{c: c}
}

func (s *TradingBotRecurringOrdersAlgoDetailsService) AlgoId(algoId string) *TradingBotRecurringOrdersAlgoDetailsService {
	s.algoId = algoId
	return s
}

var (
	errTradingBotRecurringOrdersAlgoDetailsMissingAlgoId = errors.New("okx: tradingBot recurring orders-algo-details requires algoId")
	errEmptyTradingBotRecurringOrdersAlgoDetailsResponse = errors.New("okx: empty tradingBot recurring orders-algo-details response")
)

// Do 获取定投策略委托订单详情（GET /api/v5/tradingBot/recurring/orders-algo-details）。
func (s *TradingBotRecurringOrdersAlgoDetailsService) Do(ctx context.Context) (*TradingBotRecurringOrder, error) {
	if s.algoId == "" {
		return nil, errTradingBotRecurringOrdersAlgoDetailsMissingAlgoId
	}

	q := url.Values{}
	q.Set("algoId", s.algoId)

	var data []TradingBotRecurringOrder
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/tradingBot/recurring/orders-algo-details", q, nil, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyTradingBotRecurringOrdersAlgoDetailsResponse
	}
	return &data[0], nil
}
