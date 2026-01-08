package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// AssetSubaccountBalancesService 获取子账户资金账户余额（母账户）。
type AssetSubaccountBalancesService struct {
	c       *Client
	subAcct string
	ccy     string
}

// NewAssetSubaccountBalancesService 创建 AssetSubaccountBalancesService。
func (c *Client) NewAssetSubaccountBalancesService() *AssetSubaccountBalancesService {
	return &AssetSubaccountBalancesService{c: c}
}

// SubAcct 设置子账户名称（必填）。
func (s *AssetSubaccountBalancesService) SubAcct(subAcct string) *AssetSubaccountBalancesService {
	s.subAcct = subAcct
	return s
}

// Ccy 设置币种过滤（支持多币种，逗号分隔；最多 20 个）。
func (s *AssetSubaccountBalancesService) Ccy(ccy string) *AssetSubaccountBalancesService {
	s.ccy = ccy
	return s
}

var errAssetSubaccountBalancesMissingSubAcct = errors.New("okx: subaccount balances requires subAcct")

// Do 获取子账户资金账户余额（GET /api/v5/asset/subaccount/balances）。
func (s *AssetSubaccountBalancesService) Do(ctx context.Context) ([]AssetBalance, error) {
	if s.subAcct == "" {
		return nil, errAssetSubaccountBalancesMissingSubAcct
	}

	q := url.Values{}
	q.Set("subAcct", s.subAcct)
	if s.ccy != "" {
		q.Set("ccy", s.ccy)
	}

	var data []AssetBalance
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/asset/subaccount/balances", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
