package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// RubikOptionOpenInterestVolumeStrikeService 看涨看跌持仓总量及交易总量（按执行价格分）。
type RubikOptionOpenInterestVolumeStrikeService struct {
	c *Client

	ccy     string
	expTime string
	period  string
}

// NewRubikOptionOpenInterestVolumeStrikeService 创建 RubikOptionOpenInterestVolumeStrikeService。
func (c *Client) NewRubikOptionOpenInterestVolumeStrikeService() *RubikOptionOpenInterestVolumeStrikeService {
	return &RubikOptionOpenInterestVolumeStrikeService{c: c}
}

// Ccy 设置币种（必填）。
func (s *RubikOptionOpenInterestVolumeStrikeService) Ccy(ccy string) *RubikOptionOpenInterestVolumeStrikeService {
	s.ccy = ccy
	return s
}

// ExpTime 设置到期日（必填，格式 YYYYMMDD，如 20210901）。
func (s *RubikOptionOpenInterestVolumeStrikeService) ExpTime(expTime string) *RubikOptionOpenInterestVolumeStrikeService {
	s.expTime = expTime
	return s
}

// Period 设置时间粒度（可选，默认 8H）。
func (s *RubikOptionOpenInterestVolumeStrikeService) Period(period string) *RubikOptionOpenInterestVolumeStrikeService {
	s.period = period
	return s
}

var errRubikOptionOpenInterestVolumeStrikeMissingRequired = errors.New("okx: rubik option open interest volume strike requires ccy and expTime")

// Do 看涨看跌持仓总量及交易总量（按执行价格分）（GET /api/v5/rubik/stat/option/open-interest-volume-strike）。
func (s *RubikOptionOpenInterestVolumeStrikeService) Do(ctx context.Context) ([]RubikOptionOpenInterestVolumeStrike, error) {
	if s.ccy == "" || s.expTime == "" {
		return nil, errRubikOptionOpenInterestVolumeStrikeMissingRequired
	}

	q := url.Values{}
	q.Set("ccy", s.ccy)
	q.Set("expTime", s.expTime)
	if s.period != "" {
		q.Set("period", s.period)
	}

	var data []RubikOptionOpenInterestVolumeStrike
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/rubik/stat/option/open-interest-volume-strike", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
