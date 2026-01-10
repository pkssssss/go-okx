package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// SettlementHistoryDetail 表示交割结算记录的明细项。
type SettlementHistoryDetail struct {
	InstId   string `json:"instId"`
	SettlePx string `json:"settlePx"`
}

// SettlementHistory 表示交割结算记录。
//
// 说明：价格字段保持为 string（无损）。
type SettlementHistory struct {
	Details []SettlementHistoryDetail `json:"details"`
	TS      int64                     `json:"ts,string"`
}

// PublicSettlementHistoryService 获取交割结算记录。
type PublicSettlementHistoryService struct {
	c *Client

	instFamily string
	after      string
	before     string
	limit      *int
}

// NewPublicSettlementHistoryService 创建 PublicSettlementHistoryService。
func (c *Client) NewPublicSettlementHistoryService() *PublicSettlementHistoryService {
	return &PublicSettlementHistoryService{c: c}
}

// InstFamily 设置交易品种（必填）。
func (s *PublicSettlementHistoryService) InstFamily(instFamily string) *PublicSettlementHistoryService {
	s.instFamily = instFamily
	return s
}

// After 设置请求此 ts 之前的分页内容。
func (s *PublicSettlementHistoryService) After(after string) *PublicSettlementHistoryService {
	s.after = after
	return s
}

// Before 设置请求此 ts 之后（更新的数据）的分页内容。
func (s *PublicSettlementHistoryService) Before(before string) *PublicSettlementHistoryService {
	s.before = before
	return s
}

// Limit 设置返回条数（最大 100，默认 100）。
func (s *PublicSettlementHistoryService) Limit(limit int) *PublicSettlementHistoryService {
	s.limit = &limit
	return s
}

var errPublicSettlementHistoryMissingInstFamily = errors.New("okx: public settlement history requires instFamily")

// Do 获取交割结算记录（GET /api/v5/public/settlement-history）。
func (s *PublicSettlementHistoryService) Do(ctx context.Context) ([]SettlementHistory, error) {
	if s.instFamily == "" {
		return nil, errPublicSettlementHistoryMissingInstFamily
	}

	q := url.Values{}
	q.Set("instFamily", s.instFamily)
	if s.after != "" {
		q.Set("after", s.after)
	}
	if s.before != "" {
		q.Set("before", s.before)
	}
	if s.limit != nil {
		q.Set("limit", strconv.Itoa(*s.limit))
	}

	var data []SettlementHistory
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/public/settlement-history", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
