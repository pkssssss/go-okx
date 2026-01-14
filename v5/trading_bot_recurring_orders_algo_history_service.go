package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// TradingBotRecurringOrdersAlgoHistoryService 获取历史定投策略委托单列表。
type TradingBotRecurringOrdersAlgoHistoryService struct {
	c *Client

	algoId string
	after  string
	before string
	limit  *int
}

// NewTradingBotRecurringOrdersAlgoHistoryService 创建 TradingBotRecurringOrdersAlgoHistoryService。
func (c *Client) NewTradingBotRecurringOrdersAlgoHistoryService() *TradingBotRecurringOrdersAlgoHistoryService {
	return &TradingBotRecurringOrdersAlgoHistoryService{c: c}
}

func (s *TradingBotRecurringOrdersAlgoHistoryService) AlgoId(algoId string) *TradingBotRecurringOrdersAlgoHistoryService {
	s.algoId = algoId
	return s
}

func (s *TradingBotRecurringOrdersAlgoHistoryService) After(after string) *TradingBotRecurringOrdersAlgoHistoryService {
	s.after = after
	return s
}

func (s *TradingBotRecurringOrdersAlgoHistoryService) Before(before string) *TradingBotRecurringOrdersAlgoHistoryService {
	s.before = before
	return s
}

func (s *TradingBotRecurringOrdersAlgoHistoryService) Limit(limit int) *TradingBotRecurringOrdersAlgoHistoryService {
	s.limit = &limit
	return s
}

// Do 获取历史定投策略委托单列表（GET /api/v5/tradingBot/recurring/orders-algo-history）。
func (s *TradingBotRecurringOrdersAlgoHistoryService) Do(ctx context.Context) ([]TradingBotRecurringOrder, error) {
	var q url.Values
	if s.algoId != "" || s.after != "" || s.before != "" || s.limit != nil {
		q = url.Values{}
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
	}

	var data []TradingBotRecurringOrder
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/tradingBot/recurring/orders-algo-history", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
