package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// SprdOrdersHistoryService 获取历史订单记录（近 21 天）。
type SprdOrdersHistoryService struct {
	c *Client

	sprdId  string
	ordType string
	state   string

	beginId string
	endId   string

	begin string
	end   string

	limit *int
}

// NewSprdOrdersHistoryService 创建 SprdOrdersHistoryService。
func (c *Client) NewSprdOrdersHistoryService() *SprdOrdersHistoryService {
	return &SprdOrdersHistoryService{c: c}
}

func (s *SprdOrdersHistoryService) SprdId(sprdId string) *SprdOrdersHistoryService {
	s.sprdId = sprdId
	return s
}

func (s *SprdOrdersHistoryService) OrdType(ordType string) *SprdOrdersHistoryService {
	s.ordType = ordType
	return s
}

// State 设置订单状态（canceled/filled）。
func (s *SprdOrdersHistoryService) State(state string) *SprdOrdersHistoryService {
	s.state = state
	return s
}

func (s *SprdOrdersHistoryService) BeginId(beginId string) *SprdOrdersHistoryService {
	s.beginId = beginId
	return s
}

func (s *SprdOrdersHistoryService) EndId(endId string) *SprdOrdersHistoryService {
	s.endId = endId
	return s
}

// Begin 设置开始时间（Unix 毫秒时间戳字符串）。
func (s *SprdOrdersHistoryService) Begin(begin string) *SprdOrdersHistoryService {
	s.begin = begin
	return s
}

// End 设置结束时间（Unix 毫秒时间戳字符串）。
func (s *SprdOrdersHistoryService) End(end string) *SprdOrdersHistoryService {
	s.end = end
	return s
}

func (s *SprdOrdersHistoryService) Limit(limit int) *SprdOrdersHistoryService {
	s.limit = &limit
	return s
}

// Do 获取历史订单记录（GET /api/v5/sprd/orders-history）。
func (s *SprdOrdersHistoryService) Do(ctx context.Context) ([]SprdOrder, error) {
	q := url.Values{}
	if s.sprdId != "" {
		q.Set("sprdId", s.sprdId)
	}
	if s.ordType != "" {
		q.Set("ordType", s.ordType)
	}
	if s.state != "" {
		q.Set("state", s.state)
	}
	if s.beginId != "" {
		q.Set("beginId", s.beginId)
	}
	if s.endId != "" {
		q.Set("endId", s.endId)
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

	var data []SprdOrder
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/sprd/orders-history", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
