package okx

import (
	"context"
	"errors"
	"net/http"
)

type rfqCancelBatchQuotesRequest struct {
	QuoteIds   []string `json:"quoteIds,omitempty"`
	ClQuoteIds []string `json:"clQuoteIds,omitempty"`
}

// RFQCancelBatchQuotesService 批量取消报价单。
type RFQCancelBatchQuotesService struct {
	c *Client

	quoteIds   []string
	clQuoteIds []string
}

// NewRFQCancelBatchQuotesService 创建 RFQCancelBatchQuotesService。
func (c *Client) NewRFQCancelBatchQuotesService() *RFQCancelBatchQuotesService {
	return &RFQCancelBatchQuotesService{c: c}
}

// QuoteIds 设置报价单 ID 列表（可选，最多 100 个）。
func (s *RFQCancelBatchQuotesService) QuoteIds(quoteIds []string) *RFQCancelBatchQuotesService {
	s.quoteIds = quoteIds
	return s
}

// ClQuoteIds 设置报价单自定义 ID 列表（可选，最多 100 个）。
func (s *RFQCancelBatchQuotesService) ClQuoteIds(clQuoteIds []string) *RFQCancelBatchQuotesService {
	s.clQuoteIds = clQuoteIds
	return s
}

var (
	errRFQCancelBatchQuotesMissingIds = errors.New("okx: rfq cancel batch quotes requires quoteIds or clQuoteIds")
	errRFQCancelBatchQuotesTooManyIds = errors.New("okx: rfq cancel batch quotes max 100 ids")
)

// Do 批量取消报价单（POST /api/v5/rfq/cancel-batch-quotes）。
func (s *RFQCancelBatchQuotesService) Do(ctx context.Context) ([]RFQCancelQuoteAck, error) {
	if len(s.quoteIds) == 0 && len(s.clQuoteIds) == 0 {
		return nil, errRFQCancelBatchQuotesMissingIds
	}
	if len(s.quoteIds) > rfqMaxCancelBatch || len(s.clQuoteIds) > rfqMaxCancelBatch {
		return nil, errRFQCancelBatchQuotesTooManyIds
	}

	req := rfqCancelBatchQuotesRequest{
		QuoteIds:   s.quoteIds,
		ClQuoteIds: s.clQuoteIds,
	}

	var data []RFQCancelQuoteAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/rfq/cancel-batch-quotes", nil, req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if err := rfqCheckCancelBatchQuotes(http.MethodPost, "/api/v5/rfq/cancel-batch-quotes", requestID, data); err != nil {
		return data, err
	}
	return data, nil
}
