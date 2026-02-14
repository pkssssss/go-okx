package okx

import (
	"context"
	"errors"
	"net/http"
)

type accountSetAutoRepayRequest struct {
	AutoRepay *bool `json:"autoRepay"`
}

// AccountSetAutoRepayAck 表示设置自动还币返回项。
type AccountSetAutoRepayAck struct {
	AutoRepay bool `json:"autoRepay"`
}

type accountSetAutoRepayAckRaw struct {
	AutoRepay *bool `json:"autoRepay"`
}

// AccountSetAutoRepayService 设置自动还币（现货模式）。
type AccountSetAutoRepayService struct {
	c *Client
	r accountSetAutoRepayRequest
}

// NewAccountSetAutoRepayService 创建 AccountSetAutoRepayService。
func (c *Client) NewAccountSetAutoRepayService() *AccountSetAutoRepayService {
	return &AccountSetAutoRepayService{c: c}
}

// AutoRepay 设置是否支持现货模式下自动还币（必填）。
func (s *AccountSetAutoRepayService) AutoRepay(autoRepay bool) *AccountSetAutoRepayService {
	s.r.AutoRepay = &autoRepay
	return s
}

var (
	errAccountSetAutoRepayMissingRequired = errors.New("okx: set auto repay requires autoRepay")
	errEmptyAccountSetAutoRepay           = errors.New("okx: empty set auto repay response")
	errInvalidAccountSetAutoRepay         = errors.New("okx: invalid set auto repay response")
)

func validateAccountSetAutoRepayAck(ack *accountSetAutoRepayAckRaw, req accountSetAutoRepayRequest) error {
	if ack == nil || ack.AutoRepay == nil || req.AutoRepay == nil {
		return errInvalidAccountSetAutoRepay
	}
	if *ack.AutoRepay != *req.AutoRepay {
		return errInvalidAccountSetAutoRepay
	}
	return nil
}

// Do 设置自动还币（POST /api/v5/account/set-auto-repay）。
func (s *AccountSetAutoRepayService) Do(ctx context.Context) (*AccountSetAutoRepayAck, error) {
	if s.r.AutoRepay == nil {
		return nil, errAccountSetAutoRepayMissingRequired
	}

	var data []accountSetAutoRepayAckRaw
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/account/set-auto-repay", nil, s.r, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/account/set-auto-repay", requestID, errEmptyAccountSetAutoRepay)
	}
	if err := validateAccountSetAutoRepayAck(&data[0], s.r); err != nil {
		return nil, err
	}
	return &AccountSetAutoRepayAck{AutoRepay: *data[0].AutoRepay}, nil
}
