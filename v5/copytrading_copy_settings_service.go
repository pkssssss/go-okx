package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// CopyTradingCopySettingsService 获取跟单设置。
type CopyTradingCopySettingsService struct {
	c *Client

	instType   string
	uniqueCode string
}

// NewCopyTradingCopySettingsService 创建 CopyTradingCopySettingsService。
func (c *Client) NewCopyTradingCopySettingsService() *CopyTradingCopySettingsService {
	return &CopyTradingCopySettingsService{c: c}
}

// InstType 设置产品类型（默认 SWAP）。
func (s *CopyTradingCopySettingsService) InstType(instType string) *CopyTradingCopySettingsService {
	s.instType = instType
	return s
}

// UniqueCode 设置交易员唯一标识码（必填）。
func (s *CopyTradingCopySettingsService) UniqueCode(uniqueCode string) *CopyTradingCopySettingsService {
	s.uniqueCode = uniqueCode
	return s
}

var (
	errCopyTradingCopySettingsMissingUniqueCode = errors.New("okx: copytrading copy settings requires uniqueCode")
	errEmptyCopyTradingCopySettingsResponse     = errors.New("okx: empty copytrading copy settings response")
)

// Do 获取跟单设置（GET /api/v5/copytrading/copy-settings）。
func (s *CopyTradingCopySettingsService) Do(ctx context.Context) (*CopyTradingCopySettings, error) {
	if s.uniqueCode == "" {
		return nil, errCopyTradingCopySettingsMissingUniqueCode
	}

	q := url.Values{}
	if s.instType != "" {
		q.Set("instType", s.instType)
	}
	q.Set("uniqueCode", s.uniqueCode)

	var data []CopyTradingCopySettings
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/copytrading/copy-settings", q, nil, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyCopyTradingCopySettingsResponse
	}
	return &data[0], nil
}
