package okx

import (
	"context"
	"errors"
	"net/http"
)

// RFQCancelAllRFQsService 取消所有询价单。
type RFQCancelAllRFQsService struct {
	c *Client
}

// NewRFQCancelAllRFQsService 创建 RFQCancelAllRFQsService。
func (c *Client) NewRFQCancelAllRFQsService() *RFQCancelAllRFQsService {
	return &RFQCancelAllRFQsService{c: c}
}

var errEmptyRFQCancelAllRFQsResponse = errors.New("okx: empty rfq cancel all rfqs response")

// Do 取消所有询价单（POST /api/v5/rfq/cancel-all-rfqs）。
func (s *RFQCancelAllRFQsService) Do(ctx context.Context) (*RFQTsAck, error) {
	var data []RFQTsAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/rfq/cancel-all-rfqs", nil, nil, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/rfq/cancel-all-rfqs", requestID, errEmptyRFQCancelAllRFQsResponse)
	}
	return &data[0], nil
}
