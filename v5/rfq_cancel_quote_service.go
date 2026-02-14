package okx

import (
	"context"
	"errors"
	"net/http"
)

type rfqCancelQuoteRequest struct {
	QuoteId   string `json:"quoteId,omitempty"`
	ClQuoteId string `json:"clQuoteId,omitempty"`
	RfqId     string `json:"rfqId,omitempty"`
}

// RFQCancelQuoteAck 表示取消报价单返回项。
type RFQCancelQuoteAck struct {
	QuoteId   string `json:"quoteId"`
	ClQuoteId string `json:"clQuoteId"`
	SCode     string `json:"sCode"`
	SMsg      string `json:"sMsg"`
}

// RFQCancelQuoteService 取消报价单。
type RFQCancelQuoteService struct {
	c *Client

	quoteId   string
	clQuoteId string
	rfqId     string
}

// NewRFQCancelQuoteService 创建 RFQCancelQuoteService。
func (c *Client) NewRFQCancelQuoteService() *RFQCancelQuoteService {
	return &RFQCancelQuoteService{c: c}
}

// QuoteId 设置报价单 ID（可选）。
func (s *RFQCancelQuoteService) QuoteId(quoteId string) *RFQCancelQuoteService {
	s.quoteId = quoteId
	return s
}

// ClQuoteId 设置报价单自定义 ID（可选）。
func (s *RFQCancelQuoteService) ClQuoteId(clQuoteId string) *RFQCancelQuoteService {
	s.clQuoteId = clQuoteId
	return s
}

// RfqId 设置询价单 ID（可选）。
func (s *RFQCancelQuoteService) RfqId(rfqId string) *RFQCancelQuoteService {
	s.rfqId = rfqId
	return s
}

var (
	errRFQCancelQuoteMissingId = errors.New("okx: rfq cancel quote requires quoteId or clQuoteId")
	errEmptyRFQCancelQuote     = errors.New("okx: empty rfq cancel quote response")
)

// Do 取消报价单（POST /api/v5/rfq/cancel-quote）。
func (s *RFQCancelQuoteService) Do(ctx context.Context) (*RFQCancelQuoteAck, error) {
	if s.quoteId == "" && s.clQuoteId == "" {
		return nil, errRFQCancelQuoteMissingId
	}

	req := rfqCancelQuoteRequest{
		QuoteId:   s.quoteId,
		ClQuoteId: s.clQuoteId,
		RfqId:     s.rfqId,
	}

	var data []RFQCancelQuoteAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/rfq/cancel-quote", nil, req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/rfq/cancel-quote", requestID, errEmptyRFQCancelQuote)
	}
	if data[0].SCode != "0" {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodPost,
			RequestPath: "/api/v5/rfq/cancel-quote",
			RequestID:   requestID,
			Code:        data[0].SCode,
			Message:     data[0].SMsg,
		}
	}
	return &data[0], nil
}
