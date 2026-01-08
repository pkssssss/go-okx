package okx

import (
	"context"
	"net/http"
)

// AssetConvertCurrency 表示闪兑支持的币种。
type AssetConvertCurrency struct {
	Ccy string `json:"ccy"`
	Min string `json:"min"`
	Max string `json:"max"`
}

// AssetConvertCurrenciesService 获取闪兑币种列表。
type AssetConvertCurrenciesService struct {
	c *Client
}

// NewAssetConvertCurrenciesService 创建 AssetConvertCurrenciesService。
func (c *Client) NewAssetConvertCurrenciesService() *AssetConvertCurrenciesService {
	return &AssetConvertCurrenciesService{c: c}
}

// Do 获取闪兑币种列表（GET /api/v5/asset/convert/currencies）。
func (s *AssetConvertCurrenciesService) Do(ctx context.Context) ([]AssetConvertCurrency, error) {
	var data []AssetConvertCurrency
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/asset/convert/currencies", nil, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
