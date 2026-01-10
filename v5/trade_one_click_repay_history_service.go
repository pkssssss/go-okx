package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// OneClickRepayHistoryService 获取一键还债历史记录（跨币种保证金/组合保证金）。
type OneClickRepayHistoryService struct {
	c *Client

	after  string
	before string
	limit  *int
}

// NewOneClickRepayHistoryService 创建 OneClickRepayHistoryService。
func (c *Client) NewOneClickRepayHistoryService() *OneClickRepayHistoryService {
	return &OneClickRepayHistoryService{c: c}
}

func (s *OneClickRepayHistoryService) After(after string) *OneClickRepayHistoryService {
	s.after = after
	return s
}

func (s *OneClickRepayHistoryService) Before(before string) *OneClickRepayHistoryService {
	s.before = before
	return s
}

func (s *OneClickRepayHistoryService) Limit(limit int) *OneClickRepayHistoryService {
	s.limit = &limit
	return s
}

// Do 获取一键还债历史记录（跨币种保证金/组合保证金）（GET /api/v5/trade/one-click-repay-history）。
func (s *OneClickRepayHistoryService) Do(ctx context.Context) ([]OneClickRepayHistory, error) {
	var q url.Values
	if s.after != "" || s.before != "" || s.limit != nil {
		q = url.Values{}
		if s.after != "" {
			q.Set("after", s.after)
		}
		if s.before != "" {
			q.Set("before", s.before)
		}
		if s.limit != nil {
			q.Set("limit", strconv.Itoa(*s.limit))
		}
	}

	var data []OneClickRepayHistory
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/trade/one-click-repay-history", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
