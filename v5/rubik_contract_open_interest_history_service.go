package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// RubikOpenInterestHistoryService 获取合约持仓量历史。
type RubikOpenInterestHistoryService struct {
	c *Client

	instId string
	period string
	end    string
	begin  string
	limit  *int
}

// NewRubikOpenInterestHistoryService 创建 RubikOpenInterestHistoryService。
func (c *Client) NewRubikOpenInterestHistoryService() *RubikOpenInterestHistoryService {
	return &RubikOpenInterestHistoryService{c: c}
}

// InstId 设置产品ID（必填，如 BTC-USDT-SWAP）。
func (s *RubikOpenInterestHistoryService) InstId(instId string) *RubikOpenInterestHistoryService {
	s.instId = instId
	return s
}

// Period 设置时间粒度（可选）。
func (s *RubikOpenInterestHistoryService) Period(period string) *RubikOpenInterestHistoryService {
	s.period = period
	return s
}

// Begin 设置筛选的开始时间戳（毫秒字符串，可选）。
func (s *RubikOpenInterestHistoryService) Begin(begin string) *RubikOpenInterestHistoryService {
	s.begin = begin
	return s
}

// End 设置筛选的结束时间戳（毫秒字符串，可选）。
func (s *RubikOpenInterestHistoryService) End(end string) *RubikOpenInterestHistoryService {
	s.end = end
	return s
}

// Limit 设置分页返回的结果集数量（可选，最大 100）。
func (s *RubikOpenInterestHistoryService) Limit(limit int) *RubikOpenInterestHistoryService {
	s.limit = &limit
	return s
}

var errRubikOpenInterestHistoryMissingInstId = errors.New("okx: rubik open interest history requires instId")

// Do 获取合约持仓量历史（GET /api/v5/rubik/stat/contracts/open-interest-history）。
func (s *RubikOpenInterestHistoryService) Do(ctx context.Context) ([]RubikOpenInterestHistory, error) {
	if s.instId == "" {
		return nil, errRubikOpenInterestHistoryMissingInstId
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

	var data []RubikOpenInterestHistory
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/rubik/stat/contracts/open-interest-history", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
