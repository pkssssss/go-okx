package okx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type accountLevelSwitchPresetRequest struct {
	AcctLv         string `json:"acctLv"`
	Lever          string `json:"lever,omitempty"`
	RiskOffsetType string `json:"riskOffsetType,omitempty"`
}

// AccountLevelSwitchPresetAck 表示预设置账户模式切换返回项。
type AccountLevelSwitchPresetAck struct {
	AcctLv         string `json:"acctLv"`
	CurAcctLv      string `json:"curAcctLv"`
	Lever          string `json:"lever"`
	RiskOffsetType string `json:"riskOffsetType"`
}

// AccountLevelSwitchPresetService 预设置账户模式切换。
type AccountLevelSwitchPresetService struct {
	c *Client
	r accountLevelSwitchPresetRequest
}

// NewAccountLevelSwitchPresetService 创建 AccountLevelSwitchPresetService。
func (c *Client) NewAccountLevelSwitchPresetService() *AccountLevelSwitchPresetService {
	return &AccountLevelSwitchPresetService{c: c}
}

// AcctLv 设置目标账户模式（必填）。
func (s *AccountLevelSwitchPresetService) AcctLv(acctLv string) *AccountLevelSwitchPresetService {
	s.r.AcctLv = acctLv
	return s
}

// Lever 设置杠杆倍数（可选）。
func (s *AccountLevelSwitchPresetService) Lever(lever string) *AccountLevelSwitchPresetService {
	s.r.Lever = lever
	return s
}

// RiskOffsetType 设置风险对冲模式（已弃用，仅为兼容保留）。
func (s *AccountLevelSwitchPresetService) RiskOffsetType(riskOffsetType string) *AccountLevelSwitchPresetService {
	s.r.RiskOffsetType = riskOffsetType
	return s
}

var (
	errAccountLevelSwitchPresetMissingAcctLv = errors.New("okx: account level switch preset requires acctLv")
	errEmptyAccountLevelSwitchPreset         = errors.New("okx: empty account level switch preset response")
)

// Do 预设置账户模式切换（POST /api/v5/account/account-level-switch-preset）。
func (s *AccountLevelSwitchPresetService) Do(ctx context.Context) (*AccountLevelSwitchPresetAck, error) {
	if s.r.AcctLv == "" {
		return nil, errAccountLevelSwitchPresetMissingAcctLv
	}

	var raw json.RawMessage
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/account/account-level-switch-preset", nil, s.r, true, &raw); err != nil {
		return nil, err
	}

	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 || string(raw) == "null" {
		return nil, errEmptyAccountLevelSwitchPreset
	}

	switch raw[0] {
	case '[':
		var data []AccountLevelSwitchPresetAck
		if err := json.Unmarshal(raw, &data); err != nil {
			return nil, err
		}
		if len(data) == 0 {
			return nil, errEmptyAccountLevelSwitchPreset
		}
		return &data[0], nil
	case '{':
		var data AccountLevelSwitchPresetAck
		if err := json.Unmarshal(raw, &data); err != nil {
			return nil, err
		}
		return &data, nil
	default:
		return nil, errors.New("okx: invalid account level switch preset response")
	}
}
