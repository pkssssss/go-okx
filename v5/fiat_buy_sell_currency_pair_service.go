package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// FiatBuySellCurrencyPairService 获取买卖交易币对。
type FiatBuySellCurrencyPairService struct {
	c *Client

	fromCcy string
	toCcy   string
}

// NewFiatBuySellCurrencyPairService 创建 FiatBuySellCurrencyPairService。
func (c *Client) NewFiatBuySellCurrencyPairService() *FiatBuySellCurrencyPairService {
	return &FiatBuySellCurrencyPairService{c: c}
}

// FromCcy 设置卖出币种（必填）。
func (s *FiatBuySellCurrencyPairService) FromCcy(fromCcy string) *FiatBuySellCurrencyPairService {
	s.fromCcy = fromCcy
	return s
}

// ToCcy 设置买入币种（必填）。
func (s *FiatBuySellCurrencyPairService) ToCcy(toCcy string) *FiatBuySellCurrencyPairService {
	s.toCcy = toCcy
	return s
}

var errFiatBuySellCurrencyPairMissingRequired = errors.New("okx: fiat buy-sell currency pair requires fromCcy and toCcy")

// Do 获取买卖交易币对（GET /api/v5/fiat/buy-sell/currency-pair）。
func (s *FiatBuySellCurrencyPairService) Do(ctx context.Context) ([]FiatBuySellCurrencyPair, error) {
	if s.fromCcy == "" || s.toCcy == "" {
		return nil, errFiatBuySellCurrencyPairMissingRequired
	}

	q := url.Values{}
	q.Set("fromCcy", s.fromCcy)
	q.Set("toCcy", s.toCcy)

	var data []FiatBuySellCurrencyPair
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/fiat/buy-sell/currency-pair", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
