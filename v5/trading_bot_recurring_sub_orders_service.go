package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// TradingBotRecurringSubOrdersService 获取定投策略子订单信息。
type TradingBotRecurringSubOrdersService struct {
	c *Client

	algoId string
	ordId  string
	after  string
	before string
	limit  *int
}

// NewTradingBotRecurringSubOrdersService 创建 TradingBotRecurringSubOrdersService。
func (c *Client) NewTradingBotRecurringSubOrdersService() *TradingBotRecurringSubOrdersService {
	return &TradingBotRecurringSubOrdersService{c: c}
}

func (s *TradingBotRecurringSubOrdersService) AlgoId(algoId string) *TradingBotRecurringSubOrdersService {
	s.algoId = algoId
	return s
}

func (s *TradingBotRecurringSubOrdersService) OrdId(ordId string) *TradingBotRecurringSubOrdersService {
	s.ordId = ordId
	return s
}

func (s *TradingBotRecurringSubOrdersService) After(after string) *TradingBotRecurringSubOrdersService {
	s.after = after
	return s
}

func (s *TradingBotRecurringSubOrdersService) Before(before string) *TradingBotRecurringSubOrdersService {
	s.before = before
	return s
}

func (s *TradingBotRecurringSubOrdersService) Limit(limit int) *TradingBotRecurringSubOrdersService {
	s.limit = &limit
	return s
}

var errTradingBotRecurringSubOrdersMissingAlgoId = errors.New("okx: tradingBot recurring sub-orders requires algoId")

// Do 获取定投策略子订单信息（GET /api/v5/tradingBot/recurring/sub-orders）。
func (s *TradingBotRecurringSubOrdersService) Do(ctx context.Context) ([]TradingBotRecurringSubOrder, error) {
	if s.algoId == "" {
		return nil, errTradingBotRecurringSubOrdersMissingAlgoId
	}

	q := url.Values{}
	q.Set("algoId", s.algoId)
	if s.ordId != "" {
		q.Set("ordId", s.ordId)
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

	var data []TradingBotRecurringSubOrder
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/tradingBot/recurring/sub-orders", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
