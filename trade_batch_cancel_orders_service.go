package okx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// BatchCancelOrder 表示批量撤单的单笔请求。
type BatchCancelOrder struct {
	InstId  string `json:"instId"`
	OrdId   string `json:"ordId,omitempty"`
	ClOrdId string `json:"clOrdId,omitempty"`
}

// BatchCancelOrdersService 批量撤单。
type BatchCancelOrdersService struct {
	c      *Client
	orders []BatchCancelOrder
}

// NewBatchCancelOrdersService 创建 BatchCancelOrdersService。
func (c *Client) NewBatchCancelOrdersService() *BatchCancelOrdersService {
	return &BatchCancelOrdersService{c: c}
}

// Orders 设置批量撤单列表（最多 20 个）。
func (s *BatchCancelOrdersService) Orders(orders []BatchCancelOrder) *BatchCancelOrdersService {
	s.orders = orders
	return s
}

var (
	errBatchCancelOrdersMissingOrders = errors.New("okx: batch cancel orders requires at least one order")
	errBatchCancelOrdersTooManyOrders = errors.New("okx: batch cancel orders max 20 orders")
)

// Do 批量撤单（POST /api/v5/trade/cancel-batch-orders）。
func (s *BatchCancelOrdersService) Do(ctx context.Context) ([]TradeOrderAck, error) {
	if len(s.orders) == 0 {
		return nil, errBatchCancelOrdersMissingOrders
	}
	if len(s.orders) > tradeBatchMaxOrders {
		return nil, errBatchCancelOrdersTooManyOrders
	}

	req := make([]BatchCancelOrder, 0, len(s.orders))
	for i, o := range s.orders {
		if o.InstId == "" {
			return nil, fmt.Errorf("okx: batch cancel orders[%d] missing instId", i)
		}
		if o.OrdId == "" && o.ClOrdId == "" {
			return nil, fmt.Errorf("okx: batch cancel orders[%d] missing ordId or clOrdId", i)
		}
		if o.OrdId != "" {
			o.ClOrdId = ""
		}
		req = append(req, o)
	}

	var data []TradeOrderAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/trade/cancel-batch-orders", nil, req, true, &data); err != nil {
		return nil, err
	}
	if err := tradeCheckBatchAcks(http.MethodPost, "/api/v5/trade/cancel-batch-orders", data); err != nil {
		return data, err
	}
	return data, nil
}
