package okx

import (
	"context"
	"net/http"
)

// RFQMakerInstrumentSetting 表示 maker 可报价产品设置（按 instType 分组）。
type RFQMakerInstrumentSetting struct {
	InstType   string                          `json:"instType"`
	IncludeAll *bool                           `json:"includeAll,omitempty"`
	Data       []RFQMakerInstrumentSettingItem `json:"data"`
}

// RFQMakerInstrumentSettingItem 表示单个产品（或交易品种）的设置项。
type RFQMakerInstrumentSettingItem struct {
	InstFamily  string `json:"instFamily,omitempty"`
	InstId      string `json:"instId,omitempty"`
	MaxBlockSz  string `json:"maxBlockSz,omitempty"`
	MakerPxBand string `json:"makerPxBand,omitempty"`
}

// RFQMakerInstrumentSettingsService 获取可报价产品设置。
type RFQMakerInstrumentSettingsService struct {
	c *Client
}

// NewRFQMakerInstrumentSettingsService 创建 RFQMakerInstrumentSettingsService。
func (c *Client) NewRFQMakerInstrumentSettingsService() *RFQMakerInstrumentSettingsService {
	return &RFQMakerInstrumentSettingsService{c: c}
}

// Do 获取可报价产品设置（GET /api/v5/rfq/maker-instrument-settings）。
func (s *RFQMakerInstrumentSettingsService) Do(ctx context.Context) ([]RFQMakerInstrumentSetting, error) {
	var data []RFQMakerInstrumentSetting
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/rfq/maker-instrument-settings", nil, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
