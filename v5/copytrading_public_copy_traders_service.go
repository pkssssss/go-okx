package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// CopyTradingPublicCopyTradersService 获取跟单人信息（公共）。
type CopyTradingPublicCopyTradersService struct {
	c *Client

	instType   string
	uniqueCode string
	limit      *int
}

// NewCopyTradingPublicCopyTradersService 创建 CopyTradingPublicCopyTradersService。
func (c *Client) NewCopyTradingPublicCopyTradersService() *CopyTradingPublicCopyTradersService {
	return &CopyTradingPublicCopyTradersService{c: c}
}

func (s *CopyTradingPublicCopyTradersService) InstType(instType string) *CopyTradingPublicCopyTradersService {
	s.instType = instType
	return s
}

func (s *CopyTradingPublicCopyTradersService) UniqueCode(uniqueCode string) *CopyTradingPublicCopyTradersService {
	s.uniqueCode = uniqueCode
	return s
}

func (s *CopyTradingPublicCopyTradersService) Limit(limit int) *CopyTradingPublicCopyTradersService {
	s.limit = &limit
	return s
}

var (
	errCopyTradingPublicCopyTradersMissingUniqueCode = errors.New("okx: copytrading public copy traders requires uniqueCode")
	errEmptyCopyTradingPublicCopyTradersResponse     = errors.New("okx: empty copytrading public copy traders response")
)

// Do 获取跟单人信息（GET /api/v5/copytrading/public-copy-traders）。
func (s *CopyTradingPublicCopyTradersService) Do(ctx context.Context) (*CopyTradingPublicCopyTraders, error) {
	if s.uniqueCode == "" {
		return nil, errCopyTradingPublicCopyTradersMissingUniqueCode
	}

	q := url.Values{}
	if s.instType != "" {
		q.Set("instType", s.instType)
	}
	q.Set("uniqueCode", s.uniqueCode)
	if s.limit != nil {
		q.Set("limit", strconv.Itoa(*s.limit))
	}

	var data []CopyTradingPublicCopyTraders
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/copytrading/public-copy-traders", q, nil, false, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyCopyTradingPublicCopyTradersResponse
	}
	return &data[0], nil
}
