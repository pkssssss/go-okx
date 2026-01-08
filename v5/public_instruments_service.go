package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// Instrument 表示 OKX 产品信息（精简版）。
//
// 说明：数值与枚举字段保持为 string（无损），避免 float 精度问题。
type Instrument struct {
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
	CtValCcy string `json:"ctValCcy"`

	State string `json:"state"`
}

// PublicInstrumentsService 查询产品信息。
type PublicInstrumentsService struct {
	c *Client

	instType   string
	uly        string
	instFamily string
	instId     string
}

// NewPublicInstrumentsService 创建 PublicInstrumentsService。
func (c *Client) NewPublicInstrumentsService() *PublicInstrumentsService {
	return &PublicInstrumentsService{c: c}
}

// InstType 设置产品类型（SPOT/SWAP/FUTURES/OPTION），必填。
func (s *PublicInstrumentsService) InstType(instType string) *PublicInstrumentsService {
	s.instType = instType
	return s
}

// Uly 设置标的指数（适用于交割/永续/期权）。
func (s *PublicInstrumentsService) Uly(uly string) *PublicInstrumentsService {
	s.uly = uly
	return s
}

// InstFamily 设置交易品种（适用于交割/永续/期权），如 BTC-USD。
func (s *PublicInstrumentsService) InstFamily(instFamily string) *PublicInstrumentsService {
	s.instFamily = instFamily
	return s
}

// InstId 设置产品 ID。
func (s *PublicInstrumentsService) InstId(instId string) *PublicInstrumentsService {
	s.instId = instId
	return s
}

var errPublicInstrumentsMissingInstType = errors.New("okx: public instruments requires instType")

// Do 查询产品信息（GET /api/v5/public/instruments）。
func (s *PublicInstrumentsService) Do(ctx context.Context) ([]Instrument, error) {
	if s.instType == "" {
		return nil, errPublicInstrumentsMissingInstType
	}

	q := url.Values{}
	q.Set("instType", s.instType)
	if s.uly != "" {
		q.Set("uly", s.uly)
	}
	if s.instFamily != "" {
		q.Set("instFamily", s.instFamily)
	}
	if s.instId != "" {
		q.Set("instId", s.instId)
	}

	var data []Instrument
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/public/instruments", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
