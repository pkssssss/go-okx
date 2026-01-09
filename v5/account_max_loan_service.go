package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// AccountMaxLoan 表示交易产品最大可借信息。
type AccountMaxLoan struct {
	InstId  string `json:"instId"`
	MgnMode string `json:"mgnMode"`
	MgnCcy  string `json:"mgnCcy"`
	MaxLoan string `json:"maxLoan"`
	Ccy     string `json:"ccy"`
	Side    string `json:"side"`
}

// AccountMaxLoanService 获取交易产品最大可借。
type AccountMaxLoanService struct {
	c *Client

	mgnMode       string
	instId        string
	ccy           string
	mgnCcy        string
	tradeQuoteCcy string
}

// NewAccountMaxLoanService 创建 AccountMaxLoanService。
func (c *Client) NewAccountMaxLoanService() *AccountMaxLoanService {
	return &AccountMaxLoanService{c: c}
}

// MgnMode 设置仓位类型（必填：isolated/cross）。
func (s *AccountMaxLoanService) MgnMode(mgnMode string) *AccountMaxLoanService {
	s.mgnMode = mgnMode
	return s
}

// InstId 设置产品 ID（支持多产品 ID 查询，逗号分隔；最多 5 个）。
func (s *AccountMaxLoanService) InstId(instId string) *AccountMaxLoanService {
	s.instId = instId
	return s
}

// Ccy 设置币种（仅适用于现货模式下手动借币币种最大可借）。
func (s *AccountMaxLoanService) Ccy(ccy string) *AccountMaxLoanService {
	s.ccy = ccy
	return s
}

// MgnCcy 设置保证金币种（适用于逐仓杠杆及合约模式下的全仓杠杆）。
func (s *AccountMaxLoanService) MgnCcy(mgnCcy string) *AccountMaxLoanService {
	s.mgnCcy = mgnCcy
	return s
}

// TradeQuoteCcy 设置用于交易的计价币种（仅适用于币币）。
func (s *AccountMaxLoanService) TradeQuoteCcy(tradeQuoteCcy string) *AccountMaxLoanService {
	s.tradeQuoteCcy = tradeQuoteCcy
	return s
}

var errAccountMaxLoanMissingMgnMode = errors.New("okx: max loan requires mgnMode")

// Do 获取交易产品最大可借（GET /api/v5/account/max-loan）。
func (s *AccountMaxLoanService) Do(ctx context.Context) ([]AccountMaxLoan, error) {
	if s.mgnMode == "" {
		return nil, errAccountMaxLoanMissingMgnMode
	}

	q := url.Values{}
	q.Set("mgnMode", s.mgnMode)
	if s.instId != "" {
		q.Set("instId", s.instId)
	}
	if s.ccy != "" {
		q.Set("ccy", s.ccy)
	}
	if s.mgnCcy != "" {
		q.Set("mgnCcy", s.mgnCcy)
	}
	if s.tradeQuoteCcy != "" {
		q.Set("tradeQuoteCcy", s.tradeQuoteCcy)
	}

	var data []AccountMaxLoan
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/max-loan", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
