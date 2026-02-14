package okx

import (
	"context"
	"errors"
	"net/http"
)

// CancelOrderService 撤单。
type CancelOrderService struct {
	c *Client

	instId  string
	ordId   string
	clOrdId string
}

// NewCancelOrderService 创建 CancelOrderService。
func (c *Client) NewCancelOrderService() *CancelOrderService {
	return &CancelOrderService{c: c}
}

func (s *CancelOrderService) InstId(instId string) *CancelOrderService {
	s.instId = instId
	return s
}

func (s *CancelOrderService) OrdId(ordId string) *CancelOrderService {
	s.ordId = ordId
	return s
}

func (s *CancelOrderService) ClOrdId(clOrdId string) *CancelOrderService {
	s.clOrdId = clOrdId
	return s
}

var (
	errCancelOrderMissingInstId = errors.New("okx: cancel order requires instId")
	errCancelOrderMissingId     = errors.New("okx: cancel order requires ordId or clOrdId")
	errCancelOrderTooManyId     = errors.New("okx: cancel order requires exactly one of ordId or clOrdId")
	errEmptyCancelOrderResponse = errors.New("okx: empty cancel order response")
)

type cancelOrderRequest struct {
	InstId  string `json:"instId"`
	OrdId   string `json:"ordId,omitempty"`
	ClOrdId string `json:"clOrdId,omitempty"`
}

// Do 撤单（POST /api/v5/trade/cancel-order）。
func (s *CancelOrderService) Do(ctx context.Context) (*TradeOrderAck, error) {
	if s.instId == "" {
		return nil, errCancelOrderMissingInstId
	}
	switch countNonEmptyStrings(s.ordId, s.clOrdId) {
	case 0:
		return nil, errCancelOrderMissingId
	case 2:
		return nil, errCancelOrderTooManyId
	}

	req := cancelOrderRequest{
		InstId:  s.instId,
		OrdId:   s.ordId,
		ClOrdId: s.clOrdId,
	}

	var data []TradeOrderAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/trade/cancel-order", nil, req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodPost,
			RequestPath: "/api/v5/trade/cancel-order",
			RequestID:   requestID,
			Code:        "0",
			Message:     errEmptyCancelOrderResponse.Error(),
		}
	}
	if data[0].SCode != "0" {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodPost,
			RequestPath: "/api/v5/trade/cancel-order",
			RequestID:   requestID,
			Code:        data[0].SCode,
			Message:     data[0].SMsg,
		}
	}
	return &data[0], nil
}
