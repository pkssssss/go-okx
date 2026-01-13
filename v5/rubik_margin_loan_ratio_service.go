package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// RubikMarginLoanRatioService 获取杠杆多空比。
type RubikMarginLoanRatioService struct {
	c *Client

	ccy    string
	begin  string
	end    string
	period string
}

// NewRubikMarginLoanRatioService 创建 RubikMarginLoanRatioService。
func (c *Client) NewRubikMarginLoanRatioService() *RubikMarginLoanRatioService {
	return &RubikMarginLoanRatioService{c: c}
}

// Ccy 设置币种（必填）。
func (s *RubikMarginLoanRatioService) Ccy(ccy string) *RubikMarginLoanRatioService {
	s.ccy = ccy
	return s
}

// Begin 设置开始时间（毫秒字符串，可选）。
func (s *RubikMarginLoanRatioService) Begin(begin string) *RubikMarginLoanRatioService {
	s.begin = begin
	return s
}

// End 设置结束时间（毫秒字符串，可选）。
func (s *RubikMarginLoanRatioService) End(end string) *RubikMarginLoanRatioService {
	s.end = end
	return s
}

// Period 设置时间粒度（可选，默认 5m）。
func (s *RubikMarginLoanRatioService) Period(period string) *RubikMarginLoanRatioService {
	s.period = period
	return s
}

var errRubikMarginLoanRatioMissingCcy = errors.New("okx: rubik margin loan ratio requires ccy")

// Do 获取杠杆多空比（GET /api/v5/rubik/stat/margin/loan-ratio）。
func (s *RubikMarginLoanRatioService) Do(ctx context.Context) ([]RubikTsRatio, error) {
	if s.ccy == "" {
		return nil, errRubikMarginLoanRatioMissingCcy
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
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/rubik/stat/margin/loan-ratio", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
