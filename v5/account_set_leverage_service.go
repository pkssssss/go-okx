package okx

import (
	"context"
	"errors"
	"net/http"
)

type accountSetLeverageRequest struct {
	InstId  string `json:"instId,omitempty"`
	Ccy     string `json:"ccy,omitempty"`
	Lever   string `json:"lever"`
	PosSide string `json:"posSide,omitempty"`
	MgnMode string `json:"mgnMode"`
}

// AccountSetLeverageAck 表示设置杠杆倍数返回项。
type AccountSetLeverageAck struct {
	Lever   string `json:"lever"`
	MgnMode string `json:"mgnMode"`
	InstId  string `json:"instId"`
	PosSide string `json:"posSide"`
}

// AccountSetLeverageService 设置杠杆倍数。
type AccountSetLeverageService struct {
	c   *Client
	req accountSetLeverageRequest
}

// NewAccountSetLeverageService 创建 AccountSetLeverageService。
func (c *Client) NewAccountSetLeverageService() *AccountSetLeverageService {
	return &AccountSetLeverageService{c: c}
}

// InstId 设置产品 ID（币对/合约）。
func (s *AccountSetLeverageService) InstId(instId string) *AccountSetLeverageService {
	s.req.InstId = instId
	return s
}

// Ccy 设置保证金币种（用于自动借币模式下币种维度杠杆，仅全仓币币杠杆适用）。
func (s *AccountSetLeverageService) Ccy(ccy string) *AccountSetLeverageService {
	s.req.Ccy = ccy
	return s
}

// Lever 设置杠杆倍数（必填）。
func (s *AccountSetLeverageService) Lever(lever string) *AccountSetLeverageService {
	s.req.Lever = lever
	return s
}

// MgnMode 设置保证金模式（必填：isolated/cross）。
func (s *AccountSetLeverageService) MgnMode(mgnMode string) *AccountSetLeverageService {
	s.req.MgnMode = mgnMode
	return s
}

// PosSide 设置持仓方向（可选：long/short），仅逐仓交割/永续的开平仓模式适用。
func (s *AccountSetLeverageService) PosSide(posSide string) *AccountSetLeverageService {
	s.req.PosSide = posSide
	return s
}

var (
	errAccountSetLeverageMissingRequired       = errors.New("okx: set leverage requires lever/mgnMode and one of instId or ccy")
	errAccountSetLeverageAmbiguousScope        = errors.New("okx: set leverage requires instId or ccy (not both)")
	errAccountSetLeverageInvalidMgnModeForCcy  = errors.New("okx: set leverage with ccy requires mgnMode=cross")
	errAccountSetLeverageInvalidPosSideMgnMode = errors.New("okx: set leverage with posSide requires mgnMode=isolated")
	errEmptyAccountSetLeverage                 = errors.New("okx: empty set leverage response")
	errInvalidAccountSetLeverage               = errors.New("okx: invalid set leverage response")
)

func validateAccountSetLeverageAck(ack *AccountSetLeverageAck, req accountSetLeverageRequest) error {
	if ack == nil || ack.Lever == "" || ack.MgnMode == "" {
		return errInvalidAccountSetLeverage
	}
	if ack.Lever != req.Lever || ack.MgnMode != req.MgnMode {
		return errInvalidAccountSetLeverage
	}
	if req.InstId != "" && ack.InstId != req.InstId {
		return errInvalidAccountSetLeverage
	}
	if req.PosSide != "" && ack.PosSide != req.PosSide {
		return errInvalidAccountSetLeverage
	}
	return nil
}

// Do 设置杠杆倍数（POST /api/v5/account/set-leverage）。
func (s *AccountSetLeverageService) Do(ctx context.Context) (*AccountSetLeverageAck, error) {
	if s.req.Lever == "" || s.req.MgnMode == "" || (s.req.InstId == "" && s.req.Ccy == "") {
		return nil, errAccountSetLeverageMissingRequired
	}
	if s.req.InstId != "" && s.req.Ccy != "" {
		return nil, errAccountSetLeverageAmbiguousScope
	}
	if s.req.Ccy != "" && s.req.MgnMode != "cross" {
		return nil, errAccountSetLeverageInvalidMgnModeForCcy
	}
	if s.req.PosSide != "" && s.req.MgnMode != "isolated" {
		return nil, errAccountSetLeverageInvalidPosSideMgnMode
	}

	var data []AccountSetLeverageAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/account/set-leverage", nil, s.req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyAccountSetLeverage
	}
	if err := validateAccountSetLeverageAck(&data[0], s.req); err != nil {
		return nil, err
	}
	return &data[0], nil
}
