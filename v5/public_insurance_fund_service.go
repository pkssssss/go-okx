package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// InsuranceFundDetail 表示风险保证金详情。
//
// 说明：数值字段保持为 string（无损），避免 float 精度问题。
type InsuranceFundDetail struct {
	Balance  string `json:"balance"`
	Amt      string `json:"amt"`
	Ccy      string `json:"ccy"`
	Type     string `json:"type"`
	MaxBal   string `json:"maxBal"`
	MaxBalTs string `json:"maxBalTs"`
	DecRate  string `json:"decRate"` // 已弃用
	AdlType  string `json:"adlType"`
	TS       int64  `json:"ts,string"`
}

// InsuranceFund 表示风险保证金余额信息。
//
// 说明：total 单位为 USD；数值字段保持为 string（无损）。
type InsuranceFund struct {
	Total      string                `json:"total"`
	InstFamily string                `json:"instFamily"`
	InstType   string                `json:"instType"`
	Details    []InsuranceFundDetail `json:"details"`
}

// PublicInsuranceFundService 获取风险保证金余额。
type PublicInsuranceFundService struct {
	c *Client

	instType   string
	fundType   string
	instFamily string
	uly        string
	ccy        string
	before     string
	after      string
	limit      *int
}

// NewPublicInsuranceFundService 创建 PublicInsuranceFundService。
func (c *Client) NewPublicInsuranceFundService() *PublicInsuranceFundService {
	return &PublicInsuranceFundService{c: c}
}

// InstType 设置产品类型（必填：MARGIN/SWAP/FUTURES/OPTION）。
func (s *PublicInsuranceFundService) InstType(instType string) *PublicInsuranceFundService {
	s.instType = instType
	return s
}

// Type 设置风险保证金类型（可选：regular_update/liquidation_balance_deposit/bankruptcy_loss/platform_revenue/adl）。
func (s *PublicInsuranceFundService) Type(fundType string) *PublicInsuranceFundService {
	s.fundType = fundType
	return s
}

// InstFamily 设置交易品种（可选；当 instType 为 SWAP/FUTURES/OPTION 时，instFamily/uly 至少传一个）。
func (s *PublicInsuranceFundService) InstFamily(instFamily string) *PublicInsuranceFundService {
	s.instFamily = instFamily
	return s
}

// Uly 设置标的指数（可选；当 instType 为 SWAP/FUTURES/OPTION 时，instFamily/uly 至少传一个）。
func (s *PublicInsuranceFundService) Uly(uly string) *PublicInsuranceFundService {
	s.uly = uly
	return s
}

// Ccy 设置币种（仅适用 instType=MARGIN，且必填）。
func (s *PublicInsuranceFundService) Ccy(ccy string) *PublicInsuranceFundService {
	s.ccy = ccy
	return s
}

func (s *PublicInsuranceFundService) Before(before string) *PublicInsuranceFundService {
	s.before = before
	return s
}

func (s *PublicInsuranceFundService) After(after string) *PublicInsuranceFundService {
	s.after = after
	return s
}

func (s *PublicInsuranceFundService) Limit(limit int) *PublicInsuranceFundService {
	s.limit = &limit
	return s
}

var (
	errPublicInsuranceFundMissingInstType        = errors.New("okx: public insurance fund requires instType")
	errPublicInsuranceFundMissingCcy             = errors.New("okx: public insurance fund requires ccy for MARGIN")
	errPublicInsuranceFundMissingInstFamilyOrUly = errors.New("okx: public insurance fund requires instFamily or uly")
)

// Do 获取风险保证金余额（GET /api/v5/public/insurance-fund）。
func (s *PublicInsuranceFundService) Do(ctx context.Context) ([]InsuranceFund, error) {
	if s.instType == "" {
		return nil, errPublicInsuranceFundMissingInstType
	}
	if s.instType == "MARGIN" {
		if s.ccy == "" {
			return nil, errPublicInsuranceFundMissingCcy
		}
	} else {
		if s.instFamily == "" && s.uly == "" {
			return nil, errPublicInsuranceFundMissingInstFamilyOrUly
		}
	}

	q := url.Values{}
	q.Set("instType", s.instType)
	if s.fundType != "" {
		q.Set("type", s.fundType)
	}
	if s.instFamily != "" {
		q.Set("instFamily", s.instFamily)
	}
	if s.uly != "" {
		q.Set("uly", s.uly)
	}
	if s.ccy != "" {
		q.Set("ccy", s.ccy)
	}
	if s.before != "" {
		q.Set("before", s.before)
	}
	if s.after != "" {
		q.Set("after", s.after)
	}
	if s.limit != nil {
		q.Set("limit", strconv.Itoa(*s.limit))
	}

	var data []InsuranceFund
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/public/insurance-fund", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
