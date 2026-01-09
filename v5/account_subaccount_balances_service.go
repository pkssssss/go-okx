package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// AccountSubaccountBalancesService 获取子账户交易账户余额（母账户）。
type AccountSubaccountBalancesService struct {
	c       *Client
	subAcct string
}

// NewAccountSubaccountBalancesService 创建 AccountSubaccountBalancesService。
func (c *Client) NewAccountSubaccountBalancesService() *AccountSubaccountBalancesService {
	return &AccountSubaccountBalancesService{c: c}
}

// SubAcct 设置子账户名称（必填）。
func (s *AccountSubaccountBalancesService) SubAcct(subAcct string) *AccountSubaccountBalancesService {
	s.subAcct = subAcct
	return s
}

var (
	errAccountSubaccountBalancesMissingSubAcct = errors.New("okx: subaccount balances requires subAcct")
	errEmptyAccountSubaccountBalances          = errors.New("okx: empty subaccount balances response")
)

// Do 获取子账户交易账户余额（GET /api/v5/account/subaccount/balances）。
func (s *AccountSubaccountBalancesService) Do(ctx context.Context) (*AccountBalance, error) {
	if s.subAcct == "" {
		return nil, errAccountSubaccountBalancesMissingSubAcct
	}

	q := url.Values{}
	q.Set("subAcct", s.subAcct)

	var data []AccountBalance
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/subaccount/balances", q, nil, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyAccountSubaccountBalances
	}
	return &data[0], nil
}
