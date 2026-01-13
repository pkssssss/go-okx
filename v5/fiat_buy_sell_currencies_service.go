package okx

import (
	"context"
	"errors"
	"net/http"
)

// FiatBuySellCurrenciesService 获取买卖交易币种。
type FiatBuySellCurrenciesService struct {
	c *Client
}

// NewFiatBuySellCurrenciesService 创建 FiatBuySellCurrenciesService。
func (c *Client) NewFiatBuySellCurrenciesService() *FiatBuySellCurrenciesService {
	return &FiatBuySellCurrenciesService{c: c}
}

var errEmptyFiatBuySellCurrenciesResponse = errors.New("okx: empty fiat buy-sell currencies response")

// Do 获取买卖交易币种（GET /api/v5/fiat/buy-sell/currencies）。
func (s *FiatBuySellCurrenciesService) Do(ctx context.Context) (*FiatBuySellCurrencies, error) {
	var data []FiatBuySellCurrencies
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/fiat/buy-sell/currencies", nil, nil, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyFiatBuySellCurrenciesResponse
	}
	return &data[0], nil
}
