package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// MarketOptionInstrumentFamilyTradeInfo 表示期权品种成交数据项。
// 数值字段保持为 string（无损）。
type MarketOptionInstrumentFamilyTradeInfo struct {
	InstId  string `json:"instId"`
	TradeId string `json:"tradeId"`

	Px   string `json:"px"`
	Sz   string `json:"sz"`
	Side string `json:"side"`

	TS int64 `json:"ts,string"`
}

// MarketOptionInstrumentFamilyTrades 表示期权品种公共成交数据（按期权类型聚合）。
// 数值字段保持为 string（无损）。
type MarketOptionInstrumentFamilyTrades struct {
	Vol24h  string `json:"vol24h"`
	OptType string `json:"optType"`

	TradeInfo []MarketOptionInstrumentFamilyTradeInfo `json:"tradeInfo"`
}

// MarketOptionInstrumentFamilyTradesService 获取期权同一交易品种下的公共成交数据。
type MarketOptionInstrumentFamilyTradesService struct {
	c *Client

	instFamily string
}

// NewMarketOptionInstrumentFamilyTradesService 创建 MarketOptionInstrumentFamilyTradesService。
func (c *Client) NewMarketOptionInstrumentFamilyTradesService() *MarketOptionInstrumentFamilyTradesService {
	return &MarketOptionInstrumentFamilyTradesService{c: c}
}

// InstFamily 设置交易品种（必填），如 BTC-USD（期权）。
func (s *MarketOptionInstrumentFamilyTradesService) InstFamily(instFamily string) *MarketOptionInstrumentFamilyTradesService {
	s.instFamily = instFamily
	return s
}

var errMarketOptionInstrumentFamilyTradesMissingInstFamily = errors.New("okx: market option instrument family trades requires instFamily")

// Do 获取期权品种公共成交数据（GET /api/v5/market/option/instrument-family-trades）。
func (s *MarketOptionInstrumentFamilyTradesService) Do(ctx context.Context) ([]MarketOptionInstrumentFamilyTrades, error) {
	if s.instFamily == "" {
		return nil, errMarketOptionInstrumentFamilyTradesMissingInstFamily
	}

	q := url.Values{}
	q.Set("instFamily", s.instFamily)

	var data []MarketOptionInstrumentFamilyTrades
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/market/option/instrument-family-trades", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
