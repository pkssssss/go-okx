package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// TradeFillsService 获取成交明细（近三天）。
type TradeFillsService struct {
	c *Client

	instType   string
	instFamily string
	instId     string
	ordId      string
	after      string
	before     string
	limit      *int
}

// NewTradeFillsService 创建 TradeFillsService。
func (c *Client) NewTradeFillsService() *TradeFillsService {
	return &TradeFillsService{c: c}
}

func (s *TradeFillsService) InstType(instType string) *TradeFillsService {
	s.instType = instType
	return s
}

func (s *TradeFillsService) InstFamily(instFamily string) *TradeFillsService {
	s.instFamily = instFamily
	return s
}

func (s *TradeFillsService) InstId(instId string) *TradeFillsService {
	s.instId = instId
	return s
}

func (s *TradeFillsService) OrdId(ordId string) *TradeFillsService {
	s.ordId = ordId
	return s
}

func (s *TradeFillsService) After(after string) *TradeFillsService {
	s.after = after
	return s
}

func (s *TradeFillsService) Before(before string) *TradeFillsService {
	s.before = before
	return s
}

func (s *TradeFillsService) Limit(limit int) *TradeFillsService {
	s.limit = &limit
	return s
}

// Do 获取成交明细（GET /api/v5/trade/fills）。
func (s *TradeFillsService) Do(ctx context.Context) ([]TradeFill, error) {
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

	var data []TradeFill
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/trade/fills", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
