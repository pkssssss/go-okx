package okx

import (
	"context"
	"net/http"
	"net/url"
)

// SystemStatusService 获取系统维护状态。
type SystemStatusService struct {
	c     *Client
	state string
}

// NewSystemStatusService 创建 SystemStatusService。
func (c *Client) NewSystemStatusService() *SystemStatusService {
	return &SystemStatusService{c: c}
}

// State 过滤系统状态（可选：scheduled/ongoing/pre_open/completed/canceled）。
// 不填写默认返回 scheduled/ongoing/pre_open。
func (s *SystemStatusService) State(state string) *SystemStatusService {
	s.state = state
	return s
}

// Do 获取系统维护状态（GET /api/v5/system/status）。
func (s *SystemStatusService) Do(ctx context.Context) ([]SystemStatus, error) {
	var q url.Values
	if s.state != "" {
		q = url.Values{}
		q.Set("state", s.state)
	}

	var data []SystemStatus
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/system/status", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
