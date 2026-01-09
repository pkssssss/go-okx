package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// MarketHistoryTradesService 获取交易产品公共历史成交数据。
type MarketHistoryTradesService struct {
	c *Client

	instId string

	paginationType string
	after          string
	before         string
	limit          *int
}

// NewMarketHistoryTradesService 创建 MarketHistoryTradesService。
func (c *Client) NewMarketHistoryTradesService() *MarketHistoryTradesService {
	return &MarketHistoryTradesService{c: c}
}

// InstId 设置产品 ID（必填）。
func (s *MarketHistoryTradesService) InstId(instId string) *MarketHistoryTradesService {
	s.instId = instId
	return s
}

// Type 设置分页类型（可选；1: tradeId 分页；2: 时间戳分页）。
func (s *MarketHistoryTradesService) Type(paginationType string) *MarketHistoryTradesService {
	s.paginationType = paginationType
	return s
}

// After 设置请求此 tradeId 或 ts 之前的分页内容。
func (s *MarketHistoryTradesService) After(after string) *MarketHistoryTradesService {
	s.after = after
	return s
}

// Before 设置请求此 tradeId 之后（更新的数据）的分页内容（不支持时间戳分页）。
func (s *MarketHistoryTradesService) Before(before string) *MarketHistoryTradesService {
	s.before = before
	return s
}

// Limit 设置返回条数（最大 100，默认 100）。
func (s *MarketHistoryTradesService) Limit(limit int) *MarketHistoryTradesService {
	s.limit = &limit
	return s
}

var errMarketHistoryTradesMissingInstId = errors.New("okx: market history trades requires instId")

// Do 获取交易产品公共历史成交数据（GET /api/v5/market/history-trades）。
func (s *MarketHistoryTradesService) Do(ctx context.Context) ([]MarketTrade, error) {
	if s.instId == "" {
		return nil, errMarketHistoryTradesMissingInstId
	}

	q := url.Values{}
	q.Set("instId", s.instId)
	if s.paginationType != "" {
		q.Set("type", s.paginationType)
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

	var data []MarketTrade
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/market/history-trades", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
