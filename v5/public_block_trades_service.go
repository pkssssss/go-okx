package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// BlockTrade 表示大宗交易公共单腿成交数据。
//
// 说明：价格/数量等字段保持为 string（无损）。
type BlockTrade struct {
	InstId  string `json:"instId"`
	TradeId string `json:"tradeId"`

	Px   string `json:"px"`
	Sz   string `json:"sz"`
	Side string `json:"side"`

	FillVol string `json:"fillVol"`
	FwdPx   string `json:"fwdPx"`
	IdxPx   string `json:"idxPx"`
	MarkPx  string `json:"markPx"`

	GroupId string `json:"groupId"`
	TS      int64  `json:"ts,string"`
}

// PublicBlockTradesService 获取大宗交易公共单腿成交数据。
type PublicBlockTradesService struct {
	c *Client

	instId string
}

// NewPublicBlockTradesService 创建 PublicBlockTradesService。
func (c *Client) NewPublicBlockTradesService() *PublicBlockTradesService {
	return &PublicBlockTradesService{c: c}
}

// InstId 设置产品 ID（必填）。
func (s *PublicBlockTradesService) InstId(instId string) *PublicBlockTradesService {
	s.instId = instId
	return s
}

var errPublicBlockTradesMissingInstId = errors.New("okx: public block trades requires instId")

// Do 获取大宗交易公共单腿成交数据（GET /api/v5/public/block-trades）。
func (s *PublicBlockTradesService) Do(ctx context.Context) ([]BlockTrade, error) {
	if s.instId == "" {
		return nil, errPublicBlockTradesMissingInstId
	}

	q := url.Values{}
	q.Set("instId", s.instId)

	var data []BlockTrade
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/public/block-trades", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
