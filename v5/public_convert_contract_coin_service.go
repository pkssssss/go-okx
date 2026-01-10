package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// ConvertContractCoin 表示张/币转换结果。
type ConvertContractCoin struct {
	Type   string `json:"type"`
	InstId string `json:"instId"`
	Px     string `json:"px"`
	Sz     string `json:"sz"`
	Unit   string `json:"unit"`
}

// PublicConvertContractCoinService 张/币转换。
type PublicConvertContractCoinService struct {
	c *Client

	convertType string
	instId      string
	sz          string
	px          string
	unit        string
	opType      string
}

// NewPublicConvertContractCoinService 创建 PublicConvertContractCoinService。
func (c *Client) NewPublicConvertContractCoinService() *PublicConvertContractCoinService {
	return &PublicConvertContractCoinService{c: c}
}

// Type 设置转换类型（可选：1=币转张，2=张转币；默认 1）。
func (s *PublicConvertContractCoinService) Type(convertType string) *PublicConvertContractCoinService {
	s.convertType = convertType
	return s
}

// InstId 设置产品 ID（必填），仅适用于交割/永续/期权。
func (s *PublicConvertContractCoinService) InstId(instId string) *PublicConvertContractCoinService {
	s.instId = instId
	return s
}

// Sz 设置数量（必填）。
//
// 币转张时：币的数量；张转币时：张的数量。
func (s *PublicConvertContractCoinService) Sz(sz string) *PublicConvertContractCoinService {
	s.sz = sz
	return s
}

// Px 设置委托价格（可选；具体是否必填由 OKX 规则决定）。
func (s *PublicConvertContractCoinService) Px(px string) *PublicConvertContractCoinService {
	s.px = px
	return s
}

// Unit 设置币的单位（可选：coin/usds；默认 coin；仅适用于交割/永续的 U 本位合约）。
func (s *PublicConvertContractCoinService) Unit(unit string) *PublicConvertContractCoinService {
	s.unit = unit
	return s
}

// OpType 设置将要下单的类型（可选：open/close；默认 close；适用于交割/永续）。
func (s *PublicConvertContractCoinService) OpType(opType string) *PublicConvertContractCoinService {
	s.opType = opType
	return s
}

var errPublicConvertContractCoinMissingRequired = errors.New("okx: public convert contract coin requires instId/sz")

// Do 张/币转换（GET /api/v5/public/convert-contract-coin）。
func (s *PublicConvertContractCoinService) Do(ctx context.Context) ([]ConvertContractCoin, error) {
	if s.instId == "" || s.sz == "" {
		return nil, errPublicConvertContractCoinMissingRequired
	}

	q := url.Values{}
	q.Set("instId", s.instId)
	q.Set("sz", s.sz)
	if s.convertType != "" {
		q.Set("type", s.convertType)
	}
	if s.px != "" {
		q.Set("px", s.px)
	}
	if s.unit != "" {
		q.Set("unit", s.unit)
	}
	if s.opType != "" {
		q.Set("opType", s.opType)
	}

	var data []ConvertContractCoin
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/public/convert-contract-coin", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
