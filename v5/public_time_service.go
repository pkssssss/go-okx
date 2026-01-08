package okx

import (
	"context"
	"errors"
	"net/http"
	"time"
)

// SystemTime 表示 OKX 系统时间（毫秒）。
type SystemTime struct {
	TS int64 `json:"ts,string"`
}

// Time 返回 TS 对应的 UTC 时间。
func (t SystemTime) Time() time.Time {
	return time.UnixMilli(t.TS).UTC()
}

// PublicTimeService 获取系统时间。
type PublicTimeService struct {
	c *Client
}

// NewPublicTimeService 创建 PublicTimeService。
func (c *Client) NewPublicTimeService() *PublicTimeService {
	return &PublicTimeService{c: c}
}

var errEmptySystemTime = errors.New("okx: empty system time response")

// Do 获取系统时间（GET /api/v5/public/time）。
func (s *PublicTimeService) Do(ctx context.Context) (*SystemTime, error) {
	var data []SystemTime
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/public/time", nil, nil, false, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptySystemTime
	}
	return &data[0], nil
}
