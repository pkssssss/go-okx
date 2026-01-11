package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// MarketSprdHistoryCandlesService 获取价差交易产品历史 K 线数据。
type MarketSprdHistoryCandlesService struct {
	c *Client

	sprdId string
	bar    string
	after  string
	before string
	limit  *int
}

// NewMarketSprdHistoryCandlesService 创建 MarketSprdHistoryCandlesService。
func (c *Client) NewMarketSprdHistoryCandlesService() *MarketSprdHistoryCandlesService {
	return &MarketSprdHistoryCandlesService{c: c}
}

// SprdId 设置 Spread ID（必填）。
func (s *MarketSprdHistoryCandlesService) SprdId(sprdId string) *MarketSprdHistoryCandlesService {
	s.sprdId = sprdId
	return s
}

// After 设置请求此时间戳之前（更旧的数据）的分页内容（毫秒字符串）。
func (s *MarketSprdHistoryCandlesService) After(after string) *MarketSprdHistoryCandlesService {
	s.after = after
	return s
}

// Before 设置请求此时间戳之后（更新的数据）的分页内容（毫秒字符串）。
func (s *MarketSprdHistoryCandlesService) Before(before string) *MarketSprdHistoryCandlesService {
	s.before = before
	return s
}

// Bar 设置时间粒度（默认 1m）。
func (s *MarketSprdHistoryCandlesService) Bar(bar string) *MarketSprdHistoryCandlesService {
	s.bar = bar
	return s
}

// Limit 设置返回条数（最大 100，默认 100）。
func (s *MarketSprdHistoryCandlesService) Limit(limit int) *MarketSprdHistoryCandlesService {
	s.limit = &limit
	return s
}

var errMarketSprdHistoryCandlesMissingSprdId = errors.New("okx: market sprd history candles requires sprdId")

// Do 获取价差交易产品历史 K 线数据（GET /api/v5/market/sprd-history-candles）。
func (s *MarketSprdHistoryCandlesService) Do(ctx context.Context) ([]Candle, error) {
	if s.sprdId == "" {
		return nil, errMarketSprdHistoryCandlesMissingSprdId
	}

	q := url.Values{}
	q.Set("sprdId", s.sprdId)
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
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/market/sprd-history-candles", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
