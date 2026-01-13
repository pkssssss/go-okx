package okx

import (
	"context"
	"errors"
	"net/http"
)

// CopyTradingConfigService 获取账户配置信息。
type CopyTradingConfigService struct {
	c *Client
}

// NewCopyTradingConfigService 创建 CopyTradingConfigService。
func (c *Client) NewCopyTradingConfigService() *CopyTradingConfigService {
	return &CopyTradingConfigService{c: c}
}

var errEmptyCopyTradingConfigResponse = errors.New("okx: empty copytrading config response")

// Do 获取账户配置信息（GET /api/v5/copytrading/config）。
func (s *CopyTradingConfigService) Do(ctx context.Context) (*CopyTradingConfig, error) {
	var data []CopyTradingConfig
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/copytrading/config", nil, nil, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyCopyTradingConfigResponse
	}
	return &data[0], nil
}
