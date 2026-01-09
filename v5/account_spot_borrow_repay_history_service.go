package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// AccountSpotBorrowRepayHistory 表示现货模式下借/还币历史记录。
type AccountSpotBorrowRepayHistory struct {
	AccBorrowed string    `json:"accBorrowed"`
	Amt         string    `json:"amt"`
	Ccy         string    `json:"ccy"`
	Type        string    `json:"type"`
	TS          UnixMilli `json:"ts"`
}

// AccountSpotBorrowRepayHistoryService 获取借/还币历史（现货模式）。
type AccountSpotBorrowRepayHistoryService struct {
	c *Client

	ccy       string
	eventType string
	after     string
	before    string
	limit     *int
}

// NewAccountSpotBorrowRepayHistoryService 创建 AccountSpotBorrowRepayHistoryService。
func (c *Client) NewAccountSpotBorrowRepayHistoryService() *AccountSpotBorrowRepayHistoryService {
	return &AccountSpotBorrowRepayHistoryService{c: c}
}

// Ccy 设置币种过滤（可选）。
func (s *AccountSpotBorrowRepayHistoryService) Ccy(ccy string) *AccountSpotBorrowRepayHistoryService {
	s.ccy = ccy
	return s
}

// Type 设置事件类型过滤（可选：auto_borrow/auto_repay/manual_borrow/manual_repay）。
func (s *AccountSpotBorrowRepayHistoryService) Type(eventType string) *AccountSpotBorrowRepayHistoryService {
	s.eventType = eventType
	return s
}

// After 请求发生时间 ts 之前（包含）的分页内容（Unix 毫秒字符串）。
func (s *AccountSpotBorrowRepayHistoryService) After(after string) *AccountSpotBorrowRepayHistoryService {
	s.after = after
	return s
}

// Before 请求发生时间 ts 之后（包含）的分页内容（Unix 毫秒字符串）。
func (s *AccountSpotBorrowRepayHistoryService) Before(before string) *AccountSpotBorrowRepayHistoryService {
	s.before = before
	return s
}

// Limit 设置返回条数（最大 100，默认 100）。
func (s *AccountSpotBorrowRepayHistoryService) Limit(limit int) *AccountSpotBorrowRepayHistoryService {
	s.limit = &limit
	return s
}

// Do 获取借/还币历史（GET /api/v5/account/spot-borrow-repay-history）。
func (s *AccountSpotBorrowRepayHistoryService) Do(ctx context.Context) ([]AccountSpotBorrowRepayHistory, error) {
	q := url.Values{}
	if s.ccy != "" {
		q.Set("ccy", s.ccy)
	}
	if s.eventType != "" {
		q.Set("type", s.eventType)
	}
	if s.after != "" {
		q.Set("after", s.after)
	}
	if s.before != "" {
		q.Set("before", s.before)
	}
	if s.limit != nil {
		q.Set("limit", strconv.Itoa(*s.limit))
	}
	if len(q) == 0 {
		q = nil
	}

	var data []AccountSpotBorrowRepayHistory
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/spot-borrow-repay-history", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
