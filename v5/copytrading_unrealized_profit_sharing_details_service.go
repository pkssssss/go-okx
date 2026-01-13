package okx

import (
	"context"
	"net/http"
	"net/url"
)

// CopyTradingUnrealizedProfitSharingDetailsService 交易员待分润明细。
type CopyTradingUnrealizedProfitSharingDetailsService struct {
	c *Client

	instType string
}

// NewCopyTradingUnrealizedProfitSharingDetailsService 创建 CopyTradingUnrealizedProfitSharingDetailsService。
func (c *Client) NewCopyTradingUnrealizedProfitSharingDetailsService() *CopyTradingUnrealizedProfitSharingDetailsService {
	return &CopyTradingUnrealizedProfitSharingDetailsService{c: c}
}

// InstType 设置产品类型（默认返回所有）。
func (s *CopyTradingUnrealizedProfitSharingDetailsService) InstType(instType string) *CopyTradingUnrealizedProfitSharingDetailsService {
	s.instType = instType
	return s
}

// Do 获取待分润明细（GET /api/v5/copytrading/unrealized-profit-sharing-details）。
func (s *CopyTradingUnrealizedProfitSharingDetailsService) Do(ctx context.Context) ([]CopyTradingUnrealizedProfitSharingDetail, error) {
	var q url.Values
	if s.instType != "" {
		q = url.Values{}
		q.Set("instType", s.instType)
	}

	var data []CopyTradingUnrealizedProfitSharingDetail
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/copytrading/unrealized-profit-sharing-details", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
