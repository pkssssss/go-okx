package okx

import (
	"context"
	"errors"
	"net/http"
)

// MarketExchangeRate 表示法币汇率（2 周平均）。
// 数值字段保持为 string（无损）。
type MarketExchangeRate struct {
	UsdCny string `json:"usdCny"`
}

// MarketExchangeRateService 获取法币汇率。
type MarketExchangeRateService struct {
	c *Client
}

// NewMarketExchangeRateService 创建 MarketExchangeRateService。
func (c *Client) NewMarketExchangeRateService() *MarketExchangeRateService {
	return &MarketExchangeRateService{c: c}
}

var errEmptyMarketExchangeRateResponse = errors.New("okx: empty market exchange rate response")

// Do 获取法币汇率（GET /api/v5/market/exchange-rate）。
func (s *MarketExchangeRateService) Do(ctx context.Context) (*MarketExchangeRate, error) {
	var data []MarketExchangeRate
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/market/exchange-rate", nil, nil, false, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyMarketExchangeRateResponse
	}
	return &data[0], nil
}
