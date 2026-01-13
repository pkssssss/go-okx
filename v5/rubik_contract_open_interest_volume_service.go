package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// RubikContractsOpenInterestVolumeService 获取合约持仓量及交易量。
type RubikContractsOpenInterestVolumeService struct {
	c *Client

	ccy    string
	begin  string
	end    string
	period string
}

// NewRubikContractsOpenInterestVolumeService 创建 RubikContractsOpenInterestVolumeService。
func (c *Client) NewRubikContractsOpenInterestVolumeService() *RubikContractsOpenInterestVolumeService {
	return &RubikContractsOpenInterestVolumeService{c: c}
}

// Ccy 设置币种（必填）。
func (s *RubikContractsOpenInterestVolumeService) Ccy(ccy string) *RubikContractsOpenInterestVolumeService {
	s.ccy = ccy
	return s
}

// Begin 设置开始时间（毫秒字符串，可选）。
func (s *RubikContractsOpenInterestVolumeService) Begin(begin string) *RubikContractsOpenInterestVolumeService {
	s.begin = begin
	return s
}

// End 设置结束时间（毫秒字符串，可选）。
func (s *RubikContractsOpenInterestVolumeService) End(end string) *RubikContractsOpenInterestVolumeService {
	s.end = end
	return s
}

// Period 设置时间粒度（可选，默认 5m）。
func (s *RubikContractsOpenInterestVolumeService) Period(period string) *RubikContractsOpenInterestVolumeService {
	s.period = period
	return s
}

var errRubikContractsOpenInterestVolumeMissingCcy = errors.New("okx: rubik contracts open interest volume requires ccy")

// Do 获取合约持仓量及交易量（GET /api/v5/rubik/stat/contracts/open-interest-volume）。
func (s *RubikContractsOpenInterestVolumeService) Do(ctx context.Context) ([]RubikOpenInterestVolume, error) {
	if s.ccy == "" {
		return nil, errRubikContractsOpenInterestVolumeMissingCcy
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

	var data []RubikOpenInterestVolume
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/rubik/stat/contracts/open-interest-volume", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
