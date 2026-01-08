package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// AssetDepositAddress 表示充值地址信息。
type AssetDepositAddress struct {
	Addr   string            `json:"addr"`
	Tag    string            `json:"tag"`
	Memo   string            `json:"memo"`
	PmtId  string            `json:"pmtId"`
	AddrEx map[string]string `json:"addrEx"`

	Ccy   string `json:"ccy"`
	Chain string `json:"chain"`
	To    string `json:"to"`

	VerifiedName string `json:"verifiedName"`
	Selected     bool   `json:"selected"`
	CtAddr       string `json:"ctAddr"`
}

// AssetDepositAddressService 获取充值地址信息（包含曾使用过的老地址）。
type AssetDepositAddressService struct {
	c   *Client
	ccy string
}

// NewAssetDepositAddressService 创建 AssetDepositAddressService。
func (c *Client) NewAssetDepositAddressService() *AssetDepositAddressService {
	return &AssetDepositAddressService{c: c}
}

// Ccy 设置币种（必填）。
func (s *AssetDepositAddressService) Ccy(ccy string) *AssetDepositAddressService {
	s.ccy = ccy
	return s
}

var errAssetDepositAddressMissingCcy = errors.New("okx: deposit address requires ccy")

// Do 获取充值地址信息（GET /api/v5/asset/deposit-address）。
func (s *AssetDepositAddressService) Do(ctx context.Context) ([]AssetDepositAddress, error) {
	if s.ccy == "" {
		return nil, errAssetDepositAddressMissingCcy
	}

	q := url.Values{}
	q.Set("ccy", s.ccy)

	var data []AssetDepositAddress
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/asset/deposit-address", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
