package okx

import (
	"context"
	"net/http"
	"net/url"
)

// AccountMaxWithdrawal 表示账户最大可转余额（交易账户 -> 资金账户）。
type AccountMaxWithdrawal struct {
	Ccy               string `json:"ccy"`
	MaxWd             string `json:"maxWd"`
	MaxWdEx           string `json:"maxWdEx"`
	SpotOffsetMaxWd   string `json:"spotOffsetMaxWd"`
	SpotOffsetMaxWdEx string `json:"spotOffsetMaxWdEx"`
}

// AccountMaxWithdrawalService 查看账户最大可转余额。
type AccountMaxWithdrawalService struct {
	c   *Client
	ccy string
}

// NewAccountMaxWithdrawalService 创建 AccountMaxWithdrawalService。
func (c *Client) NewAccountMaxWithdrawalService() *AccountMaxWithdrawalService {
	return &AccountMaxWithdrawalService{c: c}
}

// Ccy 设置币种过滤（支持多币种，逗号分隔；最多 20 个）。
func (s *AccountMaxWithdrawalService) Ccy(ccy string) *AccountMaxWithdrawalService {
	s.ccy = ccy
	return s
}

// Do 查看账户最大可转余额（GET /api/v5/account/max-withdrawal）。
func (s *AccountMaxWithdrawalService) Do(ctx context.Context) ([]AccountMaxWithdrawal, error) {
	var q url.Values
	if s.ccy != "" {
		q = url.Values{}
		q.Set("ccy", s.ccy)
	}

	var data []AccountMaxWithdrawal
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/max-withdrawal", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
