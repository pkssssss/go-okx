package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// AccountInstrument 表示账户可交易产品信息（精简版）。
//
// 说明：数值字段保持为 string（无损），未包含字段会被忽略，后续可按需补齐。
type AccountInstrument struct {
	InstType   string `json:"instType"`
	InstId     string `json:"instId"`
	InstFamily string `json:"instFamily"`
	Uly        string `json:"uly"`
	Category   string `json:"category"`

	BaseCcy   string `json:"baseCcy"`
	QuoteCcy  string `json:"quoteCcy"`
	SettleCcy string `json:"settleCcy"`

	TickSz string `json:"tickSz"`
	LotSz  string `json:"lotSz"`
	MinSz  string `json:"minSz"`

	CtVal    string `json:"ctVal"`
	CtMult   string `json:"ctMult"`
	CtValCcy string `json:"ctValCcy"`

	GroupId           string   `json:"groupId"`
	TradeQuoteCcyList []string `json:"tradeQuoteCcyList"`

	State string `json:"state"`
}

// AccountInstrumentsService 获取交易产品基础信息（账户可交易）。
type AccountInstrumentsService struct {
	c *Client

	instType   string
	instFamily string
	instId     string
}

// NewAccountInstrumentsService 创建 AccountInstrumentsService。
func (c *Client) NewAccountInstrumentsService() *AccountInstrumentsService {
	return &AccountInstrumentsService{c: c}
}

// InstType 设置产品类型（必填：SPOT/MARGIN/SWAP/FUTURES/OPTION）。
func (s *AccountInstrumentsService) InstType(instType string) *AccountInstrumentsService {
	s.instType = instType
	return s
}

// InstFamily 设置交易品种（可选；期权必填）。
func (s *AccountInstrumentsService) InstFamily(instFamily string) *AccountInstrumentsService {
	s.instFamily = instFamily
	return s
}

// InstId 设置产品 ID（可选）。
func (s *AccountInstrumentsService) InstId(instId string) *AccountInstrumentsService {
	s.instId = instId
	return s
}

var errAccountInstrumentsMissingInstType = errors.New("okx: account instruments requires instType")

// Do 获取交易产品基础信息（GET /api/v5/account/instruments）。
func (s *AccountInstrumentsService) Do(ctx context.Context) ([]AccountInstrument, error) {
	if s.instType == "" {
		return nil, errAccountInstrumentsMissingInstType
	}

	q := url.Values{}
	q.Set("instType", s.instType)
	if s.instFamily != "" {
		q.Set("instFamily", s.instFamily)
	}
	if s.instId != "" {
		q.Set("instId", s.instId)
	}

	var data []AccountInstrument
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/instruments", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
