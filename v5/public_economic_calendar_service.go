package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// EconomicCalendarEvent 表示经济日历数据。
//
// 说明：百分比/数值字段保持为 string（无损）；时间戳按 Unix 毫秒解析为 int64。
type EconomicCalendarEvent struct {
	CalendarId string `json:"calendarId"`
	Date       int64  `json:"date,string"`
	Region     string `json:"region"`
	Category   string `json:"category"`
	Event      string `json:"event"`
	RefDate    int64  `json:"refDate,string"`

	Actual      string `json:"actual"`
	Previous    string `json:"previous"`
	Forecast    string `json:"forecast"`
	PrevInitial string `json:"prevInitial"`
	Unit        string `json:"unit"`
	Ccy         string `json:"ccy"`

	DateSpan   string    `json:"dateSpan"`
	Importance string    `json:"importance"`
	UTime      int64     `json:"uTime,string"`
	TS         UnixMilli `json:"ts"`
}

// PublicEconomicCalendarService 获取经济日历数据。
//
// 注意：该接口需验证后使用，且仅支持实盘服务（按 OKX 文档）。
type PublicEconomicCalendarService struct {
	c *Client

	region     string
	importance string
	before     string
	after      string
	limit      *int
}

// NewPublicEconomicCalendarService 创建 PublicEconomicCalendarService。
func (c *Client) NewPublicEconomicCalendarService() *PublicEconomicCalendarService {
	return &PublicEconomicCalendarService{c: c}
}

// Region 设置国家、地区或实体（可选）。
func (s *PublicEconomicCalendarService) Region(region string) *PublicEconomicCalendarService {
	s.region = region
	return s
}

// Importance 设置重要性（可选：1=低，2=中等，3=高）。
func (s *PublicEconomicCalendarService) Importance(importance string) *PublicEconomicCalendarService {
	s.importance = importance
	return s
}

// Before 设置查询发布日期(date)之后的内容（可选；Unix 毫秒时间戳）。
func (s *PublicEconomicCalendarService) Before(before string) *PublicEconomicCalendarService {
	s.before = before
	return s
}

// After 设置查询发布日期(date)之前的内容（可选；Unix 毫秒时间戳；默认值为请求时刻的时间戳）。
func (s *PublicEconomicCalendarService) After(after string) *PublicEconomicCalendarService {
	s.after = after
	return s
}

// Limit 设置返回条数（最大 100，默认 100）。
func (s *PublicEconomicCalendarService) Limit(limit int) *PublicEconomicCalendarService {
	s.limit = &limit
	return s
}

// Do 获取经济日历数据（GET /api/v5/public/economic-calendar）。
func (s *PublicEconomicCalendarService) Do(ctx context.Context) ([]EconomicCalendarEvent, error) {
	q := url.Values{}
	if s.region != "" {
		q.Set("region", s.region)
	}
	if s.importance != "" {
		q.Set("importance", s.importance)
	}
	if s.before != "" {
		q.Set("before", s.before)
	}
	if s.after != "" {
		q.Set("after", s.after)
	}
	if s.limit != nil {
		q.Set("limit", strconv.Itoa(*s.limit))
	}

	var data []EconomicCalendarEvent
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/public/economic-calendar", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
