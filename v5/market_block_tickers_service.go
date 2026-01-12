package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// MarketBlockTickersService 获取大宗交易所有产品行情信息（最近24小时成交量信息）。
type MarketBlockTickersService struct {
	c          *Client
	instType   string
	instFamily string
}

// NewMarketBlockTickersService 创建 MarketBlockTickersService。
func (c *Client) NewMarketBlockTickersService() *MarketBlockTickersService {
	return &MarketBlockTickersService{c: c}
}

// InstType 设置产品类型（SPOT/SWAP/FUTURES/OPTION），必填。
func (s *MarketBlockTickersService) InstType(instType string) *MarketBlockTickersService {
	s.instType = instType
	return s
}

// InstFamily 设置交易品种（适用于交割/永续/期权），如 BTC-USD。
func (s *MarketBlockTickersService) InstFamily(instFamily string) *MarketBlockTickersService {
	s.instFamily = instFamily
	return s
}

var errMarketBlockTickersMissingInstType = errors.New("okx: market block tickers requires instType")

// Do 获取大宗交易所有产品行情信息（GET /api/v5/market/block-tickers）。
func (s *MarketBlockTickersService) Do(ctx context.Context) ([]MarketBlockTicker, error) {
	if s.instType == "" {
		return nil, errMarketBlockTickersMissingInstType
	}

	q := url.Values{}
	q.Set("instType", s.instType)
	if s.instFamily != "" {
		q.Set("instFamily", s.instFamily)
	}

	var data []MarketBlockTicker
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/market/block-tickers", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
