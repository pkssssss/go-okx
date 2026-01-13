package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// CopyTradingPublicPnlService 获取交易员收益日表现（公共）。
type CopyTradingPublicPnlService struct {
	c *Client

	instType   string
	uniqueCode string
	lastDays   string
}

// NewCopyTradingPublicPnlService 创建 CopyTradingPublicPnlService。
func (c *Client) NewCopyTradingPublicPnlService() *CopyTradingPublicPnlService {
	return &CopyTradingPublicPnlService{c: c}
}

func (s *CopyTradingPublicPnlService) InstType(instType string) *CopyTradingPublicPnlService {
	s.instType = instType
	return s
}

func (s *CopyTradingPublicPnlService) UniqueCode(uniqueCode string) *CopyTradingPublicPnlService {
	s.uniqueCode = uniqueCode
	return s
}

func (s *CopyTradingPublicPnlService) LastDays(lastDays string) *CopyTradingPublicPnlService {
	s.lastDays = lastDays
	return s
}

var errCopyTradingPublicPnlMissingRequired = errors.New("okx: copytrading public pnl requires uniqueCode/lastDays")

// Do 获取交易员收益日表现（GET /api/v5/copytrading/public-pnl）。
func (s *CopyTradingPublicPnlService) Do(ctx context.Context) ([]CopyTradingPnl, error) {
	if s.uniqueCode == "" || s.lastDays == "" {
		return nil, errCopyTradingPublicPnlMissingRequired
	}

	q := url.Values{}
	if s.instType != "" {
		q.Set("instType", s.instType)
	}
	q.Set("uniqueCode", s.uniqueCode)
	q.Set("lastDays", s.lastDays)

	var data []CopyTradingPnl
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/copytrading/public-pnl", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
