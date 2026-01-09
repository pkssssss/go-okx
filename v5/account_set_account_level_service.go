package okx

import (
	"context"
	"errors"
	"net/http"
)

type accountSetAccountLevelRequest struct {
	AcctLv string `json:"acctLv"`
}

// AccountSetAccountLevelAck 表示设置账户模式返回项。
type AccountSetAccountLevelAck struct {
	AcctLv string `json:"acctLv"`
}

// AccountSetAccountLevelService 设置账户模式。
type AccountSetAccountLevelService struct {
	c *Client
	r accountSetAccountLevelRequest
}

// NewAccountSetAccountLevelService 创建 AccountSetAccountLevelService。
func (c *Client) NewAccountSetAccountLevelService() *AccountSetAccountLevelService {
	return &AccountSetAccountLevelService{c: c}
}

// AcctLv 设置账户模式（必填）。
func (s *AccountSetAccountLevelService) AcctLv(acctLv string) *AccountSetAccountLevelService {
	s.r.AcctLv = acctLv
	return s
}

var (
	errAccountSetAccountLevelMissingAcctLv = errors.New("okx: set account level requires acctLv")
	errEmptyAccountSetAccountLevel         = errors.New("okx: empty set account level response")
)

// Do 设置账户模式（POST /api/v5/account/set-account-level）。
func (s *AccountSetAccountLevelService) Do(ctx context.Context) (*AccountSetAccountLevelAck, error) {
	if s.r.AcctLv == "" {
		return nil, errAccountSetAccountLevelMissingAcctLv
	}

	var data []AccountSetAccountLevelAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/account/set-account-level", nil, s.r, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyAccountSetAccountLevel
	}
	return &data[0], nil
}
