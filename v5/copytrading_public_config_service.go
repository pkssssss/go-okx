package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// CopyTradingPublicConfigService 获取跟单配置信息（公共）。
type CopyTradingPublicConfigService struct {
	c *Client

	instType string
}

// NewCopyTradingPublicConfigService 创建 CopyTradingPublicConfigService。
func (c *Client) NewCopyTradingPublicConfigService() *CopyTradingPublicConfigService {
	return &CopyTradingPublicConfigService{c: c}
}

// InstType 设置产品类型（默认 SWAP）。
func (s *CopyTradingPublicConfigService) InstType(instType string) *CopyTradingPublicConfigService {
	s.instType = instType
	return s
}

var errEmptyCopyTradingPublicConfigResponse = errors.New("okx: empty copytrading public config response")

// Do 获取跟单配置信息（GET /api/v5/copytrading/public-config）。
func (s *CopyTradingPublicConfigService) Do(ctx context.Context) (*CopyTradingPublicConfig, error) {
	var q url.Values
	if s.instType != "" {
		q = url.Values{}
		q.Set("instType", s.instType)
	}

	var data []CopyTradingPublicConfig
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/copytrading/public-config", q, nil, false, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyCopyTradingPublicConfigResponse
	}
	return &data[0], nil
}
