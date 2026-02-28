package okx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// GetOrderService 获取订单信息（单笔）。
type GetOrderService struct {
	c *Client

	instId  string
	ordId   string
	clOrdId string
}

// NewGetOrderService 创建 GetOrderService。
func (c *Client) NewGetOrderService() *GetOrderService {
	return &GetOrderService{c: c}
}

func (s *GetOrderService) InstId(instId string) *GetOrderService {
	s.instId = instId
	return s
}

func (s *GetOrderService) OrdId(ordId string) *GetOrderService {
	s.ordId = ordId
	return s
}

func (s *GetOrderService) ClOrdId(clOrdId string) *GetOrderService {
	s.clOrdId = clOrdId
	return s
}

var (
	errGetOrderMissingInstId   = errors.New("okx: get order requires instId")
	errGetOrderMissingId       = errors.New("okx: get order requires ordId or clOrdId")
	errEmptyGetOrderResponse   = errors.New("okx: empty get order response")
	errInvalidGetOrderResponse = errors.New("okx: invalid get order response")
)

// Do 获取订单信息（GET /api/v5/trade/order）。
func (s *GetOrderService) Do(ctx context.Context) (*TradeOrder, error) {
	if s.instId == "" {
		return nil, errGetOrderMissingInstId
	}
	if s.ordId == "" && s.clOrdId == "" {
		return nil, errGetOrderMissingId
	}

	q := url.Values{}
	q.Set("instId", s.instId)
	if s.ordId != "" {
		q.Set("ordId", s.ordId)
	}
	if s.clOrdId != "" && s.ordId == "" {
		q.Set("clOrdId", s.clOrdId)
	}

	var data []TradeOrder
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodGet, "/api/v5/trade/order", q, nil, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodGet, "/api/v5/trade/order", requestID, errEmptyGetOrderResponse)
	}
	if len(data) != 1 {
		return nil, newInvalidDataAPIError(
			http.MethodGet,
			"/api/v5/trade/order",
			requestID,
			fmt.Errorf("%w: expected 1 item, got %d", errInvalidGetOrderResponse, len(data)),
		)
	}
	return &data[0], nil
}
