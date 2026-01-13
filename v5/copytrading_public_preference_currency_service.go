package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// CopyTradingPublicPreferenceCurrencyService 获取交易员币种偏好（公共）。
type CopyTradingPublicPreferenceCurrencyService struct {
	c *Client

	instType   string
	uniqueCode string
}

// NewCopyTradingPublicPreferenceCurrencyService 创建 CopyTradingPublicPreferenceCurrencyService。
func (c *Client) NewCopyTradingPublicPreferenceCurrencyService() *CopyTradingPublicPreferenceCurrencyService {
	return &CopyTradingPublicPreferenceCurrencyService{c: c}
}

func (s *CopyTradingPublicPreferenceCurrencyService) InstType(instType string) *CopyTradingPublicPreferenceCurrencyService {
	s.instType = instType
	return s
}

func (s *CopyTradingPublicPreferenceCurrencyService) UniqueCode(uniqueCode string) *CopyTradingPublicPreferenceCurrencyService {
	s.uniqueCode = uniqueCode
	return s
}

var errCopyTradingPublicPreferenceCurrencyMissingUniqueCode = errors.New("okx: copytrading public preference currency requires uniqueCode")

// Do 获取交易员币种偏好（GET /api/v5/copytrading/public-preference-currency）。
func (s *CopyTradingPublicPreferenceCurrencyService) Do(ctx context.Context) ([]CopyTradingPreferenceCurrency, error) {
	if s.uniqueCode == "" {
		return nil, errCopyTradingPublicPreferenceCurrencyMissingUniqueCode
	}

	q := url.Values{}
	if s.instType != "" {
		q.Set("instType", s.instType)
	}
	q.Set("uniqueCode", s.uniqueCode)

	var data []CopyTradingPreferenceCurrency
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/copytrading/public-preference-currency", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
