package okx

import (
	"context"
	"net/http"
	"net/url"
)

// AssetCurrency 表示资金账户支持的币种/链信息（含充提与手续费/精度信息）。
// 数值字段保持为 string（无损）。
type AssetCurrency struct {
	Ccy      string `json:"ccy"`
	Name     string `json:"name"`
	LogoLink string `json:"logoLink"`
	Chain    string `json:"chain"`
	CtAddr   string `json:"ctAddr"`

	CanDep      bool `json:"canDep"`
	CanWd       bool `json:"canWd"`
	CanInternal bool `json:"canInternal"`

	DepEstOpenTime string `json:"depEstOpenTime"`
	WdEstOpenTime  string `json:"wdEstOpenTime"`

	MinDep      string `json:"minDep"`
	MinWd       string `json:"minWd"`
	MinInternal string `json:"minInternal"`
	MaxWd       string `json:"maxWd"`

	WdTickSz    string `json:"wdTickSz"`
	WdQuota     string `json:"wdQuota"`
	UsedWdQuota string `json:"usedWdQuota"`

	Fee            string `json:"fee"`
	BurningFeeRate string `json:"burningFeeRate"`

	MainNet bool `json:"mainNet"`
	NeedTag bool `json:"needTag"`

	MinDepArrivalConfirm string `json:"minDepArrivalConfirm"`
	MinWdUnlockConfirm   string `json:"minWdUnlockConfirm"`

	DepQuotaFixed       string `json:"depQuotaFixed"`
	UsedDepQuotaFixed   string `json:"usedDepQuotaFixed"`
	DepQuoteDailyLayer2 string `json:"depQuoteDailyLayer2"`
}

// AssetCurrenciesService 获取币种列表（资金账户，按用户 KYC 实体返回）。
type AssetCurrenciesService struct {
	c   *Client
	ccy string
}

// NewAssetCurrenciesService 创建 AssetCurrenciesService。
func (c *Client) NewAssetCurrenciesService() *AssetCurrenciesService {
	return &AssetCurrenciesService{c: c}
}

// Ccy 设置币种过滤（支持多币种，逗号分隔）。
func (s *AssetCurrenciesService) Ccy(ccy string) *AssetCurrenciesService {
	s.ccy = ccy
	return s
}

// Do 获取币种列表（GET /api/v5/asset/currencies）。
func (s *AssetCurrenciesService) Do(ctx context.Context) ([]AssetCurrency, error) {
	var q url.Values
	if s.ccy != "" {
		q = url.Values{}
		q.Set("ccy", s.ccy)
	}

	var data []AssetCurrency
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/asset/currencies", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
