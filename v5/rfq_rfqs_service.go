package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// RFQRfqsService 获取询价单信息。
type RFQRfqsService struct {
	c *Client

	rfqId   string
	clRfqId string
	state   string
	beginId string
	endId   string
	limit   *int
}

// NewRFQRfqsService 创建 RFQRfqsService。
func (c *Client) NewRFQRfqsService() *RFQRfqsService {
	return &RFQRfqsService{c: c}
}

func (s *RFQRfqsService) RfqId(rfqId string) *RFQRfqsService {
	s.rfqId = rfqId
	return s
}

func (s *RFQRfqsService) ClRfqId(clRfqId string) *RFQRfqsService {
	s.clRfqId = clRfqId
	return s
}

func (s *RFQRfqsService) State(state string) *RFQRfqsService {
	s.state = state
	return s
}

func (s *RFQRfqsService) BeginId(beginId string) *RFQRfqsService {
	s.beginId = beginId
	return s
}

func (s *RFQRfqsService) EndId(endId string) *RFQRfqsService {
	s.endId = endId
	return s
}

func (s *RFQRfqsService) Limit(limit int) *RFQRfqsService {
	s.limit = &limit
	return s
}

// Do 获取询价单信息（GET /api/v5/rfq/rfqs）。
func (s *RFQRfqsService) Do(ctx context.Context) ([]RFQ, error) {
	q := url.Values{}
	if s.rfqId != "" {
		q.Set("rfqId", s.rfqId)
	}
	if s.clRfqId != "" {
		q.Set("clRfqId", s.clRfqId)
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

	var data []RFQ
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/rfq/rfqs", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
