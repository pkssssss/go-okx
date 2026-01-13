package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// RubikLongShortAccountRatioService 获取多空持仓人数比（按币种维度）。
type RubikLongShortAccountRatioService struct {
	c *Client

	ccy    string
	begin  string
	end    string
	period string
}

// NewRubikLongShortAccountRatioService 创建 RubikLongShortAccountRatioService。
func (c *Client) NewRubikLongShortAccountRatioService() *RubikLongShortAccountRatioService {
	return &RubikLongShortAccountRatioService{c: c}
}

// Ccy 设置币种（必填）。
func (s *RubikLongShortAccountRatioService) Ccy(ccy string) *RubikLongShortAccountRatioService {
	s.ccy = ccy
	return s
}

// Begin 设置开始时间（毫秒字符串，可选）。
func (s *RubikLongShortAccountRatioService) Begin(begin string) *RubikLongShortAccountRatioService {
	s.begin = begin
	return s
}

// End 设置结束时间（毫秒字符串，可选）。
func (s *RubikLongShortAccountRatioService) End(end string) *RubikLongShortAccountRatioService {
	s.end = end
	return s
}

// Period 设置时间粒度（可选，默认 5m）。
func (s *RubikLongShortAccountRatioService) Period(period string) *RubikLongShortAccountRatioService {
	s.period = period
	return s
}

var errRubikLongShortAccountRatioMissingCcy = errors.New("okx: rubik long-short account ratio requires ccy")

// Do 获取多空持仓人数比（GET /api/v5/rubik/stat/contracts/long-short-account-ratio）。
func (s *RubikLongShortAccountRatioService) Do(ctx context.Context) ([]RubikTsRatio, error) {
	if s.ccy == "" {
		return nil, errRubikLongShortAccountRatioMissingCcy
	}

	q := url.Values{}
	q.Set("ccy", s.ccy)
	if s.begin != "" {
		q.Set("begin", s.begin)
	}
	if s.end != "" {
		q.Set("end", s.end)
	}
	if s.period != "" {
		q.Set("period", s.period)
	}

	var data []RubikTsRatio
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/rubik/stat/contracts/long-short-account-ratio", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
