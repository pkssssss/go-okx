package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// MarketTicker 表示产品行情。
// 数值字段保持为 string（无损）。
type MarketTicker struct {
	InstType string `json:"instType"`
	InstId   string `json:"instId"`

	Last   string `json:"last"`
	LastSz string `json:"lastSz"`

	AskPx string `json:"askPx"`
	AskSz string `json:"askSz"`
	BidPx string `json:"bidPx"`
	BidSz string `json:"bidSz"`

	Open24h   string `json:"open24h"`
	High24h   string `json:"high24h"`
	Low24h    string `json:"low24h"`
	VolCcy24h string `json:"volCcy24h"`
	Vol24h    string `json:"vol24h"`

	SodUtc0 string `json:"sodUtc0"`
	SodUtc8 string `json:"sodUtc8"`

	TS int64 `json:"ts,string"`
}

// MarketTickerService 获取单个产品行情。
type MarketTickerService struct {
	c      *Client
	instId string
}

// NewMarketTickerService 创建 MarketTickerService。
func (c *Client) NewMarketTickerService() *MarketTickerService {
	return &MarketTickerService{c: c}
}

// InstId 设置产品 ID。
func (s *MarketTickerService) InstId(instId string) *MarketTickerService {
	s.instId = instId
	return s
}

var (
	errMarketTickerMissingInstId = errors.New("okx: market ticker requires instId")
	errEmptyMarketTickerResponse = errors.New("okx: empty market ticker response")
)

// Do 获取单个产品行情（GET /api/v5/market/ticker）。
func (s *MarketTickerService) Do(ctx context.Context) (*MarketTicker, error) {
	if s.instId == "" {
		return nil, errMarketTickerMissingInstId
	}
	q := url.Values{}
	q.Set("instId", s.instId)

	var data []MarketTicker
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/market/ticker", q, nil, false, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyMarketTickerResponse
	}
	return &data[0], nil
}
