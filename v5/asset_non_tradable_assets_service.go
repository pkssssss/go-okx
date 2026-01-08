package okx

import (
	"context"
	"net/http"
	"net/url"
)

// AssetNonTradableAsset 表示不可交易资产（含可提信息）。
// 数值字段保持为 string（无损）。
type AssetNonTradableAsset struct {
	Ccy      string `json:"ccy"`
	Name     string `json:"name"`
	LogoLink string `json:"logoLink"`
	Chain    string `json:"chain"`
	CtAddr   string `json:"ctAddr"`

	Bal   string `json:"bal"`
	CanWd bool   `json:"canWd"`

	MinWd          string `json:"minWd"`
	Fee            string `json:"fee"`
	FeeCcy         string `json:"feeCcy"`
	WdTickSz       string `json:"wdTickSz"`
	WdAll          bool   `json:"wdAll"`
	NeedTag        bool   `json:"needTag"`
	MainNet        bool   `json:"mainNet"`
	BurningFeeRate string `json:"burningFeeRate"`
}

// AssetNonTradableAssetsService 获取不可交易资产列表。
type AssetNonTradableAssetsService struct {
	c   *Client
	ccy string
}

// NewAssetNonTradableAssetsService 创建 AssetNonTradableAssetsService。
func (c *Client) NewAssetNonTradableAssetsService() *AssetNonTradableAssetsService {
	return &AssetNonTradableAssetsService{c: c}
}

// Ccy 设置币种过滤（支持多币种查询，逗号分隔，不超过 20）。
func (s *AssetNonTradableAssetsService) Ccy(ccy string) *AssetNonTradableAssetsService {
	s.ccy = ccy
	return s
}

// Do 获取不可交易资产列表（GET /api/v5/asset/non-tradable-assets）。
func (s *AssetNonTradableAssetsService) Do(ctx context.Context) ([]AssetNonTradableAsset, error) {
	var q url.Values
	if s.ccy != "" {
		q = url.Values{}
		q.Set("ccy", s.ccy)
	}

	var data []AssetNonTradableAsset
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/asset/non-tradable-assets", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
