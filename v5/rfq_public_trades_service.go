package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// RFQPublicTradeLeg 表示大宗交易公共多腿成交的单腿明细。
//
// 说明：价格/数量等字段保持为 string（无损）。
type RFQPublicTradeLeg struct {
	InstId string `json:"instId"`
	Side   string `json:"side"`
	Sz     string `json:"sz"`
	Px     string `json:"px"`

	TradeId string `json:"tradeId"`
}

// RFQPublicTrade 表示大宗交易公共多腿成交数据。
//
// 说明：时间戳字段解析为 int64（整数不丢精度）；价格/数量等字段保持为 string（无损）。
type RFQPublicTrade struct {
	Strategy string `json:"strategy"`
	CTime    int64  `json:"cTime,string"`

	BlockTdId string `json:"blockTdId"`
	GroupId   string `json:"groupId"`

	Legs []RFQPublicTradeLeg `json:"legs"`
}

// RFQPublicTradesService 获取大宗交易公共多腿成交数据。
type RFQPublicTradesService struct {
	c *Client

	beginId string
	endId   string
	limit   *int
}

// NewRFQPublicTradesService 创建 RFQPublicTradesService。
func (c *Client) NewRFQPublicTradesService() *RFQPublicTradesService {
	return &RFQPublicTradesService{c: c}
}

// BeginId 设置请求的起始大宗交易 ID（请求此 ID 之后更新的数据；不包含 beginId）。
func (s *RFQPublicTradesService) BeginId(beginId string) *RFQPublicTradesService {
	s.beginId = beginId
	return s
}

// EndId 设置请求的结束大宗交易 ID（请求此 ID 之前更旧的数据；不包含 endId）。
func (s *RFQPublicTradesService) EndId(endId string) *RFQPublicTradesService {
	s.endId = endId
	return s
}

// Limit 设置返回条数（最大 100，默认 100）。
func (s *RFQPublicTradesService) Limit(limit int) *RFQPublicTradesService {
	s.limit = &limit
	return s
}

// Do 获取大宗交易公共多腿成交数据（GET /api/v5/rfq/public-trades）。
func (s *RFQPublicTradesService) Do(ctx context.Context) ([]RFQPublicTrade, error) {
	q := url.Values{}
	if s.beginId != "" {
		q.Set("beginId", s.beginId)
	}
	if s.endId != "" {
		q.Set("endId", s.endId)
	}
	if s.limit != nil {
		q.Set("limit", strconv.Itoa(*s.limit))
	}

	var data []RFQPublicTrade
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/rfq/public-trades", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
