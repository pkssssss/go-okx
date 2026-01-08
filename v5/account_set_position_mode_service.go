package okx

import (
	"context"
	"errors"
	"net/http"
)

type accountSetPositionModeRequest struct {
	PosMode string `json:"posMode"`
}

// AccountSetPositionModeAck 表示设置持仓模式返回项。
type AccountSetPositionModeAck struct {
	PosMode string `json:"posMode"`
}

// AccountSetPositionModeService 设置持仓模式。
type AccountSetPositionModeService struct {
	c       *Client
	posMode string
}

// NewAccountSetPositionModeService 创建 AccountSetPositionModeService。
func (c *Client) NewAccountSetPositionModeService() *AccountSetPositionModeService {
	return &AccountSetPositionModeService{c: c}
}

// PosMode 设置持仓方式（必填：long_short_mode / net_mode）。
func (s *AccountSetPositionModeService) PosMode(posMode string) *AccountSetPositionModeService {
	s.posMode = posMode
	return s
}

var errAccountSetPositionModeMissingPosMode = errors.New("okx: set position mode requires posMode")
var errEmptyAccountSetPositionMode = errors.New("okx: empty set position mode response")

// Do 设置持仓模式（POST /api/v5/account/set-position-mode）。
func (s *AccountSetPositionModeService) Do(ctx context.Context) (*AccountSetPositionModeAck, error) {
	if s.posMode == "" {
		return nil, errAccountSetPositionModeMissingPosMode
	}

	req := accountSetPositionModeRequest{PosMode: s.posMode}

	var data []AccountSetPositionModeAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/account/set-position-mode", nil, req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyAccountSetPositionMode
	}
	return &data[0], nil
}
