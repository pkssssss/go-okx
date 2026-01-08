package okx

import (
	"context"
	"net/http"
	"net/url"
)

// AssetBalance 表示资金账户余额信息。
// 数值字段保持为 string（无损）。
type AssetBalance struct {
	Ccy       string `json:"ccy"`
	Bal       string `json:"bal"`
	FrozenBal string `json:"frozenBal"`
	AvailBal  string `json:"availBal"`
}

// AssetBalancesService 获取资金账户余额。
type AssetBalancesService struct {
	c   *Client
	ccy string
}

// NewAssetBalancesService 创建 AssetBalancesService。
func (c *Client) NewAssetBalancesService() *AssetBalancesService {
	return &AssetBalancesService{c: c}
}

// Ccy 设置币种过滤（支持多币种，逗号分隔；最多 20 个）。
func (s *AssetBalancesService) Ccy(ccy string) *AssetBalancesService {
	s.ccy = ccy
	return s
}

// Do 获取资金账户余额（GET /api/v5/asset/balances）。
func (s *AssetBalancesService) Do(ctx context.Context) ([]AssetBalance, error) {
	var q url.Values
	if s.ccy != "" {
		q = url.Values{}
		q.Set("ccy", s.ccy)
	}

	var data []AssetBalance
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/asset/balances", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
