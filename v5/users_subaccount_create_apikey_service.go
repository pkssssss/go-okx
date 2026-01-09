package okx

import (
	"context"
	"errors"
	"net/http"
)

// UsersSubaccountCreateAPIKeyService 创建子账户 API Key（母账户）。
type UsersSubaccountCreateAPIKeyService struct {
	c *Client

	subAcct    string
	label      string
	passphrase string
	perm       string
	ip         string
}

// NewUsersSubaccountCreateAPIKeyService 创建 UsersSubaccountCreateAPIKeyService。
func (c *Client) NewUsersSubaccountCreateAPIKeyService() *UsersSubaccountCreateAPIKeyService {
	return &UsersSubaccountCreateAPIKeyService{c: c}
}

// SubAcct 设置子账户名称（必填）。
func (s *UsersSubaccountCreateAPIKeyService) SubAcct(subAcct string) *UsersSubaccountCreateAPIKeyService {
	s.subAcct = subAcct
	return s
}

// Label 设置 API Key 备注（必填）。
func (s *UsersSubaccountCreateAPIKeyService) Label(label string) *UsersSubaccountCreateAPIKeyService {
	s.label = label
	return s
}

// Passphrase 设置 API Key 密码（必填）。
func (s *UsersSubaccountCreateAPIKeyService) Passphrase(passphrase string) *UsersSubaccountCreateAPIKeyService {
	s.passphrase = passphrase
	return s
}

// Perm 设置 API Key 权限（可选，read_only / trade）。
func (s *UsersSubaccountCreateAPIKeyService) Perm(perm string) *UsersSubaccountCreateAPIKeyService {
	s.perm = perm
	return s
}

// IP 设置绑定 IP（可选，多个用半角逗号分隔）。
func (s *UsersSubaccountCreateAPIKeyService) IP(ip string) *UsersSubaccountCreateAPIKeyService {
	s.ip = ip
	return s
}

var (
	errUsersSubaccountCreateAPIKeyMissingSubAcct    = errors.New("okx: users create subaccount apikey requires subAcct")
	errUsersSubaccountCreateAPIKeyMissingLabel      = errors.New("okx: users create subaccount apikey requires label")
	errUsersSubaccountCreateAPIKeyMissingPassphrase = errors.New("okx: users create subaccount apikey requires passphrase")
	errEmptyUsersSubaccountCreateAPIKeyResponse     = errors.New("okx: empty users create subaccount apikey response")
)

type usersSubaccountCreateAPIKeyRequest struct {
	SubAcct    string `json:"subAcct"`
	Label      string `json:"label"`
	Passphrase string `json:"passphrase"`
	Perm       string `json:"perm,omitempty"`
	IP         string `json:"ip,omitempty"`
}

// Do 创建子账户 API Key（POST /api/v5/users/subaccount/apikey）。
func (s *UsersSubaccountCreateAPIKeyService) Do(ctx context.Context) (*UsersSubaccountAPIKeyCreateResult, error) {
	if s.subAcct == "" {
		return nil, errUsersSubaccountCreateAPIKeyMissingSubAcct
	}
	if s.label == "" {
		return nil, errUsersSubaccountCreateAPIKeyMissingLabel
	}
	if s.passphrase == "" {
		return nil, errUsersSubaccountCreateAPIKeyMissingPassphrase
	}

	req := usersSubaccountCreateAPIKeyRequest{
		SubAcct:    s.subAcct,
		Label:      s.label,
		Passphrase: s.passphrase,
		Perm:       s.perm,
		IP:         s.ip,
	}

	var data []UsersSubaccountAPIKeyCreateResult
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/users/subaccount/apikey", nil, req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyUsersSubaccountCreateAPIKeyResponse
	}
	return &data[0], nil
}
