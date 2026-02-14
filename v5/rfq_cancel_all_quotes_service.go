package okx

import (
	"context"
	"errors"
	"net/http"
)

// RFQCancelAllQuotesService 取消所有报价单。
type RFQCancelAllQuotesService struct {
	c *Client
}

// NewRFQCancelAllQuotesService 创建 RFQCancelAllQuotesService。
func (c *Client) NewRFQCancelAllQuotesService() *RFQCancelAllQuotesService {
	return &RFQCancelAllQuotesService{c: c}
}

var errEmptyRFQCancelAllQuotesResponse = errors.New("okx: empty rfq cancel all quotes response")

// Do 取消所有报价单（POST /api/v5/rfq/cancel-all-quotes）。
func (s *RFQCancelAllQuotesService) Do(ctx context.Context) (*RFQTsAck, error) {
	var data []RFQTsAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/rfq/cancel-all-quotes", nil, nil, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/rfq/cancel-all-quotes", requestID, errEmptyRFQCancelAllQuotesResponse)
	}
	return &data[0], nil
}
