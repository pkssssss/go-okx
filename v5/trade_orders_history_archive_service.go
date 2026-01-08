package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// OrdersHistoryArchiveService 获取历史订单记录（近三个月）。
type OrdersHistoryArchiveService struct {
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

// NewOrdersHistoryArchiveService 创建 OrdersHistoryArchiveService。
func (c *Client) NewOrdersHistoryArchiveService() *OrdersHistoryArchiveService {
	return &OrdersHistoryArchiveService{c: c}
}

// InstType 设置产品类型（必填）。
func (s *OrdersHistoryArchiveService) InstType(instType string) *OrdersHistoryArchiveService {
	s.instType = instType
	return s
}

func (s *OrdersHistoryArchiveService) InstFamily(instFamily string) *OrdersHistoryArchiveService {
	s.instFamily = instFamily
	return s
}

func (s *OrdersHistoryArchiveService) InstId(instId string) *OrdersHistoryArchiveService {
	s.instId = instId
	return s
}

// OrdType 设置订单类型（可用逗号分隔多个）。
func (s *OrdersHistoryArchiveService) OrdType(ordType string) *OrdersHistoryArchiveService {
	s.ordType = ordType
	return s
}

func (s *OrdersHistoryArchiveService) State(state string) *OrdersHistoryArchiveService {
	s.state = state
	return s
}

// Category 设置订单种类（如 twap/adl/...）。
func (s *OrdersHistoryArchiveService) Category(category string) *OrdersHistoryArchiveService {
	s.category = category
	return s
}

func (s *OrdersHistoryArchiveService) After(after string) *OrdersHistoryArchiveService {
	s.after = after
	return s
}

func (s *OrdersHistoryArchiveService) Before(before string) *OrdersHistoryArchiveService {
	s.before = before
	return s
}

// Begin 设置开始时间（Unix 毫秒时间戳字符串）。
func (s *OrdersHistoryArchiveService) Begin(begin string) *OrdersHistoryArchiveService {
	s.begin = begin
	return s
}

// End 设置结束时间（Unix 毫秒时间戳字符串）。
func (s *OrdersHistoryArchiveService) End(end string) *OrdersHistoryArchiveService {
	s.end = end
	return s
}

func (s *OrdersHistoryArchiveService) Limit(limit int) *OrdersHistoryArchiveService {
	s.limit = &limit
	return s
}

var errOrdersHistoryArchiveMissingInstType = errors.New("okx: orders history archive requires instType")

// Do 获取历史订单记录（近三个月）（GET /api/v5/trade/orders-history-archive）。
func (s *OrdersHistoryArchiveService) Do(ctx context.Context) ([]TradeOrder, error) {
	if s.instType == "" {
		return nil, errOrdersHistoryArchiveMissingInstType
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
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/trade/orders-history-archive", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
