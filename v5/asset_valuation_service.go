package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// AssetValuationDetails 表示各账户的资产估值明细。
type AssetValuationDetails struct {
	Classic string `json:"classic"`
	Earn    string `json:"earn"`
	Funding string `json:"funding"`
	Trading string `json:"trading"`
}

// AssetValuation 表示账户资产估值。
// 数值字段保持为 string（无损）。
type AssetValuation struct {
	TotalBal string                `json:"totalBal"`
	TS       int64                 `json:"ts,string"`
	Details  AssetValuationDetails `json:"details"`
}

// AssetValuationService 获取账户资产估值。
type AssetValuationService struct {
	c   *Client
	ccy string
}

// NewAssetValuationService 创建 AssetValuationService。
func (c *Client) NewAssetValuationService() *AssetValuationService {
	return &AssetValuationService{c: c}
}

// Ccy 设置估值对应的单位（默认 BTC）。
func (s *AssetValuationService) Ccy(ccy string) *AssetValuationService {
	s.ccy = ccy
	return s
}

var errEmptyAssetValuation = errors.New("okx: empty asset valuation response")

// Do 获取账户资产估值（GET /api/v5/asset/asset-valuation）。
func (s *AssetValuationService) Do(ctx context.Context) (*AssetValuation, error) {
	var q url.Values
	if s.ccy != "" {
		q = url.Values{}
		q.Set("ccy", s.ccy)
	}

	var data []AssetValuation
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/asset/asset-valuation", q, nil, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyAssetValuation
	}
	return &data[0], nil
}
