package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// MarketSprdTicker 表示价差交易（Spread Trading）产品行情。
// 数值字段保持为 string（无损）。
type MarketSprdTicker struct {
	SprdId string `json:"sprdId"`

	Last   string `json:"last"`
	LastSz string `json:"lastSz"`

	AskPx string `json:"askPx"`
	AskSz string `json:"askSz"`
	BidPx string `json:"bidPx"`
	BidSz string `json:"bidSz"`

	Open24h string `json:"open24h"`
	High24h string `json:"high24h"`
	Low24h  string `json:"low24h"`
	Vol24h  string `json:"vol24h"`

	TS int64 `json:"ts,string"`
}

// MarketSprdTickerService 获取单个价差交易产品行情。
type MarketSprdTickerService struct {
	c *Client

	sprdId string
}

// NewMarketSprdTickerService 创建 MarketSprdTickerService。
func (c *Client) NewMarketSprdTickerService() *MarketSprdTickerService {
	return &MarketSprdTickerService{c: c}
}

// SprdId 设置 Spread ID（必填），如 BTC-USDT_BTC-USDT-SWAP。
func (s *MarketSprdTickerService) SprdId(sprdId string) *MarketSprdTickerService {
	s.sprdId = sprdId
	return s
}

var (
	errMarketSprdTickerMissingSprdId = errors.New("okx: market sprd ticker requires sprdId")
	errEmptyMarketSprdTickerResponse = errors.New("okx: empty market sprd ticker response")
)

// Do 获取单个价差交易产品行情（GET /api/v5/market/sprd-ticker）。
func (s *MarketSprdTickerService) Do(ctx context.Context) (*MarketSprdTicker, error) {
	if s.sprdId == "" {
		return nil, errMarketSprdTickerMissingSprdId
	}

	q := url.Values{}
	q.Set("sprdId", s.sprdId)

	var data []MarketSprdTicker
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/market/sprd-ticker", q, nil, false, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyMarketSprdTickerResponse
	}
	return &data[0], nil
}
