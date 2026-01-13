package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// RubikTakerVolumeService 获取主动买入/卖出情况。
type RubikTakerVolumeService struct {
	c *Client

	ccy      string
	instType string
	begin    string
	end      string
	period   string
}

// NewRubikTakerVolumeService 创建 RubikTakerVolumeService。
func (c *Client) NewRubikTakerVolumeService() *RubikTakerVolumeService {
	return &RubikTakerVolumeService{c: c}
}

// Ccy 设置币种（必填）。
func (s *RubikTakerVolumeService) Ccy(ccy string) *RubikTakerVolumeService {
	s.ccy = ccy
	return s
}

// InstType 设置产品类型（必填）：SPOT/CONTRACTS。
func (s *RubikTakerVolumeService) InstType(instType string) *RubikTakerVolumeService {
	s.instType = instType
	return s
}

// Begin 设置开始时间（毫秒字符串，可选）。
func (s *RubikTakerVolumeService) Begin(begin string) *RubikTakerVolumeService {
	s.begin = begin
	return s
}

// End 设置结束时间（毫秒字符串，可选）。
func (s *RubikTakerVolumeService) End(end string) *RubikTakerVolumeService {
	s.end = end
	return s
}

// Period 设置时间粒度（可选，默认 5m）。
func (s *RubikTakerVolumeService) Period(period string) *RubikTakerVolumeService {
	s.period = period
	return s
}

var errRubikTakerVolumeMissingRequired = errors.New("okx: rubik taker volume requires ccy and instType")

// Do 获取主动买入/卖出情况（GET /api/v5/rubik/stat/taker-volume）。
func (s *RubikTakerVolumeService) Do(ctx context.Context) ([]RubikTakerVolume, error) {
	if s.ccy == "" || s.instType == "" {
		return nil, errRubikTakerVolumeMissingRequired
	}

	q := url.Values{}
	q.Set("ccy", s.ccy)
	q.Set("instType", s.instType)
	if s.begin != "" {
		q.Set("begin", s.begin)
	}
	if s.end != "" {
		q.Set("end", s.end)
	}
	if s.period != "" {
		q.Set("period", s.period)
	}

	var data []RubikTakerVolume
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/rubik/stat/taker-volume", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
