package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// RFQQuotesService 获取报价单信息。
type RFQQuotesService struct {
	c *Client

	rfqId     string
	clRfqId   string
	quoteId   string
	clQuoteId string
	state     string
	beginId   string
	endId     string
	limit     *int
}

// NewRFQQuotesService 创建 RFQQuotesService。
func (c *Client) NewRFQQuotesService() *RFQQuotesService {
	return &RFQQuotesService{c: c}
}

func (s *RFQQuotesService) RfqId(rfqId string) *RFQQuotesService {
	s.rfqId = rfqId
	return s
}

func (s *RFQQuotesService) ClRfqId(clRfqId string) *RFQQuotesService {
	s.clRfqId = clRfqId
	return s
}

func (s *RFQQuotesService) QuoteId(quoteId string) *RFQQuotesService {
	s.quoteId = quoteId
	return s
}

func (s *RFQQuotesService) ClQuoteId(clQuoteId string) *RFQQuotesService {
	s.clQuoteId = clQuoteId
	return s
}

func (s *RFQQuotesService) State(state string) *RFQQuotesService {
	s.state = state
	return s
}

func (s *RFQQuotesService) BeginId(beginId string) *RFQQuotesService {
	s.beginId = beginId
	return s
}

func (s *RFQQuotesService) EndId(endId string) *RFQQuotesService {
	s.endId = endId
	return s
}

func (s *RFQQuotesService) Limit(limit int) *RFQQuotesService {
	s.limit = &limit
	return s
}

// Do 获取报价单信息（GET /api/v5/rfq/quotes）。
func (s *RFQQuotesService) Do(ctx context.Context) ([]Quote, error) {
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
	if s.state != "" {
		q.Set("state", s.state)
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

	var data []Quote
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/rfq/quotes", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
