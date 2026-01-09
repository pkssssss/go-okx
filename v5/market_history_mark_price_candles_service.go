package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// MarketHistoryMarkPriceCandlesService 获取标记价格历史K线数据。
type MarketHistoryMarkPriceCandlesService struct {
	c *Client

	instId string
	bar    string
	after  string
	before string
	limit  *int
}

// NewMarketHistoryMarkPriceCandlesService 创建 MarketHistoryMarkPriceCandlesService。
func (c *Client) NewMarketHistoryMarkPriceCandlesService() *MarketHistoryMarkPriceCandlesService {
	return &MarketHistoryMarkPriceCandlesService{c: c}
}

// InstId 设置产品 ID（必填）。
func (s *MarketHistoryMarkPriceCandlesService) InstId(instId string) *MarketHistoryMarkPriceCandlesService {
	s.instId = instId
	return s
}

// After 设置请求此时间戳之前（更旧的数据）的分页内容（毫秒字符串）。
func (s *MarketHistoryMarkPriceCandlesService) After(after string) *MarketHistoryMarkPriceCandlesService {
	s.after = after
	return s
}

// Before 设置请求此时间戳之后（更新的数据）的分页内容（毫秒字符串）。
func (s *MarketHistoryMarkPriceCandlesService) Before(before string) *MarketHistoryMarkPriceCandlesService {
	s.before = before
	return s
}

// Bar 设置时间粒度（默认 1m）。
func (s *MarketHistoryMarkPriceCandlesService) Bar(bar string) *MarketHistoryMarkPriceCandlesService {
	s.bar = bar
	return s
}

// Limit 设置返回条数（最大 100，默认 100）。
func (s *MarketHistoryMarkPriceCandlesService) Limit(limit int) *MarketHistoryMarkPriceCandlesService {
	s.limit = &limit
	return s
}

var errMarketHistoryMarkPriceCandlesMissingInstId = errors.New("okx: market history mark price candles requires instId")

// Do 获取标记价格历史K线数据（GET /api/v5/market/history-mark-price-candles）。
func (s *MarketHistoryMarkPriceCandlesService) Do(ctx context.Context) ([]PriceCandle, error) {
	if s.instId == "" {
		return nil, errMarketHistoryMarkPriceCandlesMissingInstId
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
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/market/history-mark-price-candles", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
