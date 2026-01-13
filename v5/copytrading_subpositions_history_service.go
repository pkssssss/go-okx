package okx

import (
	"context"
	"net/http"
)

// CopyTradingSubpositionsHistoryService 获取历史带单（最近三个月已平仓）。
type CopyTradingSubpositionsHistoryService struct {
	c *Client
	q copyTradingSubpositionsQuery
}

// NewCopyTradingSubpositionsHistoryService 创建 CopyTradingSubpositionsHistoryService。
func (c *Client) NewCopyTradingSubpositionsHistoryService() *CopyTradingSubpositionsHistoryService {
	return &CopyTradingSubpositionsHistoryService{c: c}
}

// InstType 设置产品类型（默认返回所有）。
func (s *CopyTradingSubpositionsHistoryService) InstType(instType string) *CopyTradingSubpositionsHistoryService {
	s.q.instType = instType
	return s
}

// InstId 设置产品 ID 过滤（如 BTC-USDT-SWAP）。
func (s *CopyTradingSubpositionsHistoryService) InstId(instId string) *CopyTradingSubpositionsHistoryService {
	s.q.instId = instId
	return s
}

// After 请求此 id 之前（更旧数据）的分页内容（subPosId）。
func (s *CopyTradingSubpositionsHistoryService) After(after string) *CopyTradingSubpositionsHistoryService {
	s.q.after = after
	return s
}

// Before 请求此 id 之后（更新数据）的分页内容（subPosId）。
func (s *CopyTradingSubpositionsHistoryService) Before(before string) *CopyTradingSubpositionsHistoryService {
	s.q.before = before
	return s
}

// Limit 分页返回数量（最大 100，默认 100）。
func (s *CopyTradingSubpositionsHistoryService) Limit(limit int) *CopyTradingSubpositionsHistoryService {
	s.q.limit = &limit
	return s
}

// Do 获取历史带单（GET /api/v5/copytrading/subpositions-history）。
func (s *CopyTradingSubpositionsHistoryService) Do(ctx context.Context) ([]CopyTradingSubPosition, error) {
	var data []CopyTradingSubPosition
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/copytrading/subpositions-history", s.q.values(), nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
