package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// MarketIndexCandlesService 获取指数K线数据。
type MarketIndexCandlesService struct {
	c *Client

	instId string
	bar    string
	after  string
	before string
	limit  *int
}

// NewMarketIndexCandlesService 创建 MarketIndexCandlesService。
func (c *Client) NewMarketIndexCandlesService() *MarketIndexCandlesService {
	return &MarketIndexCandlesService{c: c}
}

// InstId 设置现货指数（必填），如 BTC-USD。
func (s *MarketIndexCandlesService) InstId(instId string) *MarketIndexCandlesService {
	s.instId = instId
	return s
}

// After 设置请求此时间戳之前（更旧的数据）的分页内容（毫秒字符串）。
func (s *MarketIndexCandlesService) After(after string) *MarketIndexCandlesService {
	s.after = after
	return s
}

// Before 设置请求此时间戳之后（更新的数据）的分页内容（毫秒字符串）。
func (s *MarketIndexCandlesService) Before(before string) *MarketIndexCandlesService {
	s.before = before
	return s
}

// Bar 设置时间粒度（默认 1m）。
func (s *MarketIndexCandlesService) Bar(bar string) *MarketIndexCandlesService {
	s.bar = bar
	return s
}

// Limit 设置返回条数（最大 100，默认 100）。
func (s *MarketIndexCandlesService) Limit(limit int) *MarketIndexCandlesService {
	s.limit = &limit
	return s
}

var errMarketIndexCandlesMissingInstId = errors.New("okx: market index candles requires instId")

// Do 获取指数K线数据（GET /api/v5/market/index-candles）。
func (s *MarketIndexCandlesService) Do(ctx context.Context) ([]PriceCandle, error) {
	if s.instId == "" {
		return nil, errMarketIndexCandlesMissingInstId
	}

	q := url.Values{}
	q.Set("instId", s.instId)
	if s.after != "" {
		q.Set("after", s.after)
	}
	if s.before != "" {
		q.Set("before", s.before)
	}
	if s.bar != "" {
		q.Set("bar", s.bar)
	}
	if s.limit != nil {
		q.Set("limit", strconv.Itoa(*s.limit))
	}

	var data []PriceCandle
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/market/index-candles", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
