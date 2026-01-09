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
)

// Do 设置手续费计价方式（POST /api/v5/account/set-fee-type）。
func (s *AccountSetFeeTypeService) Do(ctx context.Context) (*AccountSetFeeTypeAck, error) {
	if s.req.FeeType == "" {
		return nil, errAccountSetFeeTypeMissingFeeType
	}

	var data []AccountSetFeeTypeAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/account/set-fee-type", nil, s.req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyAccountSetFeeType
	}
	return &data[0], nil
}
