package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// TradingBotSignalOrdersAlgoDetailsService 获取信号策略详情。
type TradingBotSignalOrdersAlgoDetailsService struct {
	c *Client

	algoOrdType string
	algoId      string
}

// NewTradingBotSignalOrdersAlgoDetailsService 创建 TradingBotSignalOrdersAlgoDetailsService。
func (c *Client) NewTradingBotSignalOrdersAlgoDetailsService() *TradingBotSignalOrdersAlgoDetailsService {
	return &TradingBotSignalOrdersAlgoDetailsService{c: c}
}

func (s *TradingBotSignalOrdersAlgoDetailsService) AlgoOrdType(algoOrdType string) *TradingBotSignalOrdersAlgoDetailsService {
	s.algoOrdType = algoOrdType
	return s
}

func (s *TradingBotSignalOrdersAlgoDetailsService) AlgoId(algoId string) *TradingBotSignalOrdersAlgoDetailsService {
	s.algoId = algoId
	return s
}

var (
	errTradingBotSignalOrdersAlgoDetailsMissingRequired = errors.New("okx: tradingBot signal orders-algo-details requires algoOrdType and algoId")
	errEmptyTradingBotSignalOrdersAlgoDetailsResponse   = errors.New("okx: empty tradingBot signal orders-algo-details response")
)

// Do 获取信号策略详情（GET /api/v5/tradingBot/signal/orders-algo-details）。
func (s *TradingBotSignalOrdersAlgoDetailsService) Do(ctx context.Context) (*TradingBotSignalOrder, error) {
	if s.algoOrdType == "" || s.algoId == "" {
		return nil, errTradingBotSignalOrdersAlgoDetailsMissingRequired
	}

	q := url.Values{}
	q.Set("algoOrdType", s.algoOrdType)
	q.Set("algoId", s.algoId)

	var data []TradingBotSignalOrder
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/tradingBot/signal/orders-algo-details", q, nil, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyTradingBotSignalOrdersAlgoDetailsResponse
	}
	return &data[0], nil
}
