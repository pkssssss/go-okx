package okx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const tradeBatchMaxOrders = 20

// BatchPlaceOrder 表示批量下单的单笔请求（精简版）。
type BatchPlaceOrder struct {
	InstId  string `json:"instId"`
	TdMode  string `json:"tdMode"`
	Ccy     string `json:"ccy,omitempty"`
	ClOrdId string `json:"clOrdId,omitempty"`
	Tag     string `json:"tag,omitempty"`
	Side    string `json:"side"`
	PosSide string `json:"posSide,omitempty"`
	OrdType string `json:"ordType"`

	Px    string `json:"px,omitempty"`
	PxUsd string `json:"pxUsd,omitempty"`
	PxVol string `json:"pxVol,omitempty"`

	Sz string `json:"sz"`

	ReduceOnly    *bool  `json:"reduceOnly,omitempty"`
	TgtCcy        string `json:"tgtCcy,omitempty"`
	BanAmend      *bool  `json:"banAmend,omitempty"`
	PxAmendType   string `json:"pxAmendType,omitempty"`
	StpMode       string `json:"stpMode,omitempty"`
	TradeQuoteCcy string `json:"tradeQuoteCcy,omitempty"`
}

// BatchPlaceOrdersService 批量下单。
type BatchPlaceOrdersService struct {
	c             *Client
	orders        []BatchPlaceOrder
	expTimeHeader string
}

// NewBatchPlaceOrdersService 创建 BatchPlaceOrdersService。
func (c *Client) NewBatchPlaceOrdersService() *BatchPlaceOrdersService {
	return &BatchPlaceOrdersService{c: c}
}

// Orders 设置批量下单列表（最多 20 个）。
func (s *BatchPlaceOrdersService) Orders(orders []BatchPlaceOrder) *BatchPlaceOrdersService {
	s.orders = orders
	return s
}

// ExpTime 设置 REST 请求头 expTime（请求有效截止时间，Unix 毫秒时间戳字符串）。
func (s *BatchPlaceOrdersService) ExpTime(expTimeMillis string) *BatchPlaceOrdersService {
	s.expTimeHeader = expTimeMillis
	return s
}

var (
	errBatchPlaceOrdersMissingOrders = errors.New("okx: batch place orders requires at least one order")
	errBatchPlaceOrdersTooManyOrders = errors.New("okx: batch place orders max 20 orders")
)

// Do 批量下单（POST /api/v5/trade/batch-orders）。
func (s *BatchPlaceOrdersService) Do(ctx context.Context) ([]TradeOrderAck, error) {
	if len(s.orders) == 0 {
		return nil, errBatchPlaceOrdersMissingOrders
	}
	if len(s.orders) > tradeBatchMaxOrders {
		return nil, errBatchPlaceOrdersTooManyOrders
	}

	req := make([]BatchPlaceOrder, 0, len(s.orders))
	for i, o := range s.orders {
		if o.InstId == "" {
			return nil, fmt.Errorf("okx: batch place orders[%d] missing instId", i)
		}
		if o.TdMode == "" {
			return nil, fmt.Errorf("okx: batch place orders[%d] missing tdMode", i)
		}
		if o.Side == "" {
			return nil, fmt.Errorf("okx: batch place orders[%d] missing side", i)
		}
		if o.OrdType == "" {
			return nil, fmt.Errorf("okx: batch place orders[%d] missing ordType", i)
		}
		if o.Sz == "" {
			return nil, fmt.Errorf("okx: batch place orders[%d] missing sz", i)
		}

		if countNonEmptyStrings(o.Px, o.PxUsd, o.PxVol) > 1 {
			return nil, fmt.Errorf("okx: batch place orders[%d] requires at most one of px/pxUsd/pxVol", i)
		}
		if requiresPriceForOrderType(o.OrdType) && o.Px == "" && o.PxUsd == "" && o.PxVol == "" {
			return nil, fmt.Errorf("okx: batch place orders[%d] missing px/pxUsd/pxVol for this ordType", i)
		}
		req = append(req, o)
	}

	var data []TradeOrderAck
	var header http.Header
	if s.expTimeHeader != "" {
		header = make(http.Header)
		header.Set("expTime", s.expTimeHeader)
	}
	if err := s.c.doWithHeaders(ctx, http.MethodPost, "/api/v5/trade/batch-orders", nil, req, true, header, &data); err != nil {
		return nil, err
	}
	if err := tradeCheckBatchAcks(http.MethodPost, "/api/v5/trade/batch-orders", data); err != nil {
		return data, err
	}
	return data, nil
}
