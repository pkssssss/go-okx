package okx

import (
	"context"
	"errors"
	"net/http"
)

type accountSetIsolatedModeRequest struct {
	IsoMode string `json:"isoMode"`
	Type    string `json:"type"`
}

// AccountSetIsolatedModeAck 表示逐仓保证金划转模式设置返回项。
type AccountSetIsolatedModeAck struct {
	IsoMode string `json:"isoMode"`
}

// AccountSetIsolatedModeService 逐仓交易设置（逐仓保证金划转模式）。
type AccountSetIsolatedModeService struct {
	c   *Client
	req accountSetIsolatedModeRequest
}

// NewAccountSetIsolatedModeService 创建 AccountSetIsolatedModeService。
func (c *Client) NewAccountSetIsolatedModeService() *AccountSetIsolatedModeService {
	return &AccountSetIsolatedModeService{c: c}
}

// IsoMode 设置逐仓保证金划转模式（必填：auto_transfers_ccy/automatic）。
func (s *AccountSetIsolatedModeService) IsoMode(isoMode string) *AccountSetIsolatedModeService {
	s.req.IsoMode = isoMode
	return s
}

// Type 设置业务线类型（必填：MARGIN/CONTRACTS）。
func (s *AccountSetIsolatedModeService) Type(typ string) *AccountSetIsolatedModeService {
	s.req.Type = typ
	return s
}

var (
	errAccountSetIsolatedModeMissingRequired = errors.New("okx: set isolated mode requires isoMode and type")
	errEmptyAccountSetIsolatedMode           = errors.New("okx: empty set isolated mode response")
	errInvalidAccountSetIsolatedMode         = errors.New("okx: invalid set isolated mode response")
)

func validateAccountSetIsolatedModeAck(ack *AccountSetIsolatedModeAck, req accountSetIsolatedModeRequest) error {
	if ack == nil || ack.IsoMode == "" || ack.IsoMode != req.IsoMode {
		return errInvalidAccountSetIsolatedMode
	}
	return nil
}

// Do 设置逐仓保证金划转模式（POST /api/v5/account/set-isolated-mode）。
func (s *AccountSetIsolatedModeService) Do(ctx context.Context) (*AccountSetIsolatedModeAck, error) {
	if s.req.IsoMode == "" || s.req.Type == "" {
		return nil, errAccountSetIsolatedModeMissingRequired
	}

	var data []AccountSetIsolatedModeAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/account/set-isolated-mode", nil, s.req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/account/set-isolated-mode", requestID, errEmptyAccountSetIsolatedMode)
	}
	if err := validateAccountSetIsolatedModeAck(&data[0], s.req); err != nil {
		return nil, newInvalidDataAPIError(http.MethodPost, "/api/v5/account/set-isolated-mode", requestID, err)
	}
	return &data[0], nil
}
