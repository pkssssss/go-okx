package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// AccountCollateralAsset 表示质押币种信息。
type AccountCollateralAsset struct {
	Ccy               string `json:"ccy"`
	CollateralEnabled bool   `json:"collateralEnabled"`
}

// AccountCollateralAssetsService 查看质押币种。
type AccountCollateralAssetsService struct {
	c *Client

	ccy               string
	collateralEnabled *bool
}

// NewAccountCollateralAssetsService 创建 AccountCollateralAssetsService。
func (c *Client) NewAccountCollateralAssetsService() *AccountCollateralAssetsService {
	return &AccountCollateralAssetsService{c: c}
}

// Ccy 设置币种过滤（支持多币种，逗号分隔，如 "BTC,ETH"；最多 20 个）。
func (s *AccountCollateralAssetsService) Ccy(ccy string) *AccountCollateralAssetsService {
	s.ccy = ccy
	return s
}

// CollateralEnabled 设置是否为质押币过滤。
func (s *AccountCollateralAssetsService) CollateralEnabled(enable bool) *AccountCollateralAssetsService {
	s.collateralEnabled = &enable
	return s
}

// Do 查看质押币种（GET /api/v5/account/collateral-assets）。
func (s *AccountCollateralAssetsService) Do(ctx context.Context) ([]AccountCollateralAsset, error) {
	var q url.Values
	if s.ccy != "" || s.collateralEnabled != nil {
		q = url.Values{}
		if s.ccy != "" {
			q.Set("ccy", s.ccy)
		}
		if s.collateralEnabled != nil {
			q.Set("collateralEnabled", strconv.FormatBool(*s.collateralEnabled))
		}
	}

	var data []AccountCollateralAsset
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/collateral-assets", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
