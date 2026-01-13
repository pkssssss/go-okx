package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// RubikOptionOpenInterestVolumeExpiryService 看涨看跌持仓总量及交易总量（按到期日分）。
type RubikOptionOpenInterestVolumeExpiryService struct {
	c *Client

	ccy    string
	period string
}

// NewRubikOptionOpenInterestVolumeExpiryService 创建 RubikOptionOpenInterestVolumeExpiryService。
func (c *Client) NewRubikOptionOpenInterestVolumeExpiryService() *RubikOptionOpenInterestVolumeExpiryService {
	return &RubikOptionOpenInterestVolumeExpiryService{c: c}
}

// Ccy 设置币种（必填）。
func (s *RubikOptionOpenInterestVolumeExpiryService) Ccy(ccy string) *RubikOptionOpenInterestVolumeExpiryService {
	s.ccy = ccy
	return s
}

// Period 设置时间粒度（可选，默认 8H）。
func (s *RubikOptionOpenInterestVolumeExpiryService) Period(period string) *RubikOptionOpenInterestVolumeExpiryService {
	s.period = period
	return s
}

var errRubikOptionOpenInterestVolumeExpiryMissingCcy = errors.New("okx: rubik option open interest volume expiry requires ccy")

// Do 看涨看跌持仓总量及交易总量（按到期日分）（GET /api/v5/rubik/stat/option/open-interest-volume-expiry）。
func (s *RubikOptionOpenInterestVolumeExpiryService) Do(ctx context.Context) ([]RubikOptionOpenInterestVolumeExpiry, error) {
	if s.ccy == "" {
		return nil, errRubikOptionOpenInterestVolumeExpiryMissingCcy
	}

	q := url.Values{}
	q.Set("ccy", s.ccy)
	if s.period != "" {
		q.Set("period", s.period)
	}

	var data []RubikOptionOpenInterestVolumeExpiry
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/rubik/stat/option/open-interest-volume-expiry", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
