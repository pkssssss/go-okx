package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// PositionTier 表示衍生品仓位档位信息。
// 数值字段按 OKX 返回保持为 string（无损）。
type PositionTier struct {
	Uly        string `json:"uly"`
	InstFamily string `json:"instFamily"`
	InstId     string `json:"instId"`

	Tier         string `json:"tier"`
	MinSz        string `json:"minSz"`
	MaxSz        string `json:"maxSz"`
	MMR          string `json:"mmr"`
	IMR          string `json:"imr"`
	MaxLever     string `json:"maxLever"`
	OptMgnFactor string `json:"optMgnFactor"`
	QuoteMaxLoan string `json:"quoteMaxLoan"`
	BaseMaxLoan  string `json:"baseMaxLoan"`
}

// PublicPositionTiersService 获取衍生品仓位档位。
type PublicPositionTiersService struct {
	c *Client

	instType   string
	tdMode     string
	instFamily string
	instId     string
	ccy        string
	tier       string
}

// NewPublicPositionTiersService 创建 PublicPositionTiersService。
func (c *Client) NewPublicPositionTiersService() *PublicPositionTiersService {
	return &PublicPositionTiersService{c: c}
}

// InstType 设置产品类型（必填：MARGIN/SWAP/FUTURES/OPTION）。
func (s *PublicPositionTiersService) InstType(instType string) *PublicPositionTiersService {
	s.instType = instType
	return s
}

// TdMode 设置保证金模式（必填：isolated/cross）。
func (s *PublicPositionTiersService) TdMode(tdMode string) *PublicPositionTiersService {
	s.tdMode = tdMode
	return s
}

// InstFamily 设置交易品种（适用于永续/交割/期权；支持多值：不超过 5 个，逗号分隔）。
func (s *PublicPositionTiersService) InstFamily(instFamily string) *PublicPositionTiersService {
	s.instFamily = instFamily
	return s
}

// Uly 设置标的指数（兼容官方示例/SDK 参数名，等价于 InstFamily）。
func (s *PublicPositionTiersService) Uly(uly string) *PublicPositionTiersService {
	s.instFamily = uly
	return s
}

// InstId 设置币对（仅适用币币杠杆；支持多值：不超过 5 个，逗号分隔）。
func (s *PublicPositionTiersService) InstId(instId string) *PublicPositionTiersService {
	s.instId = instId
	return s
}

// Ccy 设置保证金币种（仅适用杠杆全仓；当 instType=MARGIN 且 instId 为空时，ccy 必填）。
func (s *PublicPositionTiersService) Ccy(ccy string) *PublicPositionTiersService {
	s.ccy = ccy
	return s
}

// Tier 设置查询指定档位（可选）。
func (s *PublicPositionTiersService) Tier(tier string) *PublicPositionTiersService {
	s.tier = tier
	return s
}

var (
	errPublicPositionTiersMissingRequired    = errors.New("okx: public position tiers requires instType and tdMode")
	errPublicPositionTiersMissingInstFamily  = errors.New("okx: public position tiers requires instFamily for derivatives")
	errPublicPositionTiersMissingInstIdOrCcy = errors.New("okx: public position tiers requires instId or ccy for MARGIN")
)

// Do 获取衍生品仓位档位（GET /api/v5/public/position-tiers）。
func (s *PublicPositionTiersService) Do(ctx context.Context) ([]PositionTier, error) {
	if s.instType == "" || s.tdMode == "" {
		return nil, errPublicPositionTiersMissingRequired
	}
	if s.instType == "MARGIN" {
		if s.instId == "" && s.ccy == "" {
			return nil, errPublicPositionTiersMissingInstIdOrCcy
		}
	} else {
		if s.instFamily == "" {
			return nil, errPublicPositionTiersMissingInstFamily
		}
	}

	q := url.Values{}
	q.Set("instType", s.instType)
	q.Set("tdMode", s.tdMode)
	if s.instFamily != "" {
		q.Set("instFamily", s.instFamily)
	}
	if s.instId != "" {
		q.Set("instId", s.instId)
	} else if s.ccy != "" {
		q.Set("ccy", s.ccy)
	}
	if s.tier != "" {
		q.Set("tier", s.tier)
	}

	var data []PositionTier
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/public/position-tiers", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
