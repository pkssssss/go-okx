package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// OneClickRepayHistoryV2Service 获取一键还债历史记录（新）。
type OneClickRepayHistoryV2Service struct {
	c *Client

	after  string
	before string
	limit  *int
}

// NewOneClickRepayHistoryV2Service 创建 OneClickRepayHistoryV2Service。
func (c *Client) NewOneClickRepayHistoryV2Service() *OneClickRepayHistoryV2Service {
	return &OneClickRepayHistoryV2Service{c: c}
}

func (s *OneClickRepayHistoryV2Service) After(after string) *OneClickRepayHistoryV2Service {
	s.after = after
	return s
}

func (s *OneClickRepayHistoryV2Service) Before(before string) *OneClickRepayHistoryV2Service {
	s.before = before
	return s
}

func (s *OneClickRepayHistoryV2Service) Limit(limit int) *OneClickRepayHistoryV2Service {
	s.limit = &limit
	return s
}

// Do 获取一键还债历史记录（新）（GET /api/v5/trade/one-click-repay-history-v2）。
func (s *OneClickRepayHistoryV2Service) Do(ctx context.Context) ([]OneClickRepayHistoryV2Item, error) {
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

	var data []OneClickRepayHistoryV2Item
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/trade/one-click-repay-history-v2", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
