package okx

import (
	"context"
	"errors"
	"net/http"
)

// SprdMassCancelAck 表示全部撤单返回项。
type SprdMassCancelAck struct {
	Result bool `json:"result"`
}

// SprdMassCancelService 全部撤单（价差交易）。
type SprdMassCancelService struct {
	c *Client

	sprdId string
}

// NewSprdMassCancelService 创建 SprdMassCancelService。
func (c *Client) NewSprdMassCancelService() *SprdMassCancelService {
	return &SprdMassCancelService{c: c}
}

// SprdId 设置 spread ID（可选）。
func (s *SprdMassCancelService) SprdId(sprdId string) *SprdMassCancelService {
	s.sprdId = sprdId
	return s
}

var errEmptySprdMassCancelResponse = errors.New("okx: empty sprd mass cancel response")

type sprdMassCancelRequest struct {
	SprdId string `json:"sprdId,omitempty"`
}

// Do 全部撤单（POST /api/v5/sprd/mass-cancel）。
func (s *SprdMassCancelService) Do(ctx context.Context) (*SprdMassCancelAck, error) {
	req := sprdMassCancelRequest{SprdId: s.sprdId}

	var data []SprdMassCancelAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/sprd/mass-cancel", nil, req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptySprdMassCancelResponse
	}
	if !data[0].Result {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodPost,
			RequestPath: "/api/v5/sprd/mass-cancel",
			RequestID:   requestID,
			Code:        "0",
			Message:     "sprd mass cancel result is false",
		}
	}
	return &data[0], nil
}
