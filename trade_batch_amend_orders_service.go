package okx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// BatchAmendOrder 表示批量改单的单笔请求（精简版）。
type BatchAmendOrder struct {
	InstId    string `json:"instId"`
	CxlOnFail *bool  `json:"cxlOnFail,omitempty"`
	OrdId     string `json:"ordId,omitempty"`
	ClOrdId   string `json:"clOrdId,omitempty"`
	ReqId     string `json:"reqId,omitempty"`
	NewSz     string `json:"newSz,omitempty"`
	NewPx     string `json:"newPx,omitempty"`
}

// BatchAmendOrdersService 批量改单。
type BatchAmendOrdersService struct {
	c      *Client
	orders []BatchAmendOrder
}

// NewBatchAmendOrdersService 创建 BatchAmendOrdersService。
func (c *Client) NewBatchAmendOrdersService() *BatchAmendOrdersService {
	return &BatchAmendOrdersService{c: c}
}

// Orders 设置批量改单列表（最多 20 个）。
func (s *BatchAmendOrdersService) Orders(orders []BatchAmendOrder) *BatchAmendOrdersService {
	s.orders = orders
	return s
}

var (
	errBatchAmendOrdersMissingOrders = errors.New("okx: batch amend orders requires at least one order")
	errBatchAmendOrdersTooManyOrders = errors.New("okx: batch amend orders max 20 orders")
)

// Do 批量改单（POST /api/v5/trade/amend-batch-orders）。
func (s *BatchAmendOrdersService) Do(ctx context.Context) ([]TradeOrderAck, error) {
	if len(s.orders) == 0 {
		return nil, errBatchAmendOrdersMissingOrders
	}
	if len(s.orders) > tradeBatchMaxOrders {
		return nil, errBatchAmendOrdersTooManyOrders
	}

	req := make([]BatchAmendOrder, 0, len(s.orders))
	for i, o := range s.orders {
		if o.InstId == "" {
			return nil, fmt.Errorf("okx: batch amend orders[%d] missing instId", i)
		}
		if o.OrdId == "" && o.ClOrdId == "" {
			return nil, fmt.Errorf("okx: batch amend orders[%d] missing ordId or clOrdId", i)
		}
		if o.NewSz == "" && o.NewPx == "" {
			return nil, fmt.Errorf("okx: batch amend orders[%d] missing newSz or newPx", i)
		}
		if o.OrdId != "" {
			o.ClOrdId = ""
		}
		req = append(req, o)
	}

	var data []TradeOrderAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/trade/amend-batch-orders", nil, req, true, &data); err != nil {
		return nil, err
	}
	if err := tradeCheckBatchAcks(http.MethodPost, "/api/v5/trade/amend-batch-orders", data); err != nil {
		return data, err
	}
	return data, nil
}
