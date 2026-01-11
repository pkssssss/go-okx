package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// RFQTradesService 获取大宗交易信息。
type RFQTradesService struct {
	c *Client

	rfqId     string
	clRfqId   string
	quoteId   string
	clQuoteId string
	blockTdId string
	beginId   string
	endId     string
	limit     *int
}

// NewRFQTradesService 创建 RFQTradesService。
func (c *Client) NewRFQTradesService() *RFQTradesService {
	return &RFQTradesService{c: c}
}

func (s *RFQTradesService) RfqId(rfqId string) *RFQTradesService {
	s.rfqId = rfqId
	return s
}

func (s *RFQTradesService) ClRfqId(clRfqId string) *RFQTradesService {
	s.clRfqId = clRfqId
	return s
}

func (s *RFQTradesService) QuoteId(quoteId string) *RFQTradesService {
	s.quoteId = quoteId
	return s
}

func (s *RFQTradesService) ClQuoteId(clQuoteId string) *RFQTradesService {
	s.clQuoteId = clQuoteId
	return s
}

func (s *RFQTradesService) BlockTdId(blockTdId string) *RFQTradesService {
	s.blockTdId = blockTdId
	return s
}

func (s *RFQTradesService) BeginId(beginId string) *RFQTradesService {
	s.beginId = beginId
	return s
}

func (s *RFQTradesService) EndId(endId string) *RFQTradesService {
	s.endId = endId
	return s
}

func (s *RFQTradesService) Limit(limit int) *RFQTradesService {
	s.limit = &limit
	return s
}

// Do 获取大宗交易信息（GET /api/v5/rfq/trades）。
func (s *RFQTradesService) Do(ctx context.Context) ([]StrucBlockTrade, error) {
	q := url.Values{}
	if s.rfqId != "" {
		q.Set("rfqId", s.rfqId)
	}
	if s.clRfqId != "" {
		q.Set("clRfqId", s.clRfqId)
	}
	if s.quoteId != "" {
		q.Set("quoteId", s.quoteId)
	}
	if s.clQuoteId != "" {
		q.Set("clQuoteId", s.clQuoteId)
	}
	if s.blockTdId != "" {
		q.Set("blockTdId", s.blockTdId)
	}
	if s.beginId != "" {
		q.Set("beginId", s.beginId)
	}
	if s.endId != "" {
		q.Set("endId", s.endId)
	}
	if s.limit != nil {
		q.Set("limit", strconv.Itoa(*s.limit))
	}

	var data []StrucBlockTrade
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/rfq/trades", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
