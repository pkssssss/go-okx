package okx

import (
	"context"
	"net/http"
	"net/url"
)

// CopyTradingTotalUnrealizedProfitSharingService 交易员待分润汇总。
type CopyTradingTotalUnrealizedProfitSharingService struct {
	c *Client

	instType string
}

// NewCopyTradingTotalUnrealizedProfitSharingService 创建 CopyTradingTotalUnrealizedProfitSharingService。
func (c *Client) NewCopyTradingTotalUnrealizedProfitSharingService() *CopyTradingTotalUnrealizedProfitSharingService {
	return &CopyTradingTotalUnrealizedProfitSharingService{c: c}
}

// InstType 设置产品类型（默认 SWAP）。
func (s *CopyTradingTotalUnrealizedProfitSharingService) InstType(instType string) *CopyTradingTotalUnrealizedProfitSharingService {
	s.instType = instType
	return s
}

// Do 获取待分润汇总（GET /api/v5/copytrading/total-unrealized-profit-sharing）。
func (s *CopyTradingTotalUnrealizedProfitSharingService) Do(ctx context.Context) ([]CopyTradingTotalUnrealizedProfitSharing, error) {
	var q url.Values
	if s.instType != "" {
		q = url.Values{}
		q.Set("instType", s.instType)
	}

	var data []CopyTradingTotalUnrealizedProfitSharing
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/copytrading/total-unrealized-profit-sharing", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
