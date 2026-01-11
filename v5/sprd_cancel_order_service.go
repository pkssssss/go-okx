package okx

import (
	"context"
	"errors"
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

// OrdId 设置订单 ID（ordId/clOrdId 必须传一个；若都传，以 ordId 为主）。
func (s *SprdCancelOrderService) OrdId(ordId string) *SprdCancelOrderService {
	s.ordId = ordId
	return s
}

// ClOrdId 设置客户自定义订单 ID（ordId/clOrdId 必须传一个）。
func (s *SprdCancelOrderService) ClOrdId(clOrdId string) *SprdCancelOrderService {
	s.clOrdId = clOrdId
	return s
}

var (
	errSprdCancelOrderMissingId     = errors.New("okx: sprd cancel order requires ordId or clOrdId")
	errEmptySprdCancelOrderResponse = errors.New("okx: empty sprd cancel order response")
)

type sprdCancelOrderRequest struct {
	OrdId   string `json:"ordId,omitempty"`
	ClOrdId string `json:"clOrdId,omitempty"`
}

// Do 价差交易撤单（POST /api/v5/sprd/cancel-order）。
func (s *SprdCancelOrderService) Do(ctx context.Context) (*TradeOrderAck, error) {
	if s.ordId == "" && s.clOrdId == "" {
		return nil, errSprdCancelOrderMissingId
	}

	req := sprdCancelOrderRequest{
		OrdId:   s.ordId,
		ClOrdId: s.clOrdId,
	}

	var data []TradeOrderAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/sprd/cancel-order", nil, req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptySprdCancelOrderResponse
	}
	if data[0].SCode != "" && data[0].SCode != "0" {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodPost,
			RequestPath: "/api/v5/sprd/cancel-order",
			Code:        data[0].SCode,
			Message:     data[0].SMsg,
		}
	}
	return &data[0], nil
}
