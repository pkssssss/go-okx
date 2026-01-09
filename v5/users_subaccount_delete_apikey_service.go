package okx

import (
	"context"
	"errors"
	"net/http"
)

// UsersSubaccountDeleteAPIKeyService 删除子账户 API Key（母账户）。
type UsersSubaccountDeleteAPIKeyService struct {
	c *Client

	subAcct string
	apiKey  string
}

// NewUsersSubaccountDeleteAPIKeyService 创建 UsersSubaccountDeleteAPIKeyService。
func (c *Client) NewUsersSubaccountDeleteAPIKeyService() *UsersSubaccountDeleteAPIKeyService {
	return &UsersSubaccountDeleteAPIKeyService{c: c}
}

// SubAcct 设置子账户名称（必填）。
func (s *UsersSubaccountDeleteAPIKeyService) SubAcct(subAcct string) *UsersSubaccountDeleteAPIKeyService {
	s.subAcct = subAcct
	return s
}

// APIKey 设置 API 公钥（必填）。
func (s *UsersSubaccountDeleteAPIKeyService) APIKey(apiKey string) *UsersSubaccountDeleteAPIKeyService {
	s.apiKey = apiKey
	return s
}

var (
	errUsersSubaccountDeleteAPIKeyMissingSubAcct = errors.New("okx: users delete subaccount apikey requires subAcct")
	errUsersSubaccountDeleteAPIKeyMissingAPIKey  = errors.New("okx: users delete subaccount apikey requires apiKey")
	errEmptyUsersSubaccountDeleteAPIKeyResponse  = errors.New("okx: empty users delete subaccount apikey response")
)

type usersSubaccountDeleteAPIKeyRequest struct {
	SubAcct string `json:"subAcct"`
	APIKey  string `json:"apiKey"`
}

// Do 删除子账户 API Key（POST /api/v5/users/subaccount/delete-apikey）。
func (s *UsersSubaccountDeleteAPIKeyService) Do(ctx context.Context) (*UsersSubaccountDeleteAPIKeyResult, error) {
	if s.subAcct == "" {
		return nil, errUsersSubaccountDeleteAPIKeyMissingSubAcct
	}
	if s.apiKey == "" {
		return nil, errUsersSubaccountDeleteAPIKeyMissingAPIKey
	}

	req := usersSubaccountDeleteAPIKeyRequest{
		SubAcct: s.subAcct,
		APIKey:  s.apiKey,
	}

	var data []UsersSubaccountDeleteAPIKeyResult
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/users/subaccount/delete-apikey", nil, req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyUsersSubaccountDeleteAPIKeyResponse
	}
	return &data[0], nil
}
