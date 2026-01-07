package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// MarketTrade 表示最近成交。
type MarketTrade struct {
	InstId  string `json:"instId"`
	TradeId string `json:"tradeId"`

	Px   string `json:"px"`
	Sz   string `json:"sz"`
	Side string `json:"side"`

	TS int64 `json:"ts,string"`
}

// MarketTradesService 获取最近成交。
type MarketTradesService struct {
	c *Client

	instId string
	limit  *int
}

// NewMarketTradesService 创建 MarketTradesService。
func (c *Client) NewMarketTradesService() *MarketTradesService {
	return &MarketTradesService{c: c}
}

// InstId 设置产品 ID（必填）。
func (s *MarketTradesService) InstId(instId string) *MarketTradesService {
	s.instId = instId
	return s
}

// Limit 设置返回条数。
func (s *MarketTradesService) Limit(limit int) *MarketTradesService {
	s.limit = &limit
	return s
}

var errMarketTradesMissingInstId = errors.New("okx: market trades requires instId")

// Do 获取最近成交（GET /api/v5/market/trades）。
func (s *MarketTradesService) Do(ctx context.Context) ([]MarketTrade, error) {
	if s.instId == "" {
		return nil, errMarketTradesMissingInstId
	}

	q := url.Values{}
	q.Set("instId", s.instId)
	if s.limit != nil {
		q.Set("limit", strconv.Itoa(*s.limit))
	}

	var data []MarketTrade
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/market/trades", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
