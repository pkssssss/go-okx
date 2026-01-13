package okx

import (
	"context"
	"errors"
	"net/http"
)

// CopyTradingPublicSubpositionsHistoryService 获取交易员历史带单（公共）。
type CopyTradingPublicSubpositionsHistoryService struct {
	c *Client
	q copyTradingSubpositionsQuery
}

// NewCopyTradingPublicSubpositionsHistoryService 创建 CopyTradingPublicSubpositionsHistoryService。
func (c *Client) NewCopyTradingPublicSubpositionsHistoryService() *CopyTradingPublicSubpositionsHistoryService {
	return &CopyTradingPublicSubpositionsHistoryService{c: c}
}

func (s *CopyTradingPublicSubpositionsHistoryService) InstType(instType string) *CopyTradingPublicSubpositionsHistoryService {
	s.q.instType = instType
	return s
}

func (s *CopyTradingPublicSubpositionsHistoryService) UniqueCode(uniqueCode string) *CopyTradingPublicSubpositionsHistoryService {
	s.q.uniqueCode = uniqueCode
	return s
}

func (s *CopyTradingPublicSubpositionsHistoryService) After(after string) *CopyTradingPublicSubpositionsHistoryService {
	s.q.after = after
	return s
}

func (s *CopyTradingPublicSubpositionsHistoryService) Before(before string) *CopyTradingPublicSubpositionsHistoryService {
	s.q.before = before
	return s
}

func (s *CopyTradingPublicSubpositionsHistoryService) Limit(limit int) *CopyTradingPublicSubpositionsHistoryService {
	s.q.limit = &limit
	return s
}

var errCopyTradingPublicSubpositionsHistoryMissingUniqueCode = errors.New("okx: copytrading public subpositions history requires uniqueCode")

// Do 获取交易员历史带单（GET /api/v5/copytrading/public-subpositions-history）。
func (s *CopyTradingPublicSubpositionsHistoryService) Do(ctx context.Context) ([]CopyTradingSubPosition, error) {
	if s.q.uniqueCode == "" {
		return nil, errCopyTradingPublicSubpositionsHistoryMissingUniqueCode
	}

	var data []CopyTradingSubPosition
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/copytrading/public-subpositions-history", s.q.values(), nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
