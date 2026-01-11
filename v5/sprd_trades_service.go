package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// SprdTradesService 获取历史成交数据（近七天）。
type SprdTradesService struct {
	c *Client

	sprdId  string
	tradeId string
	ordId   string

	beginId string
	endId   string

	begin string
	end   string

	limit *int
}

// NewSprdTradesService 创建 SprdTradesService。
func (c *Client) NewSprdTradesService() *SprdTradesService {
	return &SprdTradesService{c: c}
}

func (s *SprdTradesService) SprdId(sprdId string) *SprdTradesService {
	s.sprdId = sprdId
	return s
}

func (s *SprdTradesService) TradeId(tradeId string) *SprdTradesService {
	s.tradeId = tradeId
	return s
}

func (s *SprdTradesService) OrdId(ordId string) *SprdTradesService {
	s.ordId = ordId
	return s
}

func (s *SprdTradesService) BeginId(beginId string) *SprdTradesService {
	s.beginId = beginId
	return s
}

func (s *SprdTradesService) EndId(endId string) *SprdTradesService {
	s.endId = endId
	return s
}

func (s *SprdTradesService) Begin(begin string) *SprdTradesService {
	s.begin = begin
	return s
}

func (s *SprdTradesService) End(end string) *SprdTradesService {
	s.end = end
	return s
}

func (s *SprdTradesService) Limit(limit int) *SprdTradesService {
	s.limit = &limit
	return s
}

// Do 获取历史成交数据（GET /api/v5/sprd/trades）。
func (s *SprdTradesService) Do(ctx context.Context) ([]SprdTrade, error) {
	q := url.Values{}
	if s.sprdId != "" {
		q.Set("sprdId", s.sprdId)
	}
	if s.tradeId != "" {
		q.Set("tradeId", s.tradeId)
	}
	if s.ordId != "" {
		q.Set("ordId", s.ordId)
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

	var data []SprdTrade
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/sprd/trades", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
