package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// TradingBotSignalOrdersAlgoHistoryService 获取历史信号策略列表。
type TradingBotSignalOrdersAlgoHistoryService struct {
	c *Client

	algoOrdType string
	algoId      string
	after       string
	before      string
	limit       *int
}

// NewTradingBotSignalOrdersAlgoHistoryService 创建 TradingBotSignalOrdersAlgoHistoryService。
func (c *Client) NewTradingBotSignalOrdersAlgoHistoryService() *TradingBotSignalOrdersAlgoHistoryService {
	return &TradingBotSignalOrdersAlgoHistoryService{c: c}
}

func (s *TradingBotSignalOrdersAlgoHistoryService) AlgoOrdType(algoOrdType string) *TradingBotSignalOrdersAlgoHistoryService {
	s.algoOrdType = algoOrdType
	return s
}

func (s *TradingBotSignalOrdersAlgoHistoryService) AlgoId(algoId string) *TradingBotSignalOrdersAlgoHistoryService {
	s.algoId = algoId
	return s
}

func (s *TradingBotSignalOrdersAlgoHistoryService) After(after string) *TradingBotSignalOrdersAlgoHistoryService {
	s.after = after
	return s
}

func (s *TradingBotSignalOrdersAlgoHistoryService) Before(before string) *TradingBotSignalOrdersAlgoHistoryService {
	s.before = before
	return s
}

func (s *TradingBotSignalOrdersAlgoHistoryService) Limit(limit int) *TradingBotSignalOrdersAlgoHistoryService {
	s.limit = &limit
	return s
}

var errTradingBotSignalOrdersAlgoHistoryMissingRequired = errors.New("okx: tradingBot signal orders-algo-history requires algoOrdType and algoId")

// Do 获取历史信号策略列表（GET /api/v5/tradingBot/signal/orders-algo-history）。
func (s *TradingBotSignalOrdersAlgoHistoryService) Do(ctx context.Context) ([]TradingBotSignalOrder, error) {
	if s.algoOrdType == "" || s.algoId == "" {
		return nil, errTradingBotSignalOrdersAlgoHistoryMissingRequired
	}

	q := url.Values{}
	q.Set("algoOrdType", s.algoOrdType)
	q.Set("algoId", s.algoId)
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
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/tradingBot/signal/orders-algo-history", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
