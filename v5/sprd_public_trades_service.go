package okx

import (
	"context"
	"net/http"
	"net/url"
)

// SprdPublicTrade 表示价差交易（Spread Trading）公共成交数据。
//
// 说明：价格/数量等字段保持为 string（无损）。
type SprdPublicTrade struct {
	SprdId  string `json:"sprdId"`
	Side    string `json:"side"`
	Sz      string `json:"sz"`
	Px      string `json:"px"`
	TradeId string `json:"tradeId"`

	TS int64 `json:"ts,string"`
}

// SprdPublicTradesService 获取价差交易公共成交数据。
type SprdPublicTradesService struct {
	c *Client

	sprdId string
}

// NewSprdPublicTradesService 创建 SprdPublicTradesService。
func (c *Client) NewSprdPublicTradesService() *SprdPublicTradesService {
	return &SprdPublicTradesService{c: c}
}

// SprdId 设置 Spread ID（可选；例如 BTC-USDT_BTC-USDT-SWAP）。
func (s *SprdPublicTradesService) SprdId(sprdId string) *SprdPublicTradesService {
	s.sprdId = sprdId
	return s
}

// Do 获取价差交易公共成交数据（GET /api/v5/sprd/public-trades）。
func (s *SprdPublicTradesService) Do(ctx context.Context) ([]SprdPublicTrade, error) {
	q := url.Values{}
	if s.sprdId != "" {
		q.Set("sprdId", s.sprdId)
	}

	var data []SprdPublicTrade
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/sprd/public-trades", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
