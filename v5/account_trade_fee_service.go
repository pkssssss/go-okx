package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// AccountTradeFeeGroup 表示交易手续费分组。
type AccountTradeFeeGroup struct {
	GroupId string `json:"groupId"`
	Maker   string `json:"maker"`
	Taker   string `json:"taker"`
}

// AccountTradeFeeFiat 表示法币费率（已废弃字段仍可能返回）。
type AccountTradeFeeFiat struct {
	Ccy   string `json:"ccy"`
	Maker string `json:"maker"`
	Taker string `json:"taker"`
}

// AccountTradeFee 表示当前账户交易手续费费率。
type AccountTradeFee struct {
	Level    string `json:"level"`
	InstType string `json:"instType"`

	FeeGroup []AccountTradeFeeGroup `json:"feeGroup"`

	Delivery string `json:"delivery"`
	Exercise string `json:"exercise"`

	Maker     string `json:"maker"`
	Taker     string `json:"taker"`
	MakerU    string `json:"makerU"`
	TakerU    string `json:"takerU"`
	MakerUSDC string `json:"makerUSDC"`
	TakerUSDC string `json:"takerUSDC"`

	RuleType string                `json:"ruleType"`
	Category string                `json:"category"`
	Fiat     []AccountTradeFeeFiat `json:"fiat"`

	TS UnixMilli `json:"ts"`
}

// AccountTradeFeeService 获取当前账户交易手续费费率。
type AccountTradeFeeService struct {
	c *Client

	instType   string
	instId     string
	instFamily string
	groupId    string
}

// NewAccountTradeFeeService 创建 AccountTradeFeeService。
func (c *Client) NewAccountTradeFeeService() *AccountTradeFeeService {
	return &AccountTradeFeeService{c: c}
}

// InstType 设置产品类型（必填：SPOT/MARGIN/SWAP/FUTURES/OPTION）。
func (s *AccountTradeFeeService) InstType(instType string) *AccountTradeFeeService {
	s.instType = instType
	return s
}

// InstId 设置产品 ID（仅适用于 instType 为币币/币币杠杆）。
func (s *AccountTradeFeeService) InstId(instId string) *AccountTradeFeeService {
	s.instId = instId
	return s
}

// InstFamily 设置交易品种（适用于交割/永续/期权，如 BTC-USD）。
func (s *AccountTradeFeeService) InstFamily(instFamily string) *AccountTradeFeeService {
	s.instFamily = instFamily
	return s
}

// GroupId 设置交易产品手续费分组 ID（与 instId/instFamily 互斥）。
func (s *AccountTradeFeeService) GroupId(groupId string) *AccountTradeFeeService {
	s.groupId = groupId
	return s
}

var (
	errAccountTradeFeeMissingInstType = errors.New("okx: trade fee requires instType")
	errAccountTradeFeeGroupIdConflict = errors.New("okx: trade fee groupId is mutually exclusive with instId/instFamily")
)

// Do 获取当前账户交易手续费费率（GET /api/v5/account/trade-fee）。
func (s *AccountTradeFeeService) Do(ctx context.Context) ([]AccountTradeFee, error) {
	if s.instType == "" {
		return nil, errAccountTradeFeeMissingInstType
	}
	if s.groupId != "" && (s.instId != "" || s.instFamily != "") {
		return nil, errAccountTradeFeeGroupIdConflict
	}

	q := url.Values{}
	q.Set("instType", s.instType)
	if s.instId != "" {
		q.Set("instId", s.instId)
	}
	if s.instFamily != "" {
		q.Set("instFamily", s.instFamily)
	}
	if s.groupId != "" {
		q.Set("groupId", s.groupId)
	}

	var data []AccountTradeFee
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/trade-fee", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
