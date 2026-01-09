package okx

import (
	"context"
	"errors"
	"net/http"
)

type accountSetSettleCurrencyRequest struct {
	SettleCcy string `json:"settleCcy"`
}

// AccountSetSettleCurrencyAck 表示设置结算币种返回项。
type AccountSetSettleCurrencyAck struct {
	SettleCcy string `json:"settleCcy"`
}

// AccountSetSettleCurrencyService 设置结算币种（仅适用于 USD 本位合约）。
type AccountSetSettleCurrencyService struct {
	c   *Client
	req accountSetSettleCurrencyRequest
}

// NewAccountSetSettleCurrencyService 创建 AccountSetSettleCurrencyService。
func (c *Client) NewAccountSetSettleCurrencyService() *AccountSetSettleCurrencyService {
	return &AccountSetSettleCurrencyService{c: c}
}

// SettleCcy 设置 USD 本位合约结算币种（必填）。
func (s *AccountSetSettleCurrencyService) SettleCcy(settleCcy string) *AccountSetSettleCurrencyService {
	s.req.SettleCcy = settleCcy
	return s
}

var (
	errAccountSetSettleCurrencyMissingSettleCcy = errors.New("okx: set settle currency requires settleCcy")
	errEmptyAccountSetSettleCurrency            = errors.New("okx: empty set settle currency response")
)

// Do 设置结算币种（POST /api/v5/account/set-settle-currency）。
func (s *AccountSetSettleCurrencyService) Do(ctx context.Context) (*AccountSetSettleCurrencyAck, error) {
	if s.req.SettleCcy == "" {
		return nil, errAccountSetSettleCurrencyMissingSettleCcy
	}

	var data []AccountSetSettleCurrencyAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/account/set-settle-currency", nil, s.req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyAccountSetSettleCurrency
	}
	return &data[0], nil
}
