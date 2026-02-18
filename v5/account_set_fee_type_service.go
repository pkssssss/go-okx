package okx

import (
	"context"
	"errors"
	"net/http"
)

type accountSetFeeTypeRequest struct {
	FeeType string `json:"feeType"`
}

// AccountSetFeeTypeAck 表示设置手续费计价方式返回项。
type AccountSetFeeTypeAck struct {
	FeeType string `json:"feeType"`
}

// AccountSetFeeTypeService 设置手续费计价方式（仅对现货生效）。
type AccountSetFeeTypeService struct {
	c   *Client
	req accountSetFeeTypeRequest
}

// NewAccountSetFeeTypeService 创建 AccountSetFeeTypeService。
func (c *Client) NewAccountSetFeeTypeService() *AccountSetFeeTypeService {
	return &AccountSetFeeTypeService{c: c}
}

// FeeType 设置手续费计价方式（必填：0/1）。
func (s *AccountSetFeeTypeService) FeeType(feeType string) *AccountSetFeeTypeService {
	s.req.FeeType = feeType
	return s
}

var (
	errAccountSetFeeTypeMissingFeeType = errors.New("okx: set fee type requires feeType")
	errEmptyAccountSetFeeType          = errors.New("okx: empty set fee type response")
	errInvalidAccountSetFeeType        = errors.New("okx: invalid set fee type response")
)

func validateAccountSetFeeTypeAck(ack *AccountSetFeeTypeAck, req accountSetFeeTypeRequest) error {
	if ack == nil || ack.FeeType == "" || ack.FeeType != req.FeeType {
		return errInvalidAccountSetFeeType
	}
	return nil
}

// Do 设置手续费计价方式（POST /api/v5/account/set-fee-type）。
func (s *AccountSetFeeTypeService) Do(ctx context.Context) (*AccountSetFeeTypeAck, error) {
	if s.req.FeeType == "" {
		return nil, errAccountSetFeeTypeMissingFeeType
	}

	var data []AccountSetFeeTypeAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/account/set-fee-type", nil, s.req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/account/set-fee-type", requestID, errEmptyAccountSetFeeType)
	}
	if err := validateAccountSetFeeTypeAck(&data[0], s.req); err != nil {
		return nil, newInvalidDataAPIError(http.MethodPost, "/api/v5/account/set-fee-type", requestID, err)
	}
	return &data[0], nil
}
