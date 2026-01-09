package okx

import (
	"context"
	"errors"
	"net/http"
)

// UsersSubaccountModifyAPIKeyService 重置子账户 API Key（母账户）。
type UsersSubaccountModifyAPIKeyService struct {
	c *Client

	subAcct string
	apiKey  string
	label   string
	perm    string
	ip      *string
}

// NewUsersSubaccountModifyAPIKeyService 创建 UsersSubaccountModifyAPIKeyService。
func (c *Client) NewUsersSubaccountModifyAPIKeyService() *UsersSubaccountModifyAPIKeyService {
	return &UsersSubaccountModifyAPIKeyService{c: c}
}

// SubAcct 设置子账户名称（必填）。
func (s *UsersSubaccountModifyAPIKeyService) SubAcct(subAcct string) *UsersSubaccountModifyAPIKeyService {
	s.subAcct = subAcct
	return s
}

// APIKey 设置子账户 API 公钥（必填）。
func (s *UsersSubaccountModifyAPIKeyService) APIKey(apiKey string) *UsersSubaccountModifyAPIKeyService {
	s.apiKey = apiKey
	return s
}

// Label 重置 API Key 备注（可选；填写则会被重置）。
func (s *UsersSubaccountModifyAPIKeyService) Label(label string) *UsersSubaccountModifyAPIKeyService {
	s.label = label
	return s
}

// Perm 重置 API Key 权限（可选；read_only / trade，多个用半角逗号分隔；填写则会被重置）。
func (s *UsersSubaccountModifyAPIKeyService) Perm(perm string) *UsersSubaccountModifyAPIKeyService {
	s.perm = perm
	return s
}

// IP 重置 API Key 绑定 IP（可选；多个用半角逗号分隔；若传入空字符串则解除 IP 绑定）。
func (s *UsersSubaccountModifyAPIKeyService) IP(ip string) *UsersSubaccountModifyAPIKeyService {
	s.ip = &ip
	return s
}

var (
	errUsersSubaccountModifyAPIKeyMissingSubAcct = errors.New("okx: users modify subaccount apikey requires subAcct")
	errUsersSubaccountModifyAPIKeyMissingAPIKey  = errors.New("okx: users modify subaccount apikey requires apiKey")
	errEmptyUsersSubaccountModifyAPIKeyResponse  = errors.New("okx: empty users modify subaccount apikey response")
)

type usersSubaccountModifyAPIKeyRequest struct {
	SubAcct string  `json:"subAcct"`
	APIKey  string  `json:"apiKey"`
	Label   string  `json:"label,omitempty"`
	Perm    string  `json:"perm,omitempty"`
	IP      *string `json:"ip,omitempty"`
}

// Do 重置子账户 API Key（POST /api/v5/users/subaccount/modify-apikey）。
func (s *UsersSubaccountModifyAPIKeyService) Do(ctx context.Context) (*UsersSubaccountAPIKeyInfo, error) {
	if s.subAcct == "" {
		return nil, errUsersSubaccountModifyAPIKeyMissingSubAcct
	}
	if s.apiKey == "" {
		return nil, errUsersSubaccountModifyAPIKeyMissingAPIKey
	}

	req := usersSubaccountModifyAPIKeyRequest{
		SubAcct: s.subAcct,
		APIKey:  s.apiKey,
		Label:   s.label,
		Perm:    s.perm,
		IP:      s.ip,
	}

	var data []UsersSubaccountAPIKeyInfo
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/users/subaccount/modify-apikey", nil, req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyUsersSubaccountModifyAPIKeyResponse
	}
	return &data[0], nil
}
