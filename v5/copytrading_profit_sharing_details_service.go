package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type copyTradingProfitSharingDetailsQuery struct {
	instType string
	after    string
	before   string
	limit    *int
}

func (q copyTradingProfitSharingDetailsQuery) values() url.Values {
	v := url.Values{}
	if q.instType != "" {
		v.Set("instType", q.instType)
	}
	if q.after != "" {
		v.Set("after", q.after)
	}
	if q.before != "" {
		v.Set("before", q.before)
	}
	if q.limit != nil {
		v.Set("limit", strconv.Itoa(*q.limit))
	}

	if len(v) == 0 {
		return nil
	}
	return v
}

// CopyTradingProfitSharingDetailsService 交易员历史分润明细。
type CopyTradingProfitSharingDetailsService struct {
	c *Client
	q copyTradingProfitSharingDetailsQuery
}

// NewCopyTradingProfitSharingDetailsService 创建 CopyTradingProfitSharingDetailsService。
func (c *Client) NewCopyTradingProfitSharingDetailsService() *CopyTradingProfitSharingDetailsService {
	return &CopyTradingProfitSharingDetailsService{c: c}
}

// InstType 设置产品类型（默认返回所有）。
func (s *CopyTradingProfitSharingDetailsService) InstType(instType string) *CopyTradingProfitSharingDetailsService {
	s.q.instType = instType
	return s
}

// After 请求此 id 之前（更旧数据）的分页内容（profitSharingId）。
func (s *CopyTradingProfitSharingDetailsService) After(after string) *CopyTradingProfitSharingDetailsService {
	s.q.after = after
	return s
}

// Before 请求此 id 之后（更新数据）的分页内容（profitSharingId）。
func (s *CopyTradingProfitSharingDetailsService) Before(before string) *CopyTradingProfitSharingDetailsService {
	s.q.before = before
	return s
}

// Limit 分页返回数量（最大 100，默认 100）。
func (s *CopyTradingProfitSharingDetailsService) Limit(limit int) *CopyTradingProfitSharingDetailsService {
	s.q.limit = &limit
	return s
}

// Do 获取历史分润明细（GET /api/v5/copytrading/profit-sharing-details）。
func (s *CopyTradingProfitSharingDetailsService) Do(ctx context.Context) ([]CopyTradingProfitSharingDetail, error) {
	var data []CopyTradingProfitSharingDetail
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/copytrading/profit-sharing-details", s.q.values(), nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
