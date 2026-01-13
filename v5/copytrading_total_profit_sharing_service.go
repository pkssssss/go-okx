package okx

import (
	"context"
	"net/http"
	"net/url"
)

// CopyTradingTotalProfitSharingService 交易员历史分润汇总。
type CopyTradingTotalProfitSharingService struct {
	c *Client

	instType string
}

// NewCopyTradingTotalProfitSharingService 创建 CopyTradingTotalProfitSharingService。
func (c *Client) NewCopyTradingTotalProfitSharingService() *CopyTradingTotalProfitSharingService {
	return &CopyTradingTotalProfitSharingService{c: c}
}

// InstType 设置产品类型（默认返回所有）。
func (s *CopyTradingTotalProfitSharingService) InstType(instType string) *CopyTradingTotalProfitSharingService {
	s.instType = instType
	return s
}

// Do 获取历史分润汇总（GET /api/v5/copytrading/total-profit-sharing）。
func (s *CopyTradingTotalProfitSharingService) Do(ctx context.Context) ([]CopyTradingTotalProfitSharing, error) {
	var q url.Values
	if s.instType != "" {
		q = url.Values{}
		q.Set("instType", s.instType)
	}

	var data []CopyTradingTotalProfitSharing
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/copytrading/total-profit-sharing", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
