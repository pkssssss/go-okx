package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// AssetConvertCurrencyPair 表示闪兑币对信息。
// 数值字段保持为 string（无损）。
type AssetConvertCurrencyPair struct {
	InstId      string `json:"instId"`
	BaseCcy     string `json:"baseCcy"`
	BaseCcyMax  string `json:"baseCcyMax"`
	BaseCcyMin  string `json:"baseCcyMin"`
	QuoteCcy    string `json:"quoteCcy"`
	QuoteCcyMax string `json:"quoteCcyMax"`
	QuoteCcyMin string `json:"quoteCcyMin"`
}

// AssetConvertCurrencyPairService 获取闪兑币对信息。
type AssetConvertCurrencyPairService struct {
	c       *Client
	fromCcy string
	toCcy   string
}

// NewAssetConvertCurrencyPairService 创建 AssetConvertCurrencyPairService。
func (c *Client) NewAssetConvertCurrencyPairService() *AssetConvertCurrencyPairService {
	return &AssetConvertCurrencyPairService{c: c}
}

// FromCcy 设置消耗币种（必填）。
func (s *AssetConvertCurrencyPairService) FromCcy(fromCcy string) *AssetConvertCurrencyPairService {
	s.fromCcy = fromCcy
	return s
}

// ToCcy 设置获取币种（必填）。
func (s *AssetConvertCurrencyPairService) ToCcy(toCcy string) *AssetConvertCurrencyPairService {
	s.toCcy = toCcy
	return s
}

var errAssetConvertCurrencyPairMissingRequired = errors.New("okx: convert currency pair requires fromCcy/toCcy")

// Do 获取闪兑币对信息（GET /api/v5/asset/convert/currency-pair）。
func (s *AssetConvertCurrencyPairService) Do(ctx context.Context) ([]AssetConvertCurrencyPair, error) {
	if s.fromCcy == "" || s.toCcy == "" {
		return nil, errAssetConvertCurrencyPairMissingRequired
	}

	q := url.Values{}
	q.Set("fromCcy", s.fromCcy)
	q.Set("toCcy", s.toCcy)

	var data []AssetConvertCurrencyPair
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/asset/convert/currency-pair", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
