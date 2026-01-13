package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// RubikOptionTakerBlockVolumeService 看跌/看涨期权合约主动买入/卖出量。
type RubikOptionTakerBlockVolumeService struct {
	c *Client

	ccy    string
	period string
}

// NewRubikOptionTakerBlockVolumeService 创建 RubikOptionTakerBlockVolumeService。
func (c *Client) NewRubikOptionTakerBlockVolumeService() *RubikOptionTakerBlockVolumeService {
	return &RubikOptionTakerBlockVolumeService{c: c}
}

// Ccy 设置币种（必填）。
func (s *RubikOptionTakerBlockVolumeService) Ccy(ccy string) *RubikOptionTakerBlockVolumeService {
	s.ccy = ccy
	return s
}

// Period 设置时间粒度（可选，默认 8H）。
func (s *RubikOptionTakerBlockVolumeService) Period(period string) *RubikOptionTakerBlockVolumeService {
	s.period = period
	return s
}

var errRubikOptionTakerBlockVolumeMissingCcy = errors.New("okx: rubik option taker block volume requires ccy")

// Do 看跌/看涨期权合约主动买入/卖出量（GET /api/v5/rubik/stat/option/taker-block-volume）。
func (s *RubikOptionTakerBlockVolumeService) Do(ctx context.Context) (*RubikOptionTakerBlockVolume, error) {
	if s.ccy == "" {
		return nil, errRubikOptionTakerBlockVolumeMissingCcy
	}

	q := url.Values{}
	q.Set("ccy", s.ccy)
	if s.period != "" {
		q.Set("period", s.period)
	}

	var data RubikOptionTakerBlockVolume
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/rubik/stat/option/taker-block-volume", q, nil, false, &data); err != nil {
		return nil, err
	}
	return &data, nil
}
