package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// EasyConvertHistoryService 获取一键兑换主流币历史记录。
type EasyConvertHistoryService struct {
	c *Client

	after  string
	before string
	limit  *int
}

// NewEasyConvertHistoryService 创建 EasyConvertHistoryService。
func (c *Client) NewEasyConvertHistoryService() *EasyConvertHistoryService {
	return &EasyConvertHistoryService{c: c}
}

// After 查询在此之前（不包含）的内容（可选：Unix 毫秒时间戳字符串）。
func (s *EasyConvertHistoryService) After(after string) *EasyConvertHistoryService {
	s.after = after
	return s
}

// Before 查询在此之后（不包含）的内容（可选：Unix 毫秒时间戳字符串）。
func (s *EasyConvertHistoryService) Before(before string) *EasyConvertHistoryService {
	s.before = before
	return s
}

func (s *EasyConvertHistoryService) Limit(limit int) *EasyConvertHistoryService {
	s.limit = &limit
	return s
}

// Do 获取一键兑换主流币历史记录（GET /api/v5/trade/easy-convert-history）。
func (s *EasyConvertHistoryService) Do(ctx context.Context) ([]EasyConvertHistory, error) {
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

	var data []EasyConvertHistory
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/trade/easy-convert-history", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
