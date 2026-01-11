package okx

import (
	"context"
	"errors"
	"net/http"
)

// SprdAmendOrderService 价差交易改单。
type SprdAmendOrderService struct {
	c *Client

	ordId   string
	clOrdId string
	reqId   string
	newSz   string
	newPx   string
}

// NewSprdAmendOrderService 创建 SprdAmendOrderService。
func (c *Client) NewSprdAmendOrderService() *SprdAmendOrderService {
	return &SprdAmendOrderService{c: c}
}

func (s *SprdAmendOrderService) OrdId(ordId string) *SprdAmendOrderService {
	s.ordId = ordId
	return s
}

func (s *SprdAmendOrderService) ClOrdId(clOrdId string) *SprdAmendOrderService {
	s.clOrdId = clOrdId
	return s
}

// ReqId 设置用户自定义修改事件 ID（可选，1-32 位）。
func (s *SprdAmendOrderService) ReqId(reqId string) *SprdAmendOrderService {
	s.reqId = reqId
	return s
}

func (s *SprdAmendOrderService) NewSz(newSz string) *SprdAmendOrderService {
	s.newSz = newSz
	return s
}

func (s *SprdAmendOrderService) NewPx(newPx string) *SprdAmendOrderService {
	s.newPx = newPx
	return s
}

var (
	errSprdAmendOrderMissingId     = errors.New("okx: sprd amend order requires ordId or clOrdId")
	errSprdAmendOrderMissingChange = errors.New("okx: sprd amend order requires newSz or newPx")
	errEmptySprdAmendOrderResponse = errors.New("okx: empty sprd amend order response")
)

type sprdAmendOrderRequest struct {
	OrdId   string `json:"ordId,omitempty"`
	ClOrdId string `json:"clOrdId,omitempty"`
	ReqId   string `json:"reqId,omitempty"`
	NewSz   string `json:"newSz,omitempty"`
	NewPx   string `json:"newPx,omitempty"`
}

// Do 价差交易改单（POST /api/v5/sprd/amend-order）。
func (s *SprdAmendOrderService) Do(ctx context.Context) (*TradeOrderAck, error) {
	if s.ordId == "" && s.clOrdId == "" {
		return nil, errSprdAmendOrderMissingId
	}
	if s.newSz == "" && s.newPx == "" {
		return nil, errSprdAmendOrderMissingChange
	}

	req := sprdAmendOrderRequest{
		OrdId:   s.ordId,
		ClOrdId: s.clOrdId,
		ReqId:   s.reqId,
		NewSz:   s.newSz,
		NewPx:   s.newPx,
	}

	var data []TradeOrderAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/sprd/amend-order", nil, req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptySprdAmendOrderResponse
	}
	if data[0].SCode != "" && data[0].SCode != "0" {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodPost,
			RequestPath: "/api/v5/sprd/amend-order",
			Code:        data[0].SCode,
			Message:     data[0].SMsg,
		}
	}
	return &data[0], nil
}
