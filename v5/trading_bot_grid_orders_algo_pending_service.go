package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// TradingBotGridOrdersAlgoPendingService 获取未完成网格策略委托单列表。
type TradingBotGridOrdersAlgoPendingService struct {
	c *Client

	algoOrdType string
	algoId      string
	instId      string
	instType    string
	after       string
	before      string
	limit       *int
}

// NewTradingBotGridOrdersAlgoPendingService 创建 TradingBotGridOrdersAlgoPendingService。
func (c *Client) NewTradingBotGridOrdersAlgoPendingService() *TradingBotGridOrdersAlgoPendingService {
	return &TradingBotGridOrdersAlgoPendingService{c: c}
}

func (s *TradingBotGridOrdersAlgoPendingService) AlgoOrdType(algoOrdType string) *TradingBotGridOrdersAlgoPendingService {
	s.algoOrdType = algoOrdType
	return s
}

func (s *TradingBotGridOrdersAlgoPendingService) AlgoId(algoId string) *TradingBotGridOrdersAlgoPendingService {
	s.algoId = algoId
	return s
}

func (s *TradingBotGridOrdersAlgoPendingService) InstId(instId string) *TradingBotGridOrdersAlgoPendingService {
	s.instId = instId
	return s
}

func (s *TradingBotGridOrdersAlgoPendingService) InstType(instType string) *TradingBotGridOrdersAlgoPendingService {
	s.instType = instType
	return s
}

func (s *TradingBotGridOrdersAlgoPendingService) After(after string) *TradingBotGridOrdersAlgoPendingService {
	s.after = after
	return s
}

func (s *TradingBotGridOrdersAlgoPendingService) Before(before string) *TradingBotGridOrdersAlgoPendingService {
	s.before = before
	return s
}

func (s *TradingBotGridOrdersAlgoPendingService) Limit(limit int) *TradingBotGridOrdersAlgoPendingService {
	s.limit = &limit
	return s
}

var errTradingBotGridOrdersAlgoPendingMissingAlgoOrdType = errors.New("okx: tradingBot grid orders-algo-pending requires algoOrdType")

// Do 获取未完成网格策略委托单列表（GET /api/v5/tradingBot/grid/orders-algo-pending）。
func (s *TradingBotGridOrdersAlgoPendingService) Do(ctx context.Context) ([]TradingBotGridOrder, error) {
	if s.algoOrdType == "" {
		return nil, errTradingBotGridOrdersAlgoPendingMissingAlgoOrdType
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
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/tradingBot/grid/orders-algo-pending", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
