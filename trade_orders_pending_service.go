package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// OrdersPendingService 获取未成交订单列表。
type OrdersPendingService struct {
	c *Client

	instType   string
	instFamily string
	instId     string
	ordType    string
	state      string
	after      string
	before     string
	limit      *int
}

// NewOrdersPendingService 创建 OrdersPendingService。
func (c *Client) NewOrdersPendingService() *OrdersPendingService {
	return &OrdersPendingService{c: c}
}

func (s *OrdersPendingService) InstType(instType string) *OrdersPendingService {
	s.instType = instType
	return s
}

func (s *OrdersPendingService) InstFamily(instFamily string) *OrdersPendingService {
	s.instFamily = instFamily
	return s
}

func (s *OrdersPendingService) InstId(instId string) *OrdersPendingService {
	s.instId = instId
	return s
}

// OrdType 设置订单类型（可用逗号分隔多个）。
func (s *OrdersPendingService) OrdType(ordType string) *OrdersPendingService {
	s.ordType = ordType
	return s
}

func (s *OrdersPendingService) State(state string) *OrdersPendingService {
	s.state = state
	return s
}

func (s *OrdersPendingService) After(after string) *OrdersPendingService {
	s.after = after
	return s
}

func (s *OrdersPendingService) Before(before string) *OrdersPendingService {
	s.before = before
	return s
}

func (s *OrdersPendingService) Limit(limit int) *OrdersPendingService {
	s.limit = &limit
	return s
}

// Do 获取未成交订单列表（GET /api/v5/trade/orders-pending）。
func (s *OrdersPendingService) Do(ctx context.Context) ([]TradeOrder, error) {
	q := url.Values{}
	if s.instType != "" {
		q.Set("instType", s.instType)
	}
	if s.instFamily != "" {
		q.Set("instFamily", s.instFamily)
	}
	if s.instId != "" {
		q.Set("instId", s.instId)
	}
	if s.ordType != "" {
		q.Set("ordType", s.ordType)
	}
	if s.state != "" {
		q.Set("state", s.state)
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

	var data []TradeOrder
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/trade/orders-pending", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
