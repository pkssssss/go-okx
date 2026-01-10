package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// OptionTrade 表示期权公共成交数据。
//
// 说明：价格/数量等字段保持为 string（无损）。
type OptionTrade struct {
	InstFamily string `json:"instFamily"`
	InstId     string `json:"instId"`

	TradeId string `json:"tradeId"`
	Px      string `json:"px"`
	Sz      string `json:"sz"`
	Side    string `json:"side"`

	OptType string `json:"optType"`
	FillVol string `json:"fillVol"`
	FwdPx   string `json:"fwdPx"`
	IdxPx   string `json:"idxPx"`
	MarkPx  string `json:"markPx"`

	TS int64 `json:"ts,string"`
}

// PublicOptionTradesService 获取期权公共成交数据。
type PublicOptionTradesService struct {
	c *Client

	instId     string
	instFamily string
	optType    string
}

// NewPublicOptionTradesService 创建 PublicOptionTradesService。
func (c *Client) NewPublicOptionTradesService() *PublicOptionTradesService {
	return &PublicOptionTradesService{c: c}
}

// InstId 设置产品 ID（可选；instId/instFamily 必须传一个；若都传，以 instId 为主）。
func (s *PublicOptionTradesService) InstId(instId string) *PublicOptionTradesService {
	s.instId = instId
	return s
}

// InstFamily 设置交易品种（可选；instId/instFamily 必须传一个）。
func (s *PublicOptionTradesService) InstFamily(instFamily string) *PublicOptionTradesService {
	s.instFamily = instFamily
	return s
}

// OptType 设置期权类型（可选：C/P）。
func (s *PublicOptionTradesService) OptType(optType string) *PublicOptionTradesService {
	s.optType = optType
	return s
}

var errPublicOptionTradesMissingInstIdOrInstFamily = errors.New("okx: public option trades requires instId or instFamily")

// Do 获取期权公共成交数据（GET /api/v5/public/option-trades）。
func (s *PublicOptionTradesService) Do(ctx context.Context) ([]OptionTrade, error) {
	if s.instId == "" && s.instFamily == "" {
		return nil, errPublicOptionTradesMissingInstIdOrInstFamily
	}

	q := url.Values{}
	if s.instId != "" {
		q.Set("instId", s.instId)
	} else {
		q.Set("instFamily", s.instFamily)
	}
	if s.optType != "" {
		q.Set("optType", s.optType)
	}

	var data []OptionTrade
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/public/option-trades", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
