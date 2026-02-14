package okx

import (
	"context"
	"errors"
	"net/http"
)

// RFQMMPResetService 重设 MMP 状态为无效。
type RFQMMPResetService struct {
	c *Client
}

// NewRFQMMPResetService 创建 RFQMMPResetService。
func (c *Client) NewRFQMMPResetService() *RFQMMPResetService {
	return &RFQMMPResetService{c: c}
}

var errEmptyRFQMMPResetResponse = errors.New("okx: empty rfq mmp reset response")

// Do 重设 MMP 状态（POST /api/v5/rfq/mmp-reset）。
func (s *RFQMMPResetService) Do(ctx context.Context) (*RFQTsAck, error) {
	var data []RFQTsAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/rfq/mmp-reset", nil, nil, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/rfq/mmp-reset", requestID, errEmptyRFQMMPResetResponse)
	}
	return &data[0], nil
}
