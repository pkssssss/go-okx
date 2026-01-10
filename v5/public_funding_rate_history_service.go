package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// FundingRateHistory 表示永续合约历史资金费率。
//
// 说明：费率字段保持为 string（无损）。
type FundingRateHistory struct {
	InstType     string `json:"instType"`
	InstId       string `json:"instId"`
	FormulaType  string `json:"formulaType"`
	FundingRate  string `json:"fundingRate"`
	RealizedRate string `json:"realizedRate"`
	Method       string `json:"method"`
	FundingTime  int64  `json:"fundingTime,string"`
}

// PublicFundingRateHistoryService 获取永续合约历史资金费率。
type PublicFundingRateHistoryService struct {
	c *Client

	instId string
	before string
	after  string
	limit  *int
}

// NewPublicFundingRateHistoryService 创建 PublicFundingRateHistoryService。
func (c *Client) NewPublicFundingRateHistoryService() *PublicFundingRateHistoryService {
	return &PublicFundingRateHistoryService{c: c}
}

// InstId 设置产品 ID（必填；仅适用于永续）。
func (s *PublicFundingRateHistoryService) InstId(instId string) *PublicFundingRateHistoryService {
	s.instId = instId
	return s
}

// Before 设置请求此 fundingTime 之后（更新的数据）的分页内容。
func (s *PublicFundingRateHistoryService) Before(before string) *PublicFundingRateHistoryService {
	s.before = before
	return s
}

// After 设置请求此 fundingTime 之前（更旧的数据）的分页内容。
func (s *PublicFundingRateHistoryService) After(after string) *PublicFundingRateHistoryService {
	s.after = after
	return s
}

// Limit 设置返回条数（最大 400，默认 400）。
func (s *PublicFundingRateHistoryService) Limit(limit int) *PublicFundingRateHistoryService {
	s.limit = &limit
	return s
}

var errPublicFundingRateHistoryMissingInstId = errors.New("okx: public funding rate history requires instId")

// Do 获取永续合约历史资金费率（GET /api/v5/public/funding-rate-history）。
func (s *PublicFundingRateHistoryService) Do(ctx context.Context) ([]FundingRateHistory, error) {
	if s.instId == "" {
		return nil, errPublicFundingRateHistoryMissingInstId
	}

	q := url.Values{}
	q.Set("instId", s.instId)
	if s.before != "" {
		q.Set("before", s.before)
	}
	if s.after != "" {
		q.Set("after", s.after)
	}
	if s.limit != nil {
		q.Set("limit", strconv.Itoa(*s.limit))
	}

	var data []FundingRateHistory
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/public/funding-rate-history", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
