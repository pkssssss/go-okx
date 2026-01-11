package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// SprdGetOrderService 获取价差交易订单信息（单笔）。
type SprdGetOrderService struct {
	c *Client

	ordId   string
	clOrdId string
}

// NewSprdGetOrderService 创建 SprdGetOrderService。
func (c *Client) NewSprdGetOrderService() *SprdGetOrderService {
	return &SprdGetOrderService{c: c}
}

func (s *SprdGetOrderService) OrdId(ordId string) *SprdGetOrderService {
	s.ordId = ordId
	return s
}

func (s *SprdGetOrderService) ClOrdId(clOrdId string) *SprdGetOrderService {
	s.clOrdId = clOrdId
	return s
}

var (
	errSprdGetOrderMissingId     = errors.New("okx: sprd get order requires ordId or clOrdId")
	errEmptySprdGetOrderResponse = errors.New("okx: empty sprd get order response")
)

// Do 获取价差交易订单信息（GET /api/v5/sprd/order）。
func (s *SprdGetOrderService) Do(ctx context.Context) (*SprdOrder, error) {
	if s.ordId == "" && s.clOrdId == "" {
		return nil, errSprdGetOrderMissingId
	}

	q := url.Values{}
	if s.ordId != "" {
		q.Set("ordId", s.ordId)
	}
	if s.clOrdId != "" && s.ordId == "" {
		q.Set("clOrdId", s.clOrdId)
	}

	var data []SprdOrder
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/sprd/order", q, nil, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptySprdGetOrderResponse
	}
	return &data[0], nil
}
