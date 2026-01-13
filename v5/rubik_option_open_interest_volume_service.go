package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// RubikOptionOpenInterestVolumeService 获取期权持仓量及交易量。
type RubikOptionOpenInterestVolumeService struct {
	c *Client

	ccy    string
	period string
}

// NewRubikOptionOpenInterestVolumeService 创建 RubikOptionOpenInterestVolumeService。
func (c *Client) NewRubikOptionOpenInterestVolumeService() *RubikOptionOpenInterestVolumeService {
	return &RubikOptionOpenInterestVolumeService{c: c}
}

// Ccy 设置币种（必填）。
func (s *RubikOptionOpenInterestVolumeService) Ccy(ccy string) *RubikOptionOpenInterestVolumeService {
	s.ccy = ccy
	return s
}

// Period 设置时间粒度（可选，默认 8H）。
func (s *RubikOptionOpenInterestVolumeService) Period(period string) *RubikOptionOpenInterestVolumeService {
	s.period = period
	return s
}

var errRubikOptionOpenInterestVolumeMissingCcy = errors.New("okx: rubik option open interest volume requires ccy")

// Do 获取期权持仓量及交易量（GET /api/v5/rubik/stat/option/open-interest-volume）。
func (s *RubikOptionOpenInterestVolumeService) Do(ctx context.Context) ([]RubikOpenInterestVolume, error) {
	if s.ccy == "" {
		return nil, errRubikOptionOpenInterestVolumeMissingCcy
	}

	q := url.Values{}
	q.Set("ccy", s.ccy)
	if s.period != "" {
		q.Set("period", s.period)
	}

	var data []RubikOpenInterestVolume
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/rubik/stat/option/open-interest-volume", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
