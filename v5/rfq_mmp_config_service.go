package okx

import (
	"context"
	"net/http"
)

// RFQMMPConfig 表示 RFQ Maker 的 MMP 配置信息。
type RFQMMPConfig struct {
	TimeInterval   string `json:"timeInterval"`
	FrozenInterval string `json:"frozenInterval"`
	CountLimit     string `json:"countLimit"`

	MMPFrozen      bool   `json:"mmpFrozen"`
	MMPFrozenUntil string `json:"mmpFrozenUntil"`
}

// RFQMMPConfigService 查看 MMP 配置。
type RFQMMPConfigService struct {
	c *Client
}

// NewRFQMMPConfigService 创建 RFQMMPConfigService。
func (c *Client) NewRFQMMPConfigService() *RFQMMPConfigService {
	return &RFQMMPConfigService{c: c}
}

// Do 查看 MMP 配置（GET /api/v5/rfq/mmp-config）。
func (s *RFQMMPConfigService) Do(ctx context.Context) ([]RFQMMPConfig, error) {
	var data []RFQMMPConfig
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/rfq/mmp-config", nil, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
