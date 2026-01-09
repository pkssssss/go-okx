package okx

import (
	"context"
	"net/http"
	"net/url"
)

// AccountInterestRate 表示用户当前市场借币利率（每小时）。
type AccountInterestRate struct {
	Ccy          string `json:"ccy"`
	InterestRate string `json:"interestRate"`
}

// AccountInterestRateService 获取用户当前市场借币利率。
type AccountInterestRateService struct {
	c   *Client
	ccy string
}

// NewAccountInterestRateService 创建 AccountInterestRateService。
func (c *Client) NewAccountInterestRateService() *AccountInterestRateService {
	return &AccountInterestRateService{c: c}
}

// Ccy 设置币种过滤（可选）。
func (s *AccountInterestRateService) Ccy(ccy string) *AccountInterestRateService {
	s.ccy = ccy
	return s
}

// Do 获取用户当前市场借币利率（GET /api/v5/account/interest-rate）。
func (s *AccountInterestRateService) Do(ctx context.Context) ([]AccountInterestRate, error) {
	var q url.Values
	if s.ccy != "" {
		q = url.Values{}
		q.Set("ccy", s.ccy)
	}

	var data []AccountInterestRate
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/interest-rate", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
