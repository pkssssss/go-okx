package okx

import (
	"context"
	"errors"
	"net/http"
)

type accountSetAutoLoanRequest struct {
	AutoLoan *bool `json:"autoLoan,omitempty"`
}

// AccountSetAutoLoanAck 表示设置自动借币返回项。
type AccountSetAutoLoanAck struct {
	AutoLoan bool `json:"autoLoan"`
}

type accountSetAutoLoanAckRaw struct {
	AutoLoan *bool `json:"autoLoan"`
}

// AccountSetAutoLoanService 设置自动借币。
type AccountSetAutoLoanService struct {
	c *Client
	r accountSetAutoLoanRequest
}

// NewAccountSetAutoLoanService 创建 AccountSetAutoLoanService。
func (c *Client) NewAccountSetAutoLoanService() *AccountSetAutoLoanService {
	return &AccountSetAutoLoanService{c: c}
}

// AutoLoan 设置是否自动借币（可选，省略表示使用 OKX 默认值：true）。
func (s *AccountSetAutoLoanService) AutoLoan(autoLoan bool) *AccountSetAutoLoanService {
	s.r.AutoLoan = &autoLoan
	return s
}

var (
	errEmptyAccountSetAutoLoan   = errors.New("okx: empty set auto loan response")
	errInvalidAccountSetAutoLoan = errors.New("okx: invalid set auto loan response")
)

func validateAccountSetAutoLoanAck(ack *accountSetAutoLoanAckRaw, req accountSetAutoLoanRequest) error {
	if ack == nil || ack.AutoLoan == nil {
		return errInvalidAccountSetAutoLoan
	}
	if req.AutoLoan != nil && *ack.AutoLoan != *req.AutoLoan {
		return errInvalidAccountSetAutoLoan
	}
	return nil
}

// Do 设置自动借币（POST /api/v5/account/set-auto-loan）。
func (s *AccountSetAutoLoanService) Do(ctx context.Context) (*AccountSetAutoLoanAck, error) {
	var data []accountSetAutoLoanAckRaw
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/account/set-auto-loan", nil, s.r, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/account/set-auto-loan", requestID, errEmptyAccountSetAutoLoan)
	}
	if err := validateAccountSetAutoLoanAck(&data[0], s.r); err != nil {
		return nil, newInvalidDataAPIError(http.MethodPost, "/api/v5/account/set-auto-loan", requestID, err)
	}
	return &AccountSetAutoLoanAck{AutoLoan: *data[0].AutoLoan}, nil
}
