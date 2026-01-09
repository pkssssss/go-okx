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

var errEmptyAccountSetAutoLoan = errors.New("okx: empty set auto loan response")

// Do 设置自动借币（POST /api/v5/account/set-auto-loan）。
func (s *AccountSetAutoLoanService) Do(ctx context.Context) (*AccountSetAutoLoanAck, error) {
	var data []AccountSetAutoLoanAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/account/set-auto-loan", nil, s.r, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyAccountSetAutoLoan
	}
	return &data[0], nil
}
