package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// AccountSubaccountMaxWithdrawalService 获取子账户最大可转余额（母账户）。
type AccountSubaccountMaxWithdrawalService struct {
	c *Client

	subAcct string
	ccy     string
}

// NewAccountSubaccountMaxWithdrawalService 创建 AccountSubaccountMaxWithdrawalService。
func (c *Client) NewAccountSubaccountMaxWithdrawalService() *AccountSubaccountMaxWithdrawalService {
	return &AccountSubaccountMaxWithdrawalService{c: c}
}

// SubAcct 设置子账户名称（必填）。
func (s *AccountSubaccountMaxWithdrawalService) SubAcct(subAcct string) *AccountSubaccountMaxWithdrawalService {
	s.subAcct = subAcct
	return s
}

// Ccy 设置币种过滤（支持多币种，逗号分隔；最多 20 个）。
func (s *AccountSubaccountMaxWithdrawalService) Ccy(ccy string) *AccountSubaccountMaxWithdrawalService {
	s.ccy = ccy
	return s
}

var errAccountSubaccountMaxWithdrawalMissingSubAcct = errors.New("okx: subaccount max withdrawal requires subAcct")

// Do 获取子账户最大可转余额（GET /api/v5/account/subaccount/max-withdrawal）。
func (s *AccountSubaccountMaxWithdrawalService) Do(ctx context.Context) ([]AccountMaxWithdrawal, error) {
	if s.subAcct == "" {
		return nil, errAccountSubaccountMaxWithdrawalMissingSubAcct
	}

	q := url.Values{}
	q.Set("subAcct", s.subAcct)
	if s.ccy != "" {
		q.Set("ccy", s.ccy)
	}

	var data []AccountMaxWithdrawal
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/subaccount/max-withdrawal", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
