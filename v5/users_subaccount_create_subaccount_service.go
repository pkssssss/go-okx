package okx

import (
	"context"
	"errors"
	"net/http"
)

// UsersSubaccountCreateSubaccountService 创建子账户（母账户）。
type UsersSubaccountCreateSubaccountService struct {
	c *Client

	subAcct string
	typ     string
	label   string
	pwd     string
}

// NewUsersSubaccountCreateSubaccountService 创建 UsersSubaccountCreateSubaccountService。
func (c *Client) NewUsersSubaccountCreateSubaccountService() *UsersSubaccountCreateSubaccountService {
	return &UsersSubaccountCreateSubaccountService{c: c}
}

// SubAcct 设置子账户名称（必填）。
func (s *UsersSubaccountCreateSubaccountService) SubAcct(subAcct string) *UsersSubaccountCreateSubaccountService {
	s.subAcct = subAcct
	return s
}

// Type 设置子账户类型（必填）。
func (s *UsersSubaccountCreateSubaccountService) Type(typ string) *UsersSubaccountCreateSubaccountService {
	s.typ = typ
	return s
}

// Label 设置 API Key 备注（必填）。
func (s *UsersSubaccountCreateSubaccountService) Label(label string) *UsersSubaccountCreateSubaccountService {
	s.label = label
	return s
}

// Pwd 设置子账户登录密码（可选；KYB 账户必填）。
func (s *UsersSubaccountCreateSubaccountService) Pwd(pwd string) *UsersSubaccountCreateSubaccountService {
	s.pwd = pwd
	return s
}

var (
	errUsersSubaccountCreateSubaccountMissingSubAcct = errors.New("okx: users create subaccount requires subAcct")
	errUsersSubaccountCreateSubaccountMissingType    = errors.New("okx: users create subaccount requires type")
	errUsersSubaccountCreateSubaccountMissingLabel   = errors.New("okx: users create subaccount requires label")
	errEmptyUsersSubaccountCreateSubaccountResponse  = errors.New("okx: empty users create subaccount response")
)

type usersSubaccountCreateSubaccountRequest struct {
	SubAcct string `json:"subAcct"`
	Type    string `json:"type"`
	Label   string `json:"label"`
	Pwd     string `json:"pwd,omitempty"`
}

// Do 创建子账户（POST /api/v5/users/subaccount/create-subaccount）。
func (s *UsersSubaccountCreateSubaccountService) Do(ctx context.Context) (*UsersSubaccountCreateSubaccountResult, error) {
	if s.subAcct == "" {
		return nil, errUsersSubaccountCreateSubaccountMissingSubAcct
	}
	if s.typ == "" {
		return nil, errUsersSubaccountCreateSubaccountMissingType
	}
	if s.label == "" {
		return nil, errUsersSubaccountCreateSubaccountMissingLabel
	}

	req := usersSubaccountCreateSubaccountRequest{
		SubAcct: s.subAcct,
		Type:    s.typ,
		Label:   s.label,
		Pwd:     s.pwd,
	}

	var data []UsersSubaccountCreateSubaccountResult
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/users/subaccount/create-subaccount", nil, req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/users/subaccount/create-subaccount", requestID, errEmptyUsersSubaccountCreateSubaccountResponse)
	}
	return &data[0], nil
}
