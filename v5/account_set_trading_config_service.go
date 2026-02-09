package okx

import (
	"context"
	"errors"
	"net/http"
)

type accountSetTradingConfigRequest struct {
	Type     string `json:"type"`
	StgyType string `json:"stgyType,omitempty"`
}

// AccountSetTradingConfigAck 表示设置交易配置返回项。
type AccountSetTradingConfigAck struct {
	Type     string `json:"type"`
	StgyType string `json:"stgyType"`
}

// AccountSetTradingConfigService 设置交易配置。
type AccountSetTradingConfigService struct {
	c   *Client
	req accountSetTradingConfigRequest
}

// NewAccountSetTradingConfigService 创建 AccountSetTradingConfigService。
func (c *Client) NewAccountSetTradingConfigService() *AccountSetTradingConfigService {
	return &AccountSetTradingConfigService{c: c}
}

// Type 设置交易配置类型（必填；当前仅支持 stgyType）。
func (s *AccountSetTradingConfigService) Type(typ string) *AccountSetTradingConfigService {
	s.req.Type = typ
	return s
}

// StgyType 设置账号策略类型（可选；仅适用于 type=stgyType，0 普通策略；1 delta 中性）。
func (s *AccountSetTradingConfigService) StgyType(stgyType string) *AccountSetTradingConfigService {
	s.req.StgyType = stgyType
	return s
}

var (
	errAccountSetTradingConfigMissingType     = errors.New("okx: set trading config requires type")
	errAccountSetTradingConfigMissingStgyType = errors.New("okx: set trading config type=stgyType requires stgyType")
	errEmptyAccountSetTradingConfig           = errors.New("okx: empty set trading config response")
	errInvalidAccountSetTradingConfig         = errors.New("okx: invalid set trading config response")
)

func validateAccountSetTradingConfigAck(ack *AccountSetTradingConfigAck, req accountSetTradingConfigRequest) error {
	if ack == nil || ack.Type == "" || ack.Type != req.Type {
		return errInvalidAccountSetTradingConfig
	}
	if req.Type == "stgyType" && (ack.StgyType == "" || ack.StgyType != req.StgyType) {
		return errInvalidAccountSetTradingConfig
	}
	return nil
}

// Do 设置交易配置（POST /api/v5/account/set-trading-config）。
func (s *AccountSetTradingConfigService) Do(ctx context.Context) (*AccountSetTradingConfigAck, error) {
	if s.req.Type == "" {
		return nil, errAccountSetTradingConfigMissingType
	}
	if s.req.Type == "stgyType" && s.req.StgyType == "" {
		return nil, errAccountSetTradingConfigMissingStgyType
	}

	var data []AccountSetTradingConfigAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/account/set-trading-config", nil, s.req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyAccountSetTradingConfig
	}
	if err := validateAccountSetTradingConfigAck(&data[0], s.req); err != nil {
		return nil, err
	}
	return &data[0], nil
}
