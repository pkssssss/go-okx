package okx

import (
	"context"
	"errors"
	"net/http"
)

type rfqSetMMPConfigRequest struct {
	TimeInterval   string `json:"timeInterval"`
	FrozenInterval string `json:"frozenInterval"`
	CountLimit     string `json:"countLimit"`
}

// RFQSetMMPConfigService 设置 MMP 配置（仅适用于 RFQ maker）。
type RFQSetMMPConfigService struct {
	c *Client
	r rfqSetMMPConfigRequest
}

// NewRFQSetMMPConfigService 创建 RFQSetMMPConfigService。
func (c *Client) NewRFQSetMMPConfigService() *RFQSetMMPConfigService {
	return &RFQSetMMPConfigService{c: c}
}

// TimeInterval 设置时间窗口（毫秒，必填；"0" 代表不使用 MMP）。
func (s *RFQSetMMPConfigService) TimeInterval(timeInterval string) *RFQSetMMPConfigService {
	s.r.TimeInterval = timeInterval
	return s
}

// FrozenInterval 设置冻结时间长度（毫秒，必填；"0" 代表一直冻结直到重置）。
func (s *RFQSetMMPConfigService) FrozenInterval(frozenInterval string) *RFQSetMMPConfigService {
	s.r.FrozenInterval = frozenInterval
	return s
}

// CountLimit 设置尝试执行次数限制（必填）。
func (s *RFQSetMMPConfigService) CountLimit(countLimit string) *RFQSetMMPConfigService {
	s.r.CountLimit = countLimit
	return s
}

var (
	errRFQSetMMPConfigMissingRequired = errors.New("okx: rfq set mmp config requires timeInterval/frozenInterval/countLimit")
	errEmptyRFQSetMMPConfigResponse   = errors.New("okx: empty rfq set mmp config response")
)

// Do 设置 MMP 配置（POST /api/v5/rfq/mmp-config）。
func (s *RFQSetMMPConfigService) Do(ctx context.Context) (*RFQMMPConfig, error) {
	if s.r.TimeInterval == "" || s.r.FrozenInterval == "" || s.r.CountLimit == "" {
		return nil, errRFQSetMMPConfigMissingRequired
	}

	var data []RFQMMPConfig
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/rfq/mmp-config", nil, s.r, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyRFQSetMMPConfigResponse
	}
	return &data[0], nil
}
