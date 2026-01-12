package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// MarketBlockTicker 表示大宗交易产品行情（最近24小时成交量信息）。
// 数值字段保持为 string（无损）。
type MarketBlockTicker struct {
	InstType string `json:"instType"`
	InstId   string `json:"instId"`

	VolCcy24h string `json:"volCcy24h"`
	Vol24h    string `json:"vol24h"`

	TS int64 `json:"ts,string"`
}

// MarketBlockTickerService 获取大宗交易单个产品行情信息。
type MarketBlockTickerService struct {
	c      *Client
	instId string
}

// NewMarketBlockTickerService 创建 MarketBlockTickerService。
func (c *Client) NewMarketBlockTickerService() *MarketBlockTickerService {
	return &MarketBlockTickerService{c: c}
}

// InstId 设置产品 ID（必填）。
func (s *MarketBlockTickerService) InstId(instId string) *MarketBlockTickerService {
	s.instId = instId
	return s
}

var (
	errMarketBlockTickerMissingInstId = errors.New("okx: market block ticker requires instId")
	errEmptyMarketBlockTickerResponse = errors.New("okx: empty market block ticker response")
)

// Do 获取大宗交易单个产品行情信息（GET /api/v5/market/block-ticker）。
func (s *MarketBlockTickerService) Do(ctx context.Context) (*MarketBlockTicker, error) {
	if s.instId == "" {
		return nil, errMarketBlockTickerMissingInstId
	}

	q := url.Values{}
	q.Set("instId", s.instId)

	var data []MarketBlockTicker
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/market/block-ticker", q, nil, false, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyMarketBlockTickerResponse
	}
	return &data[0], nil
}
