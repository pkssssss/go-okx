package okx

import (
	"context"
	"errors"
	"net/http"
)

// SprdPlaceOrderService 价差交易下单。
type SprdPlaceOrderService struct {
	c *Client

	sprdId  string
	clOrdId string
	tag     string
	side    string
	ordType string
	px      string
	sz      string
}

// NewSprdPlaceOrderService 创建 SprdPlaceOrderService。
func (c *Client) NewSprdPlaceOrderService() *SprdPlaceOrderService {
	return &SprdPlaceOrderService{c: c}
}

// SprdId 设置 Spread ID（必填）。
func (s *SprdPlaceOrderService) SprdId(sprdId string) *SprdPlaceOrderService {
	s.sprdId = sprdId
	return s
}

// ClOrdId 设置客户自定义订单 ID（可选，1-32 位）。
func (s *SprdPlaceOrderService) ClOrdId(clOrdId string) *SprdPlaceOrderService {
	s.clOrdId = clOrdId
	return s
}

// Tag 设置订单标签（可选，1-16 位）。
func (s *SprdPlaceOrderService) Tag(tag string) *SprdPlaceOrderService {
	s.tag = tag
	return s
}

// Side 设置订单方向（必填）：buy/sell。
func (s *SprdPlaceOrderService) Side(side string) *SprdPlaceOrderService {
	s.side = side
	return s
}

// OrdType 设置订单类型（必填）：market/limit/post_only/ioc。
func (s *SprdPlaceOrderService) OrdType(ordType string) *SprdPlaceOrderService {
	s.ordType = ordType
	return s
}

// Px 设置委托价格（是否必填以 OKX 服务端规则为准）。
func (s *SprdPlaceOrderService) Px(px string) *SprdPlaceOrderService {
	s.px = px
	return s
}

// Sz 设置委托数量（必填）。
func (s *SprdPlaceOrderService) Sz(sz string) *SprdPlaceOrderService {
	s.sz = sz
	return s
}

var (
	errSprdPlaceOrderMissingSprdId  = errors.New("okx: sprd place order requires sprdId")
	errSprdPlaceOrderMissingSide    = errors.New("okx: sprd place order requires side")
	errSprdPlaceOrderMissingOrdType = errors.New("okx: sprd place order requires ordType")
	errSprdPlaceOrderMissingSz      = errors.New("okx: sprd place order requires sz")
	errEmptySprdPlaceOrderResponse  = errors.New("okx: empty sprd place order response")
)

type sprdPlaceOrderRequest struct {
	SprdId  string `json:"sprdId"`
	ClOrdId string `json:"clOrdId,omitempty"`
	Tag     string `json:"tag,omitempty"`
	Side    string `json:"side"`
	OrdType string `json:"ordType"`
	Px      string `json:"px,omitempty"`
	Sz      string `json:"sz"`
}

// Do 价差交易下单（POST /api/v5/sprd/order）。
func (s *SprdPlaceOrderService) Do(ctx context.Context) (*TradeOrderAck, error) {
	if s.sprdId == "" {
		return nil, errSprdPlaceOrderMissingSprdId
	}
	if s.side == "" {
		return nil, errSprdPlaceOrderMissingSide
	}
	if s.ordType == "" {
		return nil, errSprdPlaceOrderMissingOrdType
	}
	if s.sz == "" {
		return nil, errSprdPlaceOrderMissingSz
	}

	req := sprdPlaceOrderRequest{
		SprdId:  s.sprdId,
		ClOrdId: s.clOrdId,
		Tag:     s.tag,
		Side:    s.side,
		OrdType: s.ordType,
		Px:      s.px,
		Sz:      s.sz,
	}

	var data []TradeOrderAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/sprd/order", nil, req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptySprdPlaceOrderResponse
	}
	if data[0].SCode != "" && data[0].SCode != "0" {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodPost,
			RequestPath: "/api/v5/sprd/order",
			Code:        data[0].SCode,
			Message:     data[0].SMsg,
		}
	}
	return &data[0], nil
}
