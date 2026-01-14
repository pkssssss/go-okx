package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// TradingBotGridOrdersAlgoDetailsService 获取网格策略委托订单详情。
type TradingBotGridOrdersAlgoDetailsService struct {
	c *Client

	algoOrdType string
	algoId      string
}

// NewTradingBotGridOrdersAlgoDetailsService 创建 TradingBotGridOrdersAlgoDetailsService。
func (c *Client) NewTradingBotGridOrdersAlgoDetailsService() *TradingBotGridOrdersAlgoDetailsService {
	return &TradingBotGridOrdersAlgoDetailsService{c: c}
}

func (s *TradingBotGridOrdersAlgoDetailsService) AlgoOrdType(algoOrdType string) *TradingBotGridOrdersAlgoDetailsService {
	s.algoOrdType = algoOrdType
	return s
}

func (s *TradingBotGridOrdersAlgoDetailsService) AlgoId(algoId string) *TradingBotGridOrdersAlgoDetailsService {
	s.algoId = algoId
	return s
}

var (
	errTradingBotGridOrdersAlgoDetailsMissingRequired = errors.New("okx: tradingBot grid orders-algo-details requires algoOrdType and algoId")
	errEmptyTradingBotGridOrdersAlgoDetailsResponse   = errors.New("okx: empty tradingBot grid orders-algo-details response")
)

// Do 获取网格策略委托订单详情（GET /api/v5/tradingBot/grid/orders-algo-details）。
func (s *TradingBotGridOrdersAlgoDetailsService) Do(ctx context.Context) (*TradingBotGridOrder, error) {
	if s.algoOrdType == "" || s.algoId == "" {
		return nil, errTradingBotGridOrdersAlgoDetailsMissingRequired
	}

	q := url.Values{}
	q.Set("algoOrdType", s.algoOrdType)
	q.Set("algoId", s.algoId)

	var data []TradingBotGridOrder
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/tradingBot/grid/orders-algo-details", q, nil, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyTradingBotGridOrdersAlgoDetailsResponse
	}
	return &data[0], nil
}
