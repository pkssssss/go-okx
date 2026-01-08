package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// MarketTickersService 获取所有产品行情信息。
type MarketTickersService struct {
	c          *Client
	instType   string
	instFamily string
}

// NewMarketTickersService 创建 MarketTickersService。
func (c *Client) NewMarketTickersService() *MarketTickersService {
	return &MarketTickersService{c: c}
}

// InstType 设置产品类型（SPOT/SWAP/FUTURES/OPTION），必填。
func (s *MarketTickersService) InstType(instType string) *MarketTickersService {
	s.instType = instType
	return s
}

// InstFamily 设置交易品种（适用于交割/永续/期权），如 BTC-USD。
func (s *MarketTickersService) InstFamily(instFamily string) *MarketTickersService {
	s.instFamily = instFamily
	return s
}

var errMarketTickersMissingInstType = errors.New("okx: market tickers requires instType")

// Do 获取所有产品行情信息（GET /api/v5/market/tickers）。
func (s *MarketTickersService) Do(ctx context.Context) ([]MarketTicker, error) {
	if s.instType == "" {
		return nil, errMarketTickersMissingInstType
	}

	q := url.Values{}
	q.Set("instType", s.instType)
	if s.instFamily != "" {
		q.Set("instFamily", s.instFamily)
	}

	var data []MarketTicker
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/market/tickers", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
