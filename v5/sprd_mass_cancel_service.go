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
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/sprd/mass-cancel", nil, req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptySprdMassCancelResponse
	}
	return &data[0], nil
}
