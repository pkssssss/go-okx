package okx

import (
	"context"
	"errors"
	"net/http"
)

const rfqMaxCancelBatch = 100

type rfqCancelBatchRFQsRequest struct {
	RfqIds   []string `json:"rfqIds,omitempty"`
	ClRfqIds []string `json:"clRfqIds,omitempty"`
}

// RFQCancelBatchRFQsService 批量取消询价单。
type RFQCancelBatchRFQsService struct {
	c *Client

	rfqIds   []string
	clRfqIds []string
}

// NewRFQCancelBatchRFQsService 创建 RFQCancelBatchRFQsService。
func (c *Client) NewRFQCancelBatchRFQsService() *RFQCancelBatchRFQsService {
	return &RFQCancelBatchRFQsService{c: c}
}

// RfqIds 设置询价单 ID 列表（可选，最多 100 个）。
func (s *RFQCancelBatchRFQsService) RfqIds(rfqIds []string) *RFQCancelBatchRFQsService {
	s.rfqIds = rfqIds
	return s
}

// ClRfqIds 设置询价单自定义 ID 列表（可选，最多 100 个）。
func (s *RFQCancelBatchRFQsService) ClRfqIds(clRfqIds []string) *RFQCancelBatchRFQsService {
	s.clRfqIds = clRfqIds
	return s
}

var (
	errRFQCancelBatchRFQsMissingIds = errors.New("okx: rfq cancel batch rfqs requires rfqIds or clRfqIds")
	errRFQCancelBatchRFQsTooManyIds = errors.New("okx: rfq cancel batch rfqs max 100 ids")
)

// Do 批量取消询价单（POST /api/v5/rfq/cancel-batch-rfqs）。
func (s *RFQCancelBatchRFQsService) Do(ctx context.Context) ([]RFQCancelAck, error) {
	if len(s.rfqIds) == 0 && len(s.clRfqIds) == 0 {
		return nil, errRFQCancelBatchRFQsMissingIds
	}
	if len(s.rfqIds) > rfqMaxCancelBatch || len(s.clRfqIds) > rfqMaxCancelBatch {
		return nil, errRFQCancelBatchRFQsTooManyIds
	}

	req := rfqCancelBatchRFQsRequest{
		RfqIds:   s.rfqIds,
		ClRfqIds: s.clRfqIds,
	}

	var data []RFQCancelAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/rfq/cancel-batch-rfqs", nil, req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if err := rfqCheckCancelBatchRFQs(http.MethodPost, "/api/v5/rfq/cancel-batch-rfqs", requestID, data); err != nil {
		return data, err
	}
	return data, nil
}
