package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// MarketIndexTickersService 获取指数行情数据。
type MarketIndexTickersService struct {
	c *Client

	quoteCcy string
	instId   string
}

// NewMarketIndexTickersService 创建 MarketIndexTickersService。
func (c *Client) NewMarketIndexTickersService() *MarketIndexTickersService {
	return &MarketIndexTickersService{c: c}
}

// QuoteCcy 设置指数计价单位（USD/USDT/BTC/USDC），与 InstId 必须填写一个。
func (s *MarketIndexTickersService) QuoteCcy(quoteCcy string) *MarketIndexTickersService {
	s.quoteCcy = quoteCcy
	return s
}

// InstId 设置指数，如 BTC-USD，与 QuoteCcy 必须填写一个。
func (s *MarketIndexTickersService) InstId(instId string) *MarketIndexTickersService {
	s.instId = instId
	return s
}

var errMarketIndexTickersMissingQuoteCcyOrInstId = errors.New("okx: market index tickers requires quoteCcy or instId")

// Do 获取指数行情数据（GET /api/v5/market/index-tickers）。
func (s *MarketIndexTickersService) Do(ctx context.Context) ([]IndexTicker, error) {
	if s.quoteCcy == "" && s.instId == "" {
		return nil, errMarketIndexTickersMissingQuoteCcyOrInstId
	}

	q := url.Values{}
	if s.quoteCcy != "" {
		q.Set("quoteCcy", s.quoteCcy)
	}
	if s.instId != "" {
		q.Set("instId", s.instId)
	}

	var data []IndexTicker
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/market/index-tickers", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
