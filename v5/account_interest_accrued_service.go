package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// AccountInterestAccrued 表示计息记录。
type AccountInterestAccrued struct {
	Type string `json:"type"`

	Ccy     string `json:"ccy"`
	InstId  string `json:"instId"`
	MgnMode string `json:"mgnMode"`

	Interest     string `json:"interest"`
	InterestRate string `json:"interestRate"`
	Liab         string `json:"liab"`
	TotalLiab    string `json:"totalLiab"`

	InterestFreeLiab string `json:"interestFreeLiab"`

	TS UnixMilli `json:"ts"`
}

// AccountInterestAccruedService 获取计息记录（过去一年）。
type AccountInterestAccruedService struct {
	c *Client

	borrowType string
	ccy        string
	instId     string
	mgnMode    string
	after      string
	before     string
	limit      *int
}

// NewAccountInterestAccruedService 创建 AccountInterestAccruedService。
func (c *Client) NewAccountInterestAccruedService() *AccountInterestAccruedService {
	return &AccountInterestAccruedService{c: c}
}

// Type 设置借币类型（默认 2=市场借币）。
func (s *AccountInterestAccruedService) Type(borrowType string) *AccountInterestAccruedService {
	s.borrowType = borrowType
	return s
}

// Ccy 设置借贷币种（仅适用于市场借币/币币杠杆）。
func (s *AccountInterestAccruedService) Ccy(ccy string) *AccountInterestAccruedService {
	s.ccy = ccy
	return s
}

// InstId 设置产品 ID（仅适用于市场借币）。
func (s *AccountInterestAccruedService) InstId(instId string) *AccountInterestAccruedService {
	s.instId = instId
	return s
}

// MgnMode 设置保证金模式（cross/isolated，仅适用于市场借币）。
func (s *AccountInterestAccruedService) MgnMode(mgnMode string) *AccountInterestAccruedService {
	s.mgnMode = mgnMode
	return s
}

// After 请求此时间戳之前（更旧的数据）的分页内容（Unix 毫秒字符串）。
func (s *AccountInterestAccruedService) After(after string) *AccountInterestAccruedService {
	s.after = after
	return s
}

// Before 请求此时间戳之后（更新的数据）的分页内容（Unix 毫秒字符串）。
func (s *AccountInterestAccruedService) Before(before string) *AccountInterestAccruedService {
	s.before = before
	return s
}

// Limit 设置返回条数（最大 100，默认 100）。
func (s *AccountInterestAccruedService) Limit(limit int) *AccountInterestAccruedService {
	s.limit = &limit
	return s
}

// Do 获取计息记录（GET /api/v5/account/interest-accrued）。
func (s *AccountInterestAccruedService) Do(ctx context.Context) ([]AccountInterestAccrued, error) {
	q := url.Values{}
	if s.borrowType != "" {
		q.Set("type", s.borrowType)
	}
	if s.ccy != "" {
		q.Set("ccy", s.ccy)
	}
	if s.instId != "" {
		q.Set("instId", s.instId)
	}
	if s.mgnMode != "" {
		q.Set("mgnMode", s.mgnMode)
	}
	if s.after != "" {
		q.Set("after", s.after)
	}
	if s.before != "" {
		q.Set("before", s.before)
	}
	if s.limit != nil {
		q.Set("limit", strconv.Itoa(*s.limit))
	}

	if len(q) == 0 {
		q = nil
	}

	var data []AccountInterestAccrued
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/interest-accrued", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
