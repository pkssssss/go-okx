package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// UsersSubaccountListService 查看子账户列表（母账户）。
type UsersSubaccountListService struct {
	c *Client

	enable  *bool
	subAcct string
	after   string
	before  string
	limit   string
}

// NewUsersSubaccountListService 创建 UsersSubaccountListService。
func (c *Client) NewUsersSubaccountListService() *UsersSubaccountListService {
	return &UsersSubaccountListService{c: c}
}

// Enable 设置子账户状态过滤（true: 正常使用，false: 冻结）。
func (s *UsersSubaccountListService) Enable(enable bool) *UsersSubaccountListService {
	s.enable = &enable
	return s
}

// SubAcct 设置子账户名称过滤。
func (s *UsersSubaccountListService) SubAcct(subAcct string) *UsersSubaccountListService {
	s.subAcct = subAcct
	return s
}

// After 查询在此之前的内容（子账户创建时间戳，Unix 毫秒）。
func (s *UsersSubaccountListService) After(after string) *UsersSubaccountListService {
	s.after = after
	return s
}

// Before 查询在此之后的内容（子账户创建时间戳，Unix 毫秒）。
func (s *UsersSubaccountListService) Before(before string) *UsersSubaccountListService {
	s.before = before
	return s
}

// Limit 分页返回数量（最大 100，默认 100）。
func (s *UsersSubaccountListService) Limit(limit string) *UsersSubaccountListService {
	s.limit = limit
	return s
}

// Do 查看子账户列表（GET /api/v5/users/subaccount/list）。
func (s *UsersSubaccountListService) Do(ctx context.Context) ([]UsersSubaccount, error) {
	var q url.Values
	if s.enable != nil || s.subAcct != "" || s.after != "" || s.before != "" || s.limit != "" {
		q = url.Values{}
		if s.enable != nil {
			q.Set("enable", strconv.FormatBool(*s.enable))
		}
		if s.subAcct != "" {
			q.Set("subAcct", s.subAcct)
		}
		if s.after != "" {
			q.Set("after", s.after)
		}
		if s.before != "" {
			q.Set("before", s.before)
		}
		if s.limit != "" {
			q.Set("limit", s.limit)
		}
	}

	var data []UsersSubaccount
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/users/subaccount/list", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
