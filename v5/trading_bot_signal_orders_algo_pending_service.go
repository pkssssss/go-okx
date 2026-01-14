package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// TradingBotSignalOrdersAlgoPendingService 获取活跃信号策略列表。
type TradingBotSignalOrdersAlgoPendingService struct {
	c *Client

	algoOrdType string
	algoId      string
	after       string
	before      string
	limit       *int
}

// NewTradingBotSignalOrdersAlgoPendingService 创建 TradingBotSignalOrdersAlgoPendingService。
func (c *Client) NewTradingBotSignalOrdersAlgoPendingService() *TradingBotSignalOrdersAlgoPendingService {
	return &TradingBotSignalOrdersAlgoPendingService{c: c}
}

func (s *TradingBotSignalOrdersAlgoPendingService) AlgoOrdType(algoOrdType string) *TradingBotSignalOrdersAlgoPendingService {
	s.algoOrdType = algoOrdType
	return s
}

func (s *TradingBotSignalOrdersAlgoPendingService) AlgoId(algoId string) *TradingBotSignalOrdersAlgoPendingService {
	s.algoId = algoId
	return s
}

func (s *TradingBotSignalOrdersAlgoPendingService) After(after string) *TradingBotSignalOrdersAlgoPendingService {
	s.after = after
	return s
}

func (s *TradingBotSignalOrdersAlgoPendingService) Before(before string) *TradingBotSignalOrdersAlgoPendingService {
	s.before = before
	return s
}

func (s *TradingBotSignalOrdersAlgoPendingService) Limit(limit int) *TradingBotSignalOrdersAlgoPendingService {
	s.limit = &limit
	return s
}

var errTradingBotSignalOrdersAlgoPendingMissingAlgoOrdType = errors.New("okx: tradingBot signal orders-algo-pending requires algoOrdType")

// Do 获取活跃信号策略列表（GET /api/v5/tradingBot/signal/orders-algo-pending）。
func (s *TradingBotSignalOrdersAlgoPendingService) Do(ctx context.Context) ([]TradingBotSignalOrder, error) {
	if s.algoOrdType == "" {
		return nil, errTradingBotSignalOrdersAlgoPendingMissingAlgoOrdType
	}

	q := url.Values{}
	q.Set("algoOrdType", s.algoOrdType)
	if s.algoId != "" {
		q.Set("algoId", s.algoId)
	}
	if s.after != "" {
		q.Set("after", s.after)
	}
	if s.before != "" {
		q.Set("before", s.before)
	}
	if s.limit != nil {
		q.Set("limit", strconv.Itoa(*s.limit))
	}

	var data []TradingBotSignalOrder
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/tradingBot/signal/orders-algo-pending", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
