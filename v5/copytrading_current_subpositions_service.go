package okx

import (
	"context"
	"net/http"
)

// CopyTradingCurrentSubpositionsService 获取当前带单（未平仓的带单仓位）。
type CopyTradingCurrentSubpositionsService struct {
	c *Client
	q copyTradingSubpositionsQuery
}

// NewCopyTradingCurrentSubpositionsService 创建 CopyTradingCurrentSubpositionsService。
func (c *Client) NewCopyTradingCurrentSubpositionsService() *CopyTradingCurrentSubpositionsService {
	return &CopyTradingCurrentSubpositionsService{c: c}
}

// InstType 设置产品类型（默认返回所有）。
func (s *CopyTradingCurrentSubpositionsService) InstType(instType string) *CopyTradingCurrentSubpositionsService {
	s.q.instType = instType
	return s
}

// InstId 设置产品 ID 过滤（如 BTC-USDT-SWAP）。
func (s *CopyTradingCurrentSubpositionsService) InstId(instId string) *CopyTradingCurrentSubpositionsService {
	s.q.instId = instId
	return s
}

// After 请求此 id 之前（更旧数据）的分页内容（subPosId）。
func (s *CopyTradingCurrentSubpositionsService) After(after string) *CopyTradingCurrentSubpositionsService {
	s.q.after = after
	return s
}

// Before 请求此 id 之后（更新数据）的分页内容（subPosId）。
func (s *CopyTradingCurrentSubpositionsService) Before(before string) *CopyTradingCurrentSubpositionsService {
	s.q.before = before
	return s
}

// Limit 分页返回数量（最大 500，默认 500）。
func (s *CopyTradingCurrentSubpositionsService) Limit(limit int) *CopyTradingCurrentSubpositionsService {
	s.q.limit = &limit
	return s
}

// Do 获取当前带单（GET /api/v5/copytrading/current-subpositions）。
func (s *CopyTradingCurrentSubpositionsService) Do(ctx context.Context) ([]CopyTradingSubPosition, error) {
	var data []CopyTradingSubPosition
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/copytrading/current-subpositions", s.q.values(), nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
