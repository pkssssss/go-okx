package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// UsersSubaccountAPIKeysService 查询子账户 API Key（母账户）。
type UsersSubaccountAPIKeysService struct {
	c *Client

	subAcct string
	apiKey  string
}

// NewUsersSubaccountAPIKeysService 创建 UsersSubaccountAPIKeysService。
func (c *Client) NewUsersSubaccountAPIKeysService() *UsersSubaccountAPIKeysService {
	return &UsersSubaccountAPIKeysService{c: c}
}

// SubAcct 设置子账户名称（必填）。
func (s *UsersSubaccountAPIKeysService) SubAcct(subAcct string) *UsersSubaccountAPIKeysService {
	s.subAcct = subAcct
	return s
}

// APIKey 设置 API 公钥过滤（可选）。
func (s *UsersSubaccountAPIKeysService) APIKey(apiKey string) *UsersSubaccountAPIKeysService {
	s.apiKey = apiKey
	return s
}

var errUsersSubaccountAPIKeysMissingSubAcct = errors.New("okx: users subaccount apikeys requires subAcct")

// Do 查询子账户 API Key（GET /api/v5/users/subaccount/apikey）。
func (s *UsersSubaccountAPIKeysService) Do(ctx context.Context) ([]UsersSubaccountAPIKeyInfo, error) {
	if s.subAcct == "" {
		return nil, errUsersSubaccountAPIKeysMissingSubAcct
	}

	q := url.Values{}
	q.Set("subAcct", s.subAcct)
	if s.apiKey != "" {
		q.Set("apiKey", s.apiKey)
	}

	var data []UsersSubaccountAPIKeyInfo
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/users/subaccount/apikey", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
