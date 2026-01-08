package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// OrdersHistoryService 获取历史订单记录（近七天）。
type OrdersHistoryService struct {
	c *Client

	instType   string
	instFamily string
	instId     string
	ordType    string
	state      string
	category   string
	after      string
	before     string
	begin      string
	end        string
	limit      *int
}

// NewOrdersHistoryService 创建 OrdersHistoryService。
func (c *Client) NewOrdersHistoryService() *OrdersHistoryService {
	return &OrdersHistoryService{c: c}
}

// InstType 设置产品类型（必填）。
func (s *OrdersHistoryService) InstType(instType string) *OrdersHistoryService {
	s.instType = instType
	return s
}

func (s *OrdersHistoryService) InstFamily(instFamily string) *OrdersHistoryService {
	s.instFamily = instFamily
	return s
}

func (s *OrdersHistoryService) InstId(instId string) *OrdersHistoryService {
	s.instId = instId
	return s
}

// OrdType 设置订单类型（可用逗号分隔多个）。
func (s *OrdersHistoryService) OrdType(ordType string) *OrdersHistoryService {
	s.ordType = ordType
	return s
}

func (s *OrdersHistoryService) State(state string) *OrdersHistoryService {
	s.state = state
	return s
}

// Category 设置订单种类（如 twap/adl/...）。
func (s *OrdersHistoryService) Category(category string) *OrdersHistoryService {
	s.category = category
	return s
}

func (s *OrdersHistoryService) After(after string) *OrdersHistoryService {
	s.after = after
	return s
}

func (s *OrdersHistoryService) Before(before string) *OrdersHistoryService {
	s.before = before
	return s
}

// Begin 设置开始时间（Unix 毫秒时间戳字符串）。
func (s *OrdersHistoryService) Begin(begin string) *OrdersHistoryService {
	s.begin = begin
	return s
}

// End 设置结束时间（Unix 毫秒时间戳字符串）。
func (s *OrdersHistoryService) End(end string) *OrdersHistoryService {
	s.end = end
	return s
}

func (s *OrdersHistoryService) Limit(limit int) *OrdersHistoryService {
	s.limit = &limit
	return s
}

var errOrdersHistoryMissingInstType = errors.New("okx: orders history requires instType")

// Do 获取历史订单记录（GET /api/v5/trade/orders-history）。
func (s *OrdersHistoryService) Do(ctx context.Context) ([]TradeOrder, error) {
	if s.instType == "" {
		return nil, errOrdersHistoryMissingInstType
	}

	q := url.Values{}
	q.Set("instType", s.instType)
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
	if s.category != "" {
		q.Set("category", s.category)
	}
	if s.after != "" {
		q.Set("after", s.after)
	}
	if s.before != "" {
		q.Set("before", s.before)
	}
	if s.begin != "" {
		q.Set("begin", s.begin)
	}
	if s.end != "" {
		q.Set("end", s.end)
	}
	if s.limit != nil {
		q.Set("limit", strconv.Itoa(*s.limit))
	}

	var data []TradeOrder
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/trade/orders-history", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
