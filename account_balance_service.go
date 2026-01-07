package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// AccountBalance 表示交易账户余额信息。
// 数值字段按 OKX 返回保持为 string（无损）。
type AccountBalance struct {
	UTime   int64                  `json:"uTime,string"`
	TotalEq string                 `json:"totalEq"`
	AdjEq   string                 `json:"adjEq"`
	AvailEq string                 `json:"availEq"`
	Details []AccountBalanceDetail `json:"details"`
}

// AccountBalanceDetail 表示单币种余额明细。
type AccountBalanceDetail struct {
	Ccy       string `json:"ccy"`
	Eq        string `json:"eq"`
	EqUsd     string `json:"eqUsd"`
	CashBal   string `json:"cashBal"`
	AvailBal  string `json:"availBal"`
	AvailEq   string `json:"availEq"`
	FrozenBal string `json:"frozenBal"`
	Liab      string `json:"liab"`
}

// AccountBalanceService 查看账户余额。
type AccountBalanceService struct {
	c   *Client
	ccy string
}

// NewAccountBalanceService 创建 AccountBalanceService。
func (c *Client) NewAccountBalanceService() *AccountBalanceService {
	return &AccountBalanceService{c: c}
}

// Ccy 设置币种过滤（支持多币种，逗号分隔，如 "BTC,ETH"；最多 20 个）。
func (s *AccountBalanceService) Ccy(ccy string) *AccountBalanceService {
	s.ccy = ccy
	return s
}

var errEmptyAccountBalance = errors.New("okx: empty account balance response")

// Do 查看账户余额（GET /api/v5/account/balance）。
func (s *AccountBalanceService) Do(ctx context.Context) (*AccountBalance, error) {
	var q url.Values
	if s.ccy != "" {
		q = url.Values{}
		q.Set("ccy", s.ccy)
	}

	var data []AccountBalance
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/balance", q, nil, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyAccountBalance
	}
	return &data[0], nil
}
