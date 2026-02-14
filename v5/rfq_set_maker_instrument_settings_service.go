package okx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// RFQSetMakerInstrumentSettingsAck 表示设置可报价产品的返回项。
type RFQSetMakerInstrumentSettingsAck struct {
	Result bool `json:"result"`
}

// RFQSetMakerInstrumentSettingsService 设置可报价产品。
type RFQSetMakerInstrumentSettingsService struct {
	c *Client

	settings []RFQMakerInstrumentSetting
}

// NewRFQSetMakerInstrumentSettingsService 创建 RFQSetMakerInstrumentSettingsService。
func (c *Client) NewRFQSetMakerInstrumentSettingsService() *RFQSetMakerInstrumentSettingsService {
	return &RFQSetMakerInstrumentSettingsService{c: c}
}

// Settings 设置可报价产品配置（必填）。
func (s *RFQSetMakerInstrumentSettingsService) Settings(settings []RFQMakerInstrumentSetting) *RFQSetMakerInstrumentSettingsService {
	s.settings = settings
	return s
}

var (
	errRFQSetMakerInstrumentSettingsMissingSettings = errors.New("okx: rfq set maker instrument settings requires settings")
	errEmptyRFQSetMakerInstrumentSettingsResponse   = errors.New("okx: empty rfq set maker instrument settings response")
)

// Do 设置可报价产品（POST /api/v5/rfq/maker-instrument-settings）。
func (s *RFQSetMakerInstrumentSettingsService) Do(ctx context.Context) (*RFQSetMakerInstrumentSettingsAck, error) {
	if len(s.settings) == 0 {
		return nil, errRFQSetMakerInstrumentSettingsMissingSettings
	}

	for i, setting := range s.settings {
		if setting.InstType == "" {
			return nil, fmt.Errorf("okx: rfq set maker instrument settings[%d] missing instType", i)
		}
		if len(setting.Data) == 0 {
			return nil, fmt.Errorf("okx: rfq set maker instrument settings[%d] requires at least one data item", i)
		}
		for j, item := range setting.Data {
			if item.InstFamily == "" && item.InstId == "" {
				return nil, fmt.Errorf("okx: rfq set maker instrument settings[%d].data[%d] requires instFamily or instId", i, j)
			}
		}
	}

	var data []RFQSetMakerInstrumentSettingsAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/rfq/maker-instrument-settings", nil, s.settings, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/rfq/maker-instrument-settings", requestID, errEmptyRFQSetMakerInstrumentSettingsResponse)
	}
	if !data[0].Result {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodPost,
			RequestPath: "/api/v5/rfq/maker-instrument-settings",
			RequestID:   requestID,
			Code:        "0",
			Message:     "rfq maker instrument settings result is false",
		}
	}
	return &data[0], nil
}
