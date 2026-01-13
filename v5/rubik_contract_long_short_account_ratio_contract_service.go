package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// RubikLongShortAccountRatioContractService 获取合约多空持仓人数比（按合约维度）。
type RubikLongShortAccountRatioContractService struct {
	c *Client

	instId string
	period string
	end    string
	begin  string
	limit  *int
}

// NewRubikLongShortAccountRatioContractService 创建 RubikLongShortAccountRatioContractService。
func (c *Client) NewRubikLongShortAccountRatioContractService() *RubikLongShortAccountRatioContractService {
	return &RubikLongShortAccountRatioContractService{c: c}
}

// InstId 设置产品ID（必填，如 BTC-USDT-SWAP）。
func (s *RubikLongShortAccountRatioContractService) InstId(instId string) *RubikLongShortAccountRatioContractService {
	s.instId = instId
	return s
}

// Period 设置时间粒度（可选）。
func (s *RubikLongShortAccountRatioContractService) Period(period string) *RubikLongShortAccountRatioContractService {
	s.period = period
	return s
}

// Begin 设置筛选的开始时间戳（毫秒字符串，可选）。
func (s *RubikLongShortAccountRatioContractService) Begin(begin string) *RubikLongShortAccountRatioContractService {
	s.begin = begin
	return s
}

// End 设置筛选的结束时间戳（毫秒字符串，可选）。
func (s *RubikLongShortAccountRatioContractService) End(end string) *RubikLongShortAccountRatioContractService {
	s.end = end
	return s
}

// Limit 设置分页返回的结果集数量（可选，最大 100）。
func (s *RubikLongShortAccountRatioContractService) Limit(limit int) *RubikLongShortAccountRatioContractService {
	s.limit = &limit
	return s
}

var errRubikLongShortAccountRatioContractMissingInstId = errors.New("okx: rubik long-short account ratio contract requires instId")

// Do 获取合约多空持仓人数比（GET /api/v5/rubik/stat/contracts/long-short-account-ratio-contract）。
func (s *RubikLongShortAccountRatioContractService) Do(ctx context.Context) ([]RubikTsRatio, error) {
	if s.instId == "" {
		return nil, errRubikLongShortAccountRatioContractMissingInstId
	}

	q := url.Values{}
	q.Set("instId", s.instId)
	if s.period != "" {
		q.Set("period", s.period)
	}
	if s.end != "" {
		q.Set("end", s.end)
	}
	if s.begin != "" {
		q.Set("begin", s.begin)
	}
	if s.limit != nil {
		q.Set("limit", strconv.Itoa(*s.limit))
	}

	var data []RubikTsRatio
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/rubik/stat/contracts/long-short-account-ratio-contract", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
