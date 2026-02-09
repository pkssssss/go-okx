package okx

import (
	"context"
	"errors"
	"net/http"
)

// UsersSubaccountSetTransferOutService 设置子账户主动转出权限（母账户）。
type UsersSubaccountSetTransferOutService struct {
	c *Client

	subAcct     string
	canTransOut *bool
}

// NewUsersSubaccountSetTransferOutService 创建 UsersSubaccountSetTransferOutService。
func (c *Client) NewUsersSubaccountSetTransferOutService() *UsersSubaccountSetTransferOutService {
	return &UsersSubaccountSetTransferOutService{c: c}
}

// SubAcct 设置子账户名称（必填；支持多个，半角逗号分隔，最多 20 个）。
func (s *UsersSubaccountSetTransferOutService) SubAcct(subAcct string) *UsersSubaccountSetTransferOutService {
	s.subAcct = subAcct
	return s
}

// CanTransOut 设置是否可以主动转出（可选；默认 true）。
func (s *UsersSubaccountSetTransferOutService) CanTransOut(enable bool) *UsersSubaccountSetTransferOutService {
	s.canTransOut = &enable
	return s
}

var (
	errUsersSubaccountSetTransferOutMissingSubAcct = errors.New("okx: users set transfer out requires subAcct")
	errEmptyUsersSubaccountSetTransferOutResponse  = errors.New("okx: empty users set transfer out response")
)

type usersSubaccountSetTransferOutRequest struct {
	SubAcct     string `json:"subAcct"`
	CanTransOut *bool  `json:"canTransOut,omitempty"`
}

// Do 设置子账户主动转出权限（POST /api/v5/users/subaccount/set-transfer-out）。
func (s *UsersSubaccountSetTransferOutService) Do(ctx context.Context) ([]UsersSubaccountTransferOutPermission, error) {
	if s.subAcct == "" {
		return nil, errUsersSubaccountSetTransferOutMissingSubAcct
	}

	req := usersSubaccountSetTransferOutRequest{
		SubAcct:     s.subAcct,
		CanTransOut: s.canTransOut,
	}

	var data []UsersSubaccountTransferOutPermission
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/users/subaccount/set-transfer-out", nil, req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyUsersSubaccountSetTransferOutResponse
	}
	return data, nil
}
