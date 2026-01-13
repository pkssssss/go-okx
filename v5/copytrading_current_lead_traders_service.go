package okx

import (
	"context"
	"net/http"
	"net/url"
)

// CopyTradingCurrentLeadTradersService 获取当前跟随的交易员。
type CopyTradingCurrentLeadTradersService struct {
	c *Client

	instType string
}

// NewCopyTradingCurrentLeadTradersService 创建 CopyTradingCurrentLeadTradersService。
func (c *Client) NewCopyTradingCurrentLeadTradersService() *CopyTradingCurrentLeadTradersService {
	return &CopyTradingCurrentLeadTradersService{c: c}
}

// InstType 设置产品类型（默认 SWAP）。
func (s *CopyTradingCurrentLeadTradersService) InstType(instType string) *CopyTradingCurrentLeadTradersService {
	s.instType = instType
	return s
}

// Do 获取当前跟随的交易员（GET /api/v5/copytrading/current-lead-traders）。
func (s *CopyTradingCurrentLeadTradersService) Do(ctx context.Context) ([]CopyTradingCurrentLeadTrader, error) {
	var q url.Values
	if s.instType != "" {
		q = url.Values{}
		q.Set("instType", s.instType)
	}

	var data []CopyTradingCurrentLeadTrader
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/copytrading/current-lead-traders", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
