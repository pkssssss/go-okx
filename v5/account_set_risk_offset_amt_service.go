package okx

import (
	"context"
	"errors"
	"net/http"
)

type accountSetRiskOffsetAmtRequest struct {
	Ccy            string `json:"ccy"`
	ClSpotInUseAmt string `json:"clSpotInUseAmt"`
}

// AccountSetRiskOffsetAmtAck 表示设置现货对冲占用返回项。
type AccountSetRiskOffsetAmtAck struct {
	Ccy            string `json:"ccy"`
	ClSpotInUseAmt string `json:"clSpotInUseAmt"`
}

// AccountSetRiskOffsetAmtService 设置现货对冲占用（仅适用于组合保证金模式）。
type AccountSetRiskOffsetAmtService struct {
	c   *Client
	req accountSetRiskOffsetAmtRequest
}

// NewAccountSetRiskOffsetAmtService 创建 AccountSetRiskOffsetAmtService。
func (c *Client) NewAccountSetRiskOffsetAmtService() *AccountSetRiskOffsetAmtService {
	return &AccountSetRiskOffsetAmtService{c: c}
}

// Ccy 设置币种（必填）。
func (s *AccountSetRiskOffsetAmtService) Ccy(ccy string) *AccountSetRiskOffsetAmtService {
	s.req.Ccy = ccy
	return s
}

// ClSpotInUseAmt 设置用户自定义现货对冲数量（必填）。
func (s *AccountSetRiskOffsetAmtService) ClSpotInUseAmt(amt string) *AccountSetRiskOffsetAmtService {
	s.req.ClSpotInUseAmt = amt
	return s
}

var (
	errAccountSetRiskOffsetAmtMissingRequired = errors.New("okx: set risk offset amt requires ccy and clSpotInUseAmt")
	errEmptyAccountSetRiskOffsetAmt           = errors.New("okx: empty set risk offset amt response")
	errInvalidAccountSetRiskOffsetAmt         = errors.New("okx: invalid set risk offset amt response")
)

func validateAccountSetRiskOffsetAmtAck(ack *AccountSetRiskOffsetAmtAck, req accountSetRiskOffsetAmtRequest) error {
	if ack == nil || ack.Ccy == "" || ack.ClSpotInUseAmt == "" {
		return errInvalidAccountSetRiskOffsetAmt
	}
	if ack.Ccy != req.Ccy || ack.ClSpotInUseAmt != req.ClSpotInUseAmt {
		return errInvalidAccountSetRiskOffsetAmt
	}
	return nil
}

// Do 设置现货对冲占用（POST /api/v5/account/set-riskOffset-amt）。
func (s *AccountSetRiskOffsetAmtService) Do(ctx context.Context) (*AccountSetRiskOffsetAmtAck, error) {
	if s.req.Ccy == "" || s.req.ClSpotInUseAmt == "" {
		return nil, errAccountSetRiskOffsetAmtMissingRequired
	}

	var data []AccountSetRiskOffsetAmtAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/account/set-riskOffset-amt", nil, s.req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/account/set-riskOffset-amt", requestID, errEmptyAccountSetRiskOffsetAmt)
	}
	if err := validateAccountSetRiskOffsetAmtAck(&data[0], s.req); err != nil {
		return nil, newInvalidDataAPIError(http.MethodPost, "/api/v5/account/set-riskOffset-amt", requestID, err)
	}
	return &data[0], nil
}
