package okx

import (
	"context"
	"net/http"
	"net/url"
)

// UsersEntrustSubaccountListService 查看被托管的子账户列表。
type UsersEntrustSubaccountListService struct {
	c *Client

	subAcct string
}

// NewUsersEntrustSubaccountListService 创建 UsersEntrustSubaccountListService。
func (c *Client) NewUsersEntrustSubaccountListService() *UsersEntrustSubaccountListService {
	return &UsersEntrustSubaccountListService{c: c}
}

// SubAcct 设置子账户名称过滤（可选）。
func (s *UsersEntrustSubaccountListService) SubAcct(subAcct string) *UsersEntrustSubaccountListService {
	s.subAcct = subAcct
	return s
}

// Do 查看被托管的子账户列表（GET /api/v5/users/entrust-subaccount-list）。
func (s *UsersEntrustSubaccountListService) Do(ctx context.Context) ([]UsersEntrustSubaccount, error) {
	var q url.Values
	if s.subAcct != "" {
		q = url.Values{}
		q.Set("subAcct", s.subAcct)
	}

	var data []UsersEntrustSubaccount
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/users/entrust-subaccount-list", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
