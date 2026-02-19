package okx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// SprdCancelOrderService 价差交易撤单。
type SprdCancelOrderService struct {
	c *Client

	ordId   string
	clOrdId string
}

// NewSprdCancelOrderService 创建 SprdCancelOrderService。
func (c *Client) NewSprdCancelOrderService() *SprdCancelOrderService {
	return &SprdCancelOrderService{c: c}
}

// OrdId 设置订单 ID（ordId/clOrdId 必须且只能传一个）。
func (s *SprdCancelOrderService) OrdId(ordId string) *SprdCancelOrderService {
	s.ordId = ordId
	return s
}

// ClOrdId 设置客户自定义订单 ID（ordId/clOrdId 必须且只能传一个）。
func (s *SprdCancelOrderService) ClOrdId(clOrdId string) *SprdCancelOrderService {
	s.clOrdId = clOrdId
	return s
}

var (
	errSprdCancelOrderMissingId     = errors.New("okx: sprd cancel order requires ordId or clOrdId")
	errSprdCancelOrderTooManyId     = errors.New("okx: sprd cancel order requires exactly one of ordId or clOrdId")
	errEmptySprdCancelOrderResponse = errors.New("okx: empty sprd cancel order response")
	errInvalidSprdCancelOrderResp   = errors.New("okx: invalid sprd cancel order response")
)

type sprdCancelOrderRequest struct {
	OrdId   string `json:"ordId,omitempty"`
	ClOrdId string `json:"clOrdId,omitempty"`
}

// Do 价差交易撤单（POST /api/v5/sprd/cancel-order）。
func (s *SprdCancelOrderService) Do(ctx context.Context) (*TradeOrderAck, error) {
	switch countNonEmptyStrings(s.ordId, s.clOrdId) {
	case 0:
		return nil, errSprdCancelOrderMissingId
	case 2:
		return nil, errSprdCancelOrderTooManyId
	}

	req := sprdCancelOrderRequest{
		OrdId:   s.ordId,
		ClOrdId: s.clOrdId,
	}

	var data []TradeOrderAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/sprd/cancel-order", nil, req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/sprd/cancel-order", requestID, errEmptySprdCancelOrderResponse)
	}
	if len(data) != 1 {
		return nil, newInvalidDataAPIError(
			http.MethodPost,
			"/api/v5/sprd/cancel-order",
			requestID,
			fmt.Errorf("%w: expected 1 ack, got %d", errInvalidSprdCancelOrderResp, len(data)),
		)
	}
	if data[0].SCode != "0" {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodPost,
			RequestPath: "/api/v5/sprd/cancel-order",
			RequestID:   requestID,
			Code:        data[0].SCode,
			Message:     data[0].SMsg,
		}
	}
	return &data[0], nil
}
