package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// AccountPositionTier 表示组合保证金模式仓位限制（档位）。
// 数值字段按 OKX 返回保持为 string（无损）。
type AccountPositionTier struct {
	Uly        string `json:"uly"`
	InstFamily string `json:"instFamily"`
	MaxSz      string `json:"maxSz"`
	PosType    string `json:"posType"`
}

// AccountPositionTiersService 获取组合保证金模式仓位限制。
type AccountPositionTiersService struct {
	c *Client

	instType   string
	instFamily string
}

// NewAccountPositionTiersService 创建 AccountPositionTiersService。
func (c *Client) NewAccountPositionTiersService() *AccountPositionTiersService {
	return &AccountPositionTiersService{c: c}
}

// InstType 设置产品类型（SWAP/FUTURES/OPTION），必填。
func (s *AccountPositionTiersService) InstType(instType string) *AccountPositionTiersService {
	s.instType = instType
	return s
}

// InstFamily 设置交易品种（必填，支持多值：不超过 5 个，逗号分隔）。
func (s *AccountPositionTiersService) InstFamily(instFamily string) *AccountPositionTiersService {
	s.instFamily = instFamily
	return s
}

// Uly 设置标的指数（兼容官方示例/SDK 参数名，等价于 InstFamily）。
func (s *AccountPositionTiersService) Uly(uly string) *AccountPositionTiersService {
	s.instFamily = uly
	return s
}

var errAccountPositionTiersMissingRequired = errors.New("okx: account position tiers requires instType and instFamily")

// Do 获取组合保证金模式仓位限制（GET /api/v5/account/position-tiers）。
func (s *AccountPositionTiersService) Do(ctx context.Context) ([]AccountPositionTier, error) {
	if s.instType == "" || s.instFamily == "" {
		return nil, errAccountPositionTiersMissingRequired
	}

	q := url.Values{}
	q.Set("instType", s.instType)
	q.Set("instFamily", s.instFamily)

	var data []AccountPositionTier
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/position-tiers", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
