package okx

import (
	"context"
	"net/http"
	"net/url"
)

// AccountMMPConfig 表示 MMP 配置信息。
type AccountMMPConfig struct {
	FrozenInterval string `json:"frozenInterval"`
	InstFamily     string `json:"instFamily"`
	MMPFrozen      bool   `json:"mmpFrozen"`
	MMPFrozenUntil string `json:"mmpFrozenUntil"`
	QtyLimit       string `json:"qtyLimit"`
	TimeInterval   string `json:"timeInterval"`
}

// AccountMMPConfigService 查看 MMP 配置。
type AccountMMPConfigService struct {
	c *Client

	instFamily string
}

// NewAccountMMPConfigService 创建 AccountMMPConfigService。
func (c *Client) NewAccountMMPConfigService() *AccountMMPConfigService {
	return &AccountMMPConfigService{c: c}
}

// InstFamily 设置交易品种过滤（可选）。
func (s *AccountMMPConfigService) InstFamily(instFamily string) *AccountMMPConfigService {
	s.instFamily = instFamily
	return s
}

// Do 查看 MMP 配置（GET /api/v5/account/mmp-config）。
func (s *AccountMMPConfigService) Do(ctx context.Context) ([]AccountMMPConfig, error) {
	var q url.Values
	if s.instFamily != "" {
		q = url.Values{}
		q.Set("instFamily", s.instFamily)
	}

	var data []AccountMMPConfig
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/mmp-config", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
