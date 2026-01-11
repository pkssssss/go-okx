package okx

import (
	"context"
	"errors"
	"net/http"
)

type rfqCancelRFQRequest struct {
	RfqId   string `json:"rfqId,omitempty"`
	ClRfqId string `json:"clRfqId,omitempty"`
}

// RFQCancelAck 表示取消询价单返回项。
type RFQCancelAck struct {
	RfqId   string `json:"rfqId"`
	ClRfqId string `json:"clRfqId"`
	SCode   string `json:"sCode"`
	SMsg    string `json:"sMsg"`
}

// RFQCancelRFQService 取消询价单。
type RFQCancelRFQService struct {
	c *Client

	rfqId   string
	clRfqId string
}

// NewRFQCancelRFQService 创建 RFQCancelRFQService。
func (c *Client) NewRFQCancelRFQService() *RFQCancelRFQService {
	return &RFQCancelRFQService{c: c}
}

// RfqId 设置询价单 ID（可选）。
func (s *RFQCancelRFQService) RfqId(rfqId string) *RFQCancelRFQService {
	s.rfqId = rfqId
	return s
}

// ClRfqId 设置询价单自定义 ID（可选）。
func (s *RFQCancelRFQService) ClRfqId(clRfqId string) *RFQCancelRFQService {
	s.clRfqId = clRfqId
	return s
}

var (
	errRFQCancelRFQMissingId = errors.New("okx: rfq cancel rfq requires rfqId or clRfqId")
	errEmptyRFQCancelRFQ     = errors.New("okx: empty rfq cancel rfq response")
)

// Do 取消询价单（POST /api/v5/rfq/cancel-rfq）。
func (s *RFQCancelRFQService) Do(ctx context.Context) (*RFQCancelAck, error) {
	if s.rfqId == "" && s.clRfqId == "" {
		return nil, errRFQCancelRFQMissingId
	}

	req := rfqCancelRFQRequest{
		RfqId:   s.rfqId,
		ClRfqId: s.clRfqId,
	}

	var data []RFQCancelAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/rfq/cancel-rfq", nil, req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyRFQCancelRFQ
	}
	if data[0].SCode != "" && data[0].SCode != "0" {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodPost,
			RequestPath: "/api/v5/rfq/cancel-rfq",
			Code:        data[0].SCode,
			Message:     data[0].SMsg,
		}
	}
	return &data[0], nil
}
