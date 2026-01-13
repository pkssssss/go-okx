package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// CopyTradingPublicStatsService 获取交易员带单情况（公共）。
type CopyTradingPublicStatsService struct {
	c *Client

	instType   string
	uniqueCode string
	lastDays   string
}

// NewCopyTradingPublicStatsService 创建 CopyTradingPublicStatsService。
func (c *Client) NewCopyTradingPublicStatsService() *CopyTradingPublicStatsService {
	return &CopyTradingPublicStatsService{c: c}
}

func (s *CopyTradingPublicStatsService) InstType(instType string) *CopyTradingPublicStatsService {
	s.instType = instType
	return s
}

func (s *CopyTradingPublicStatsService) UniqueCode(uniqueCode string) *CopyTradingPublicStatsService {
	s.uniqueCode = uniqueCode
	return s
}

func (s *CopyTradingPublicStatsService) LastDays(lastDays string) *CopyTradingPublicStatsService {
	s.lastDays = lastDays
	return s
}

var (
	errCopyTradingPublicStatsMissingRequired = errors.New("okx: copytrading public stats requires uniqueCode/lastDays")
	errEmptyCopyTradingPublicStatsResponse   = errors.New("okx: empty copytrading public stats response")
)

// Do 获取交易员带单情况（GET /api/v5/copytrading/public-stats）。
func (s *CopyTradingPublicStatsService) Do(ctx context.Context) (*CopyTradingPublicStats, error) {
	if s.uniqueCode == "" || s.lastDays == "" {
		return nil, errCopyTradingPublicStatsMissingRequired
	}

	q := url.Values{}
	if s.instType != "" {
		q.Set("instType", s.instType)
	}
	q.Set("uniqueCode", s.uniqueCode)
	q.Set("lastDays", s.lastDays)

	var data []CopyTradingPublicStats
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/copytrading/public-stats", q, nil, false, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyCopyTradingPublicStatsResponse
	}
	return &data[0], nil
}
