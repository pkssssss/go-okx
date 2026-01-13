package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// RubikLongShortAccountRatioContractTopTraderService 获取精英交易员合约多空持仓人数比。
type RubikLongShortAccountRatioContractTopTraderService struct {
	c *Client

	instId string
	period string
	end    string
	begin  string
	limit  *int
}

// NewRubikLongShortAccountRatioContractTopTraderService 创建 RubikLongShortAccountRatioContractTopTraderService。
func (c *Client) NewRubikLongShortAccountRatioContractTopTraderService() *RubikLongShortAccountRatioContractTopTraderService {
	return &RubikLongShortAccountRatioContractTopTraderService{c: c}
}

// InstId 设置产品ID（必填，如 BTC-USDT-SWAP）。
func (s *RubikLongShortAccountRatioContractTopTraderService) InstId(instId string) *RubikLongShortAccountRatioContractTopTraderService {
	s.instId = instId
	return s
}

// Period 设置时间粒度（可选）。
func (s *RubikLongShortAccountRatioContractTopTraderService) Period(period string) *RubikLongShortAccountRatioContractTopTraderService {
	s.period = period
	return s
}

// Begin 设置筛选的开始时间戳（毫秒字符串，可选）。
func (s *RubikLongShortAccountRatioContractTopTraderService) Begin(begin string) *RubikLongShortAccountRatioContractTopTraderService {
	s.begin = begin
	return s
}

// End 设置筛选的结束时间戳（毫秒字符串，可选）。
func (s *RubikLongShortAccountRatioContractTopTraderService) End(end string) *RubikLongShortAccountRatioContractTopTraderService {
	s.end = end
	return s
}

// Limit 设置分页返回的结果集数量（可选，最大 100）。
func (s *RubikLongShortAccountRatioContractTopTraderService) Limit(limit int) *RubikLongShortAccountRatioContractTopTraderService {
	s.limit = &limit
	return s
}

var errRubikLongShortAccountRatioContractTopTraderMissingInstId = errors.New("okx: rubik long-short account ratio contract top trader requires instId")

// Do 获取精英交易员合约多空持仓人数比（GET /api/v5/rubik/stat/contracts/long-short-account-ratio-contract-top-trader）。
func (s *RubikLongShortAccountRatioContractTopTraderService) Do(ctx context.Context) ([]RubikTsRatio, error) {
	if s.instId == "" {
		return nil, errRubikLongShortAccountRatioContractTopTraderMissingInstId
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
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/rubik/stat/contracts/long-short-account-ratio-contract-top-trader", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
