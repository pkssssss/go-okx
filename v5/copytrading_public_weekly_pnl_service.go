package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// CopyTradingPublicWeeklyPnlService 获取交易员收益周表现（公共）。
type CopyTradingPublicWeeklyPnlService struct {
	c *Client

	instType   string
	uniqueCode string
}

// NewCopyTradingPublicWeeklyPnlService 创建 CopyTradingPublicWeeklyPnlService。
func (c *Client) NewCopyTradingPublicWeeklyPnlService() *CopyTradingPublicWeeklyPnlService {
	return &CopyTradingPublicWeeklyPnlService{c: c}
}

func (s *CopyTradingPublicWeeklyPnlService) InstType(instType string) *CopyTradingPublicWeeklyPnlService {
	s.instType = instType
	return s
}

func (s *CopyTradingPublicWeeklyPnlService) UniqueCode(uniqueCode string) *CopyTradingPublicWeeklyPnlService {
	s.uniqueCode = uniqueCode
	return s
}

var errCopyTradingPublicWeeklyPnlMissingUniqueCode = errors.New("okx: copytrading public weekly pnl requires uniqueCode")

// Do 获取交易员收益周表现（GET /api/v5/copytrading/public-weekly-pnl）。
func (s *CopyTradingPublicWeeklyPnlService) Do(ctx context.Context) ([]CopyTradingPnl, error) {
	if s.uniqueCode == "" {
		return nil, errCopyTradingPublicWeeklyPnlMissingUniqueCode
	}

	q := url.Values{}
	if s.instType != "" {
		q.Set("instType", s.instType)
	}
	q.Set("uniqueCode", s.uniqueCode)

	var data []CopyTradingPnl
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/copytrading/public-weekly-pnl", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
