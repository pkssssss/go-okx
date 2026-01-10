package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// InstrumentTickBand 表示期权价格梯度中的单个价格区间。
//
// 说明：数值字段保持为 string（无损），避免 float 精度问题。
type InstrumentTickBand struct {
	MinPx  string `json:"minPx"`
	MaxPx  string `json:"maxPx"`
	TickSz string `json:"tickSz"`
}

// InstrumentTickBandInfo 表示期权价格梯度信息。
type InstrumentTickBandInfo struct {
	InstType   string               `json:"instType"`
	InstFamily string               `json:"instFamily"`
	TickBands  []InstrumentTickBand `json:"tickBand"`
}

// PublicInstrumentTickBandsService 获取期权价格梯度。
type PublicInstrumentTickBandsService struct {
	c *Client

	instType   string
	instFamily string
}

// NewPublicInstrumentTickBandsService 创建 PublicInstrumentTickBandsService。
func (c *Client) NewPublicInstrumentTickBandsService() *PublicInstrumentTickBandsService {
	return &PublicInstrumentTickBandsService{c: c}
}

// InstType 设置产品类型（必填；当前仅支持 OPTION）。
func (s *PublicInstrumentTickBandsService) InstType(instType string) *PublicInstrumentTickBandsService {
	s.instType = instType
	return s
}

// InstFamily 设置交易品种（可选，仅适用于期权），如 BTC-USD。
func (s *PublicInstrumentTickBandsService) InstFamily(instFamily string) *PublicInstrumentTickBandsService {
	s.instFamily = instFamily
	return s
}

var errPublicInstrumentTickBandsMissingInstType = errors.New("okx: public instrument tick bands requires instType")

// Do 获取期权价格梯度（GET /api/v5/public/instrument-tick-bands）。
func (s *PublicInstrumentTickBandsService) Do(ctx context.Context) ([]InstrumentTickBandInfo, error) {
	if s.instType == "" {
		return nil, errPublicInstrumentTickBandsMissingInstType
	}

	q := url.Values{}
	q.Set("instType", s.instType)
	if s.instFamily != "" {
		q.Set("instFamily", s.instFamily)
	}

	var data []InstrumentTickBandInfo
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/public/instrument-tick-bands", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
