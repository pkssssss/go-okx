package okx

import (
	"context"
	"net/http"
	"net/url"
)

// AccountInterestLimitsSurplusLmtDetails 表示剩余可借额度详情（已废弃字段仍可能返回）。
type AccountInterestLimitsSurplusLmtDetails struct {
	AllAcctRemainingQuota string `json:"allAcctRemainingQuota"`
	CurAcctRemainingQuota string `json:"curAcctRemainingQuota"`
	PlatRemainingQuota    string `json:"platRemainingQuota"`
}

// AccountInterestLimitsRecord 表示各币种借币利率与限额详情。
type AccountInterestLimitsRecord struct {
	Ccy      string `json:"ccy"`
	Rate     string `json:"rate"`
	Interest string `json:"interest"`

	LoanQuota  string `json:"loanQuota"`
	UsedLmt    string `json:"usedLmt"`
	SurplusLmt string `json:"surplusLmt"`

	InterestFreeLiab      string `json:"interestFreeLiab"`
	PotentialBorrowingAmt string `json:"potentialBorrowingAmt"`

	SurplusLmtDetails AccountInterestLimitsSurplusLmtDetails `json:"surplusLmtDetails"`

	PosLoan   string `json:"posLoan"`
	AvailLoan string `json:"availLoan"`
	UsedLoan  string `json:"usedLoan"`
	AvgRate   string `json:"avgRate"`
}

// AccountInterestLimits 表示借币利率与限额。
type AccountInterestLimits struct {
	Debt      string `json:"debt"`
	Interest  string `json:"interest"`
	LoanAlloc string `json:"loanAlloc"`

	NextDiscountTime UnixMilli `json:"nextDiscountTime"`
	NextInterestTime UnixMilli `json:"nextInterestTime"`

	Records []AccountInterestLimitsRecord `json:"records"`
}

// AccountInterestLimitsService 获取借币利率与限额。
type AccountInterestLimitsService struct {
	c *Client

	borrowType string
	ccy        string
}

// NewAccountInterestLimitsService 创建 AccountInterestLimitsService。
func (c *Client) NewAccountInterestLimitsService() *AccountInterestLimitsService {
	return &AccountInterestLimitsService{c: c}
}

// Type 设置借币类型（默认 2=市场借币）。
func (s *AccountInterestLimitsService) Type(borrowType string) *AccountInterestLimitsService {
	s.borrowType = borrowType
	return s
}

// Ccy 设置借贷币种过滤（可选）。
func (s *AccountInterestLimitsService) Ccy(ccy string) *AccountInterestLimitsService {
	s.ccy = ccy
	return s
}

// Do 获取借币利率与限额（GET /api/v5/account/interest-limits）。
func (s *AccountInterestLimitsService) Do(ctx context.Context) ([]AccountInterestLimits, error) {
	q := url.Values{}
	if s.borrowType != "" {
		q.Set("type", s.borrowType)
	}
	if s.ccy != "" {
		q.Set("ccy", s.ccy)
	}
	if len(q) == 0 {
		q = nil
	}

	var data []AccountInterestLimits
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/interest-limits", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
