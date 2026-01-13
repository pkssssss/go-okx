package okx

import (
	"context"
	"net/http"
	"net/url"
)

// CopyTradingInstrumentsService 获取带单产品。
type CopyTradingInstrumentsService struct {
	c *Client

	instType string
}

// NewCopyTradingInstrumentsService 创建 CopyTradingInstrumentsService。
func (c *Client) NewCopyTradingInstrumentsService() *CopyTradingInstrumentsService {
	return &CopyTradingInstrumentsService{c: c}
}

// InstType 设置产品类型（默认 SWAP）。
func (s *CopyTradingInstrumentsService) InstType(instType string) *CopyTradingInstrumentsService {
	s.instType = instType
	return s
}

// Do 获取带单产品（GET /api/v5/copytrading/instruments）。
func (s *CopyTradingInstrumentsService) Do(ctx context.Context) ([]CopyTradingInstrument, error) {
	var q url.Values
	if s.instType != "" {
		q = url.Values{}
		q.Set("instType", s.instType)
	}

	var data []CopyTradingInstrument
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/copytrading/instruments", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
