package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// TradingBotRecurringOrdersAlgoPendingService 获取未完成定投策略委托单列表。
type TradingBotRecurringOrdersAlgoPendingService struct {
	c *Client

	algoId string
	after  string
	before string
	limit  *int
}

// NewTradingBotRecurringOrdersAlgoPendingService 创建 TradingBotRecurringOrdersAlgoPendingService。
func (c *Client) NewTradingBotRecurringOrdersAlgoPendingService() *TradingBotRecurringOrdersAlgoPendingService {
	return &TradingBotRecurringOrdersAlgoPendingService{c: c}
}

func (s *TradingBotRecurringOrdersAlgoPendingService) AlgoId(algoId string) *TradingBotRecurringOrdersAlgoPendingService {
	s.algoId = algoId
	return s
}

func (s *TradingBotRecurringOrdersAlgoPendingService) After(after string) *TradingBotRecurringOrdersAlgoPendingService {
	s.after = after
	return s
}

func (s *TradingBotRecurringOrdersAlgoPendingService) Before(before string) *TradingBotRecurringOrdersAlgoPendingService {
	s.before = before
	return s
}

func (s *TradingBotRecurringOrdersAlgoPendingService) Limit(limit int) *TradingBotRecurringOrdersAlgoPendingService {
	s.limit = &limit
	return s
}

// Do 获取未完成定投策略委托单列表（GET /api/v5/tradingBot/recurring/orders-algo-pending）。
func (s *TradingBotRecurringOrdersAlgoPendingService) Do(ctx context.Context) ([]TradingBotRecurringOrder, error) {
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
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/tradingBot/recurring/orders-algo-pending", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
