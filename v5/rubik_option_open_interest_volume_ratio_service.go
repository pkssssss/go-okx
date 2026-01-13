package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// RubikOptionOpenInterestVolumeRatioService 看涨/看跌期权合约持仓总量比/交易总量比。
type RubikOptionOpenInterestVolumeRatioService struct {
	c *Client

	ccy    string
	period string
}

// NewRubikOptionOpenInterestVolumeRatioService 创建 RubikOptionOpenInterestVolumeRatioService。
func (c *Client) NewRubikOptionOpenInterestVolumeRatioService() *RubikOptionOpenInterestVolumeRatioService {
	return &RubikOptionOpenInterestVolumeRatioService{c: c}
}

// Ccy 设置币种（必填）。
func (s *RubikOptionOpenInterestVolumeRatioService) Ccy(ccy string) *RubikOptionOpenInterestVolumeRatioService {
	s.ccy = ccy
	return s
}

// Period 设置时间粒度（可选，默认 8H）。
func (s *RubikOptionOpenInterestVolumeRatioService) Period(period string) *RubikOptionOpenInterestVolumeRatioService {
	s.period = period
	return s
}

var errRubikOptionOpenInterestVolumeRatioMissingCcy = errors.New("okx: rubik option open interest volume ratio requires ccy")

// Do 看涨/看跌期权合约持仓总量比/交易总量比（GET /api/v5/rubik/stat/option/open-interest-volume-ratio）。
func (s *RubikOptionOpenInterestVolumeRatioService) Do(ctx context.Context) ([]RubikOptionOpenInterestVolumeRatio, error) {
	if s.ccy == "" {
		return nil, errRubikOptionOpenInterestVolumeRatioMissingCcy
	}

	q := url.Values{}
	q.Set("ccy", s.ccy)
	if s.period != "" {
		q.Set("period", s.period)
	}

	var data []RubikOptionOpenInterestVolumeRatio
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/rubik/stat/option/open-interest-volume-ratio", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
