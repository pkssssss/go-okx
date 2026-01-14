package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// TradingBotGridOrdersAlgoHistoryService 获取历史网格策略委托单列表。
type TradingBotGridOrdersAlgoHistoryService struct {
	c *Client

	algoOrdType string
	algoId      string
	instId      string
	instType    string
	after       string
	before      string
	limit       *int
}

// NewTradingBotGridOrdersAlgoHistoryService 创建 TradingBotGridOrdersAlgoHistoryService。
func (c *Client) NewTradingBotGridOrdersAlgoHistoryService() *TradingBotGridOrdersAlgoHistoryService {
	return &TradingBotGridOrdersAlgoHistoryService{c: c}
}

func (s *TradingBotGridOrdersAlgoHistoryService) AlgoOrdType(algoOrdType string) *TradingBotGridOrdersAlgoHistoryService {
	s.algoOrdType = algoOrdType
	return s
}

func (s *TradingBotGridOrdersAlgoHistoryService) AlgoId(algoId string) *TradingBotGridOrdersAlgoHistoryService {
	s.algoId = algoId
	return s
}

func (s *TradingBotGridOrdersAlgoHistoryService) InstId(instId string) *TradingBotGridOrdersAlgoHistoryService {
	s.instId = instId
	return s
}

func (s *TradingBotGridOrdersAlgoHistoryService) InstType(instType string) *TradingBotGridOrdersAlgoHistoryService {
	s.instType = instType
	return s
}

func (s *TradingBotGridOrdersAlgoHistoryService) After(after string) *TradingBotGridOrdersAlgoHistoryService {
	s.after = after
	return s
}

func (s *TradingBotGridOrdersAlgoHistoryService) Before(before string) *TradingBotGridOrdersAlgoHistoryService {
	s.before = before
	return s
}

func (s *TradingBotGridOrdersAlgoHistoryService) Limit(limit int) *TradingBotGridOrdersAlgoHistoryService {
	s.limit = &limit
	return s
}

var errTradingBotGridOrdersAlgoHistoryMissingAlgoOrdType = errors.New("okx: tradingBot grid orders-algo-history requires algoOrdType")

// Do 获取历史网格策略委托单列表（GET /api/v5/tradingBot/grid/orders-algo-history）。
func (s *TradingBotGridOrdersAlgoHistoryService) Do(ctx context.Context) ([]TradingBotGridOrder, error) {
	if s.algoOrdType == "" {
		return nil, errTradingBotGridOrdersAlgoHistoryMissingAlgoOrdType
	}

	q := url.Values{}
	q.Set("algoOrdType", s.algoOrdType)
	if s.algoId != "" {
		q.Set("algoId", s.algoId)
	}
	if s.instId != "" {
		q.Set("instId", s.instId)
	}
	if s.instType != "" {
		q.Set("instType", s.instType)
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

	var data []TradingBotGridOrder
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/tradingBot/grid/orders-algo-history", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
