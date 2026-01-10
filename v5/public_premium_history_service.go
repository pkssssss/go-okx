package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// PremiumHistory 表示溢价指数历史数据。
//
// 说明：premium 保持为 string（无损）。
type PremiumHistory struct {
	InstId  string `json:"instId"`
	Premium string `json:"premium"`
	TS      int64  `json:"ts,string"`
}

// PublicPremiumHistoryService 获取溢价指数历史数据（近 6 个月，仅适用于永续）。
type PublicPremiumHistoryService struct {
	c *Client

	instId string
	after  string
	before string
	limit  *int
}

// NewPublicPremiumHistoryService 创建 PublicPremiumHistoryService。
func (c *Client) NewPublicPremiumHistoryService() *PublicPremiumHistoryService {
	return &PublicPremiumHistoryService{c: c}
}

// InstId 设置产品 ID（必填；仅适用于永续），如 BTC-USDT-SWAP。
func (s *PublicPremiumHistoryService) InstId(instId string) *PublicPremiumHistoryService {
	s.instId = instId
	return s
}

// After 设置请求此时间戳（不包含）之前的分页内容（传对应接口的 ts）。
func (s *PublicPremiumHistoryService) After(after string) *PublicPremiumHistoryService {
	s.after = after
	return s
}

// Before 设置请求此时间戳（不包含）之后的分页内容（传对应接口的 ts）。
func (s *PublicPremiumHistoryService) Before(before string) *PublicPremiumHistoryService {
	s.before = before
	return s
}

// Limit 设置返回条数（最大 100，默认 100）。
func (s *PublicPremiumHistoryService) Limit(limit int) *PublicPremiumHistoryService {
	s.limit = &limit
	return s
}

var errPublicPremiumHistoryMissingInstId = errors.New("okx: public premium history requires instId")

// Do 获取溢价指数历史数据（GET /api/v5/public/premium-history）。
func (s *PublicPremiumHistoryService) Do(ctx context.Context) ([]PremiumHistory, error) {
	if s.instId == "" {
		return nil, errPublicPremiumHistoryMissingInstId
	}

	q := url.Values{}
	q.Set("instId", s.instId)
	if s.after != "" {
		q.Set("after", s.after)
	}
	if s.before != "" {
		q.Set("before", s.before)
	}
	if s.limit != nil {
		q.Set("limit", strconv.Itoa(*s.limit))
	}

	var data []PremiumHistory
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/public/premium-history", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
