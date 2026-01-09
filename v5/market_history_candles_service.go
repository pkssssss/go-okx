package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// MarketHistoryCandlesService 获取交易产品历史K线数据。
type MarketHistoryCandlesService struct {
	c *Client

	instId string
	bar    string
	after  string
	before string
	limit  *int
}

// NewMarketHistoryCandlesService 创建 MarketHistoryCandlesService。
func (c *Client) NewMarketHistoryCandlesService() *MarketHistoryCandlesService {
	return &MarketHistoryCandlesService{c: c}
}

// InstId 设置产品 ID（必填）。
func (s *MarketHistoryCandlesService) InstId(instId string) *MarketHistoryCandlesService {
	s.instId = instId
	return s
}

// Bar 设置 K 线粒度，如 1s/1m/5m/1H/1D 等。
func (s *MarketHistoryCandlesService) Bar(bar string) *MarketHistoryCandlesService {
	s.bar = bar
	return s
}

// After 设置请求此时间戳之前（更旧的数据）的分页内容（毫秒字符串）。
func (s *MarketHistoryCandlesService) After(after string) *MarketHistoryCandlesService {
	s.after = after
	return s
}

// Before 设置请求此时间戳之后（更新的数据）的分页内容（毫秒字符串）。
func (s *MarketHistoryCandlesService) Before(before string) *MarketHistoryCandlesService {
	s.before = before
	return s
}

// Limit 设置返回条数（最大 300，默认 100）。
func (s *MarketHistoryCandlesService) Limit(limit int) *MarketHistoryCandlesService {
	s.limit = &limit
	return s
}

var errMarketHistoryCandlesMissingInstId = errors.New("okx: market history candles requires instId")

// Do 获取交易产品历史K线数据（GET /api/v5/market/history-candles）。
func (s *MarketHistoryCandlesService) Do(ctx context.Context) ([]Candle, error) {
	if s.instId == "" {
		return nil, errMarketHistoryCandlesMissingInstId
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

	var data []Candle
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/market/history-candles", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
