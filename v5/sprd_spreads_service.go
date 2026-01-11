package okx

import (
	"context"
	"net/http"
	"net/url"
)

// SprdLeg 表示 Spread 的腿信息。
type SprdLeg struct {
	InstId string `json:"instId"`
	Side   string `json:"side"`
}

// SprdSpread 表示可交易的 Spread 产品信息。
//
// 说明：价格/数量等字段保持为 string（无损）；时间戳使用 int64（json string）。
type SprdSpread struct {
	SprdId   string `json:"sprdId"`
	SprdType string `json:"sprdType"`
	State    string `json:"state"`

	BaseCcy  string `json:"baseCcy"`
	SzCcy    string `json:"szCcy"`
	QuoteCcy string `json:"quoteCcy"`

	TickSz string `json:"tickSz"`
	MinSz  string `json:"minSz"`
	LotSz  string `json:"lotSz"`

	ListTime int64 `json:"listTime,string"`
	ExpTime  int64 `json:"expTime,string"`
	UTime    int64 `json:"uTime,string"`

	Legs []SprdLeg `json:"legs"`
}

// SprdSpreadsService 获取可交易的 Spreads（公共）。
type SprdSpreadsService struct {
	c *Client

	baseCcy string
	instId  string
	sprdId  string
	state   string
}

// NewSprdSpreadsService 创建 SprdSpreadsService。
func (c *Client) NewSprdSpreadsService() *SprdSpreadsService {
	return &SprdSpreadsService{c: c}
}

// BaseCcy 设置 Spread 币种（可选），如 BTC。
func (s *SprdSpreadsService) BaseCcy(baseCcy string) *SprdSpreadsService {
	s.baseCcy = baseCcy
	return s
}

// InstId 设置 Spread 里包含的产品 ID（可选），如 BTC-USDT。
func (s *SprdSpreadsService) InstId(instId string) *SprdSpreadsService {
	s.instId = instId
	return s
}

// SprdId 设置 Spread ID（可选）。
func (s *SprdSpreadsService) SprdId(sprdId string) *SprdSpreadsService {
	s.sprdId = sprdId
	return s
}

// State 设置 Spread 状态（可选）：live/suspend/expired。
func (s *SprdSpreadsService) State(state string) *SprdSpreadsService {
	s.state = state
	return s
}

// Do 获取可交易的 Spreads（GET /api/v5/sprd/spreads）。
func (s *SprdSpreadsService) Do(ctx context.Context) ([]SprdSpread, error) {
	q := url.Values{}
	if s.baseCcy != "" {
		q.Set("baseCcy", s.baseCcy)
	}
	if s.instId != "" {
		q.Set("instId", s.instId)
	}
	if s.sprdId != "" {
		q.Set("sprdId", s.sprdId)
	}
	if s.state != "" {
		q.Set("state", s.state)
	}

	var data []SprdSpread
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/sprd/spreads", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
