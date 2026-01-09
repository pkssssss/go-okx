package okx

import (
	"context"
	"errors"
	"net/http"
)

type accountSetAutoEarnRequest struct {
	EarnType string `json:"earnType,omitempty"`
	Ccy      string `json:"ccy"`
	Action   string `json:"action"`
	Apr      string `json:"apr,omitempty"`
}

// AccountSetAutoEarnAck 表示设置自动赚币返回项。
type AccountSetAutoEarnAck struct {
	EarnType string `json:"earnType"`
	Ccy      string `json:"ccy"`
	Action   string `json:"action"`
	Apr      string `json:"apr"`
}

// AccountSetAutoEarnService 设置自动赚币（开启/关闭）。
type AccountSetAutoEarnService struct {
	c   *Client
	req accountSetAutoEarnRequest
}

// NewAccountSetAutoEarnService 创建 AccountSetAutoEarnService。
func (c *Client) NewAccountSetAutoEarnService() *AccountSetAutoEarnService {
	return &AccountSetAutoEarnService{c: c}
}

// EarnType 设置自动赚币类型（可选：0 自动赚币；1 自动赚币（USDG 赚币））。
func (s *AccountSetAutoEarnService) EarnType(earnType string) *AccountSetAutoEarnService {
	s.req.EarnType = earnType
	return s
}

// Ccy 设置币种（必填）。
func (s *AccountSetAutoEarnService) Ccy(ccy string) *AccountSetAutoEarnService {
	s.req.Ccy = ccy
	return s
}

// Action 设置自动赚币操作类型（必填：turn_on/turn_off）。
func (s *AccountSetAutoEarnService) Action(action string) *AccountSetAutoEarnService {
	s.req.Action = action
	return s
}

// Apr 设置最低年化收益率（可选；已弃用，仅为兼容保留）。
func (s *AccountSetAutoEarnService) Apr(apr string) *AccountSetAutoEarnService {
	s.req.Apr = apr
	return s
}

var (
	errAccountSetAutoEarnMissingRequired = errors.New("okx: set auto earn requires ccy and action")
	errEmptyAccountSetAutoEarn           = errors.New("okx: empty set auto earn response")
)

// Do 设置自动赚币（POST /api/v5/account/set-auto-earn）。
func (s *AccountSetAutoEarnService) Do(ctx context.Context) (*AccountSetAutoEarnAck, error) {
	if s.req.Ccy == "" || s.req.Action == "" {
		return nil, errAccountSetAutoEarnMissingRequired
	}

	var data []AccountSetAutoEarnAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/account/set-auto-earn", nil, s.req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyAccountSetAutoEarn
	}
	return &data[0], nil
}
