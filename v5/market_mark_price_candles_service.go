package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// MarketMarkPriceCandlesService 获取标记价格K线数据。
type MarketMarkPriceCandlesService struct {
	c *Client

	instId string
	bar    string
	after  string
	before string
	limit  *int
}

// NewMarketMarkPriceCandlesService 创建 MarketMarkPriceCandlesService。
func (c *Client) NewMarketMarkPriceCandlesService() *MarketMarkPriceCandlesService {
	return &MarketMarkPriceCandlesService{c: c}
}

// InstId 设置产品 ID（必填）。
func (s *MarketMarkPriceCandlesService) InstId(instId string) *MarketMarkPriceCandlesService {
	s.instId = instId
	return s
}

// After 设置请求此时间戳之前（更旧的数据）的分页内容（毫秒字符串）。
func (s *MarketMarkPriceCandlesService) After(after string) *MarketMarkPriceCandlesService {
	s.after = after
	return s
}

// Before 设置请求此时间戳之后（更新的数据）的分页内容（毫秒字符串）。
func (s *MarketMarkPriceCandlesService) Before(before string) *MarketMarkPriceCandlesService {
	s.before = before
	return s
}

// Bar 设置时间粒度（默认 1m）。
func (s *MarketMarkPriceCandlesService) Bar(bar string) *MarketMarkPriceCandlesService {
	s.bar = bar
	return s
}

// Limit 设置返回条数（最大 100，默认 100）。
func (s *MarketMarkPriceCandlesService) Limit(limit int) *MarketMarkPriceCandlesService {
	s.limit = &limit
	return s
}

var errMarketMarkPriceCandlesMissingInstId = errors.New("okx: market mark price candles requires instId")

// Do 获取标记价格K线数据（GET /api/v5/market/mark-price-candles）。
func (s *MarketMarkPriceCandlesService) Do(ctx context.Context) ([]PriceCandle, error) {
	if s.instId == "" {
		return nil, errMarketMarkPriceCandlesMissingInstId
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
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/market/mark-price-candles", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
