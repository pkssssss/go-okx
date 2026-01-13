package okx

import (
	"context"
	"errors"
	"net/http"
)

// CopyTradingPublicCurrentSubpositionsService 获取交易员当前带单（公共）。
type CopyTradingPublicCurrentSubpositionsService struct {
	c *Client
	q copyTradingSubpositionsQuery
}

// NewCopyTradingPublicCurrentSubpositionsService 创建 CopyTradingPublicCurrentSubpositionsService。
func (c *Client) NewCopyTradingPublicCurrentSubpositionsService() *CopyTradingPublicCurrentSubpositionsService {
	return &CopyTradingPublicCurrentSubpositionsService{c: c}
}

func (s *CopyTradingPublicCurrentSubpositionsService) InstType(instType string) *CopyTradingPublicCurrentSubpositionsService {
	s.q.instType = instType
	return s
}

func (s *CopyTradingPublicCurrentSubpositionsService) UniqueCode(uniqueCode string) *CopyTradingPublicCurrentSubpositionsService {
	s.q.uniqueCode = uniqueCode
	return s
}

func (s *CopyTradingPublicCurrentSubpositionsService) After(after string) *CopyTradingPublicCurrentSubpositionsService {
	s.q.after = after
	return s
}

func (s *CopyTradingPublicCurrentSubpositionsService) Before(before string) *CopyTradingPublicCurrentSubpositionsService {
	s.q.before = before
	return s
}

func (s *CopyTradingPublicCurrentSubpositionsService) Limit(limit int) *CopyTradingPublicCurrentSubpositionsService {
	s.q.limit = &limit
	return s
}

var errCopyTradingPublicCurrentSubpositionsMissingUniqueCode = errors.New("okx: copytrading public current subpositions requires uniqueCode")

// Do 获取交易员当前带单（GET /api/v5/copytrading/public-current-subpositions）。
func (s *CopyTradingPublicCurrentSubpositionsService) Do(ctx context.Context) ([]CopyTradingSubPosition, error) {
	if s.q.uniqueCode == "" {
		return nil, errCopyTradingPublicCurrentSubpositionsMissingUniqueCode
	}

	var data []CopyTradingSubPosition
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/copytrading/public-current-subpositions", s.q.values(), nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
