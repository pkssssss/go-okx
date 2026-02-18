package okx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const tradeAlgoMaxCancelOrders = 10

// CancelAlgoOrder 表示撤销策略委托单的单笔请求。
type CancelAlgoOrder struct {
	InstId      string `json:"instId"`
	AlgoId      string `json:"algoId,omitempty"`
	AlgoClOrdId string `json:"algoClOrdId,omitempty"`
}

// CancelAlgoOrdersService 撤销策略委托订单（最多 10 个）。
type CancelAlgoOrdersService struct {
	c      *Client
	orders []CancelAlgoOrder
}

// NewCancelAlgoOrdersService 创建 CancelAlgoOrdersService。
func (c *Client) NewCancelAlgoOrdersService() *CancelAlgoOrdersService {
	return &CancelAlgoOrdersService{c: c}
}

// Orders 设置撤销策略委托单列表（最多 10 个）。
func (s *CancelAlgoOrdersService) Orders(orders []CancelAlgoOrder) *CancelAlgoOrdersService {
	s.orders = orders
	return s
}

var (
	errCancelAlgoOrdersMissingOrders = errors.New("okx: cancel algos requires at least one order")
	errCancelAlgoOrdersTooManyOrders = errors.New("okx: cancel algos max 10 orders")
)

// Do 撤销策略委托订单（POST /api/v5/trade/cancel-algos）。
func (s *CancelAlgoOrdersService) Do(ctx context.Context) ([]TradeAlgoOrderAck, error) {
	if len(s.orders) == 0 {
		return nil, errCancelAlgoOrdersMissingOrders
	}
	if len(s.orders) > tradeAlgoMaxCancelOrders {
		return nil, errCancelAlgoOrdersTooManyOrders
	}

	req := make([]CancelAlgoOrder, 0, len(s.orders))
	for i, o := range s.orders {
		if o.InstId == "" {
			return nil, fmt.Errorf("okx: cancel algos[%d] missing instId", i)
		}
		if o.AlgoId == "" && o.AlgoClOrdId == "" {
			return nil, fmt.Errorf("okx: cancel algos[%d] requires algoId or algoClOrdId", i)
		}

		normalized := o
		if normalized.AlgoId != "" {
			normalized.AlgoClOrdId = ""
		}
		req = append(req, normalized)
	}

	var data []TradeAlgoOrderAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/trade/cancel-algos", nil, req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if err := tradeCheckAlgoAcks(http.MethodPost, "/api/v5/trade/cancel-algos", requestID, len(req), data); err != nil {
		return data, err
	}
	return data, nil
}
