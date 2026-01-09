package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// MarketHistoryIndexCandlesService 获取指数历史K线数据。
type MarketHistoryIndexCandlesService struct {
	c *Client

	instId string
	bar    string
	after  string
	before string
	limit  *int
}

// NewMarketHistoryIndexCandlesService 创建 MarketHistoryIndexCandlesService。
func (c *Client) NewMarketHistoryIndexCandlesService() *MarketHistoryIndexCandlesService {
	return &MarketHistoryIndexCandlesService{c: c}
}

// InstId 设置现货指数（必填），如 BTC-USD。
func (s *MarketHistoryIndexCandlesService) InstId(instId string) *MarketHistoryIndexCandlesService {
	s.instId = instId
	return s
}

// After 设置请求此时间戳之前（更旧的数据）的分页内容（毫秒字符串）。
func (s *MarketHistoryIndexCandlesService) After(after string) *MarketHistoryIndexCandlesService {
	s.after = after
	return s
}

// Before 设置请求此时间戳之后（更新的数据）的分页内容（毫秒字符串）。
func (s *MarketHistoryIndexCandlesService) Before(before string) *MarketHistoryIndexCandlesService {
	s.before = before
	return s
}

// Bar 设置时间粒度（默认 1m）。
func (s *MarketHistoryIndexCandlesService) Bar(bar string) *MarketHistoryIndexCandlesService {
	s.bar = bar
	return s
}

// Limit 设置返回条数（最大 100，默认 100）。
func (s *MarketHistoryIndexCandlesService) Limit(limit int) *MarketHistoryIndexCandlesService {
	s.limit = &limit
	return s
}

var errMarketHistoryIndexCandlesMissingInstId = errors.New("okx: market history index candles requires instId")

// Do 获取指数历史K线数据（GET /api/v5/market/history-index-candles）。
func (s *MarketHistoryIndexCandlesService) Do(ctx context.Context) ([]PriceCandle, error) {
	if s.instId == "" {
		return nil, errMarketHistoryIndexCandlesMissingInstId
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
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/market/history-index-candles", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
