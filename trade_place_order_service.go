package okx

import (
	"context"
	"errors"
	"net/http"
)

// PlaceOrderService 下单。
type PlaceOrderService struct {
	c *Client

	instId  string
	tdMode  string
	clOrdId string
	tag     string
	side    string
	ordType string
	px      string
	pxUsd   string
	pxVol   string
	pxType  string
	sz      string
}

// NewPlaceOrderService 创建 PlaceOrderService。
func (c *Client) NewPlaceOrderService() *PlaceOrderService {
	return &PlaceOrderService{c: c}
}

func (s *PlaceOrderService) InstId(instId string) *PlaceOrderService {
	s.instId = instId
	return s
}

func (s *PlaceOrderService) TdMode(tdMode string) *PlaceOrderService {
	s.tdMode = tdMode
	return s
}

func (s *PlaceOrderService) ClOrdId(clOrdId string) *PlaceOrderService {
	s.clOrdId = clOrdId
	return s
}

func (s *PlaceOrderService) Tag(tag string) *PlaceOrderService {
	s.tag = tag
	return s
}

func (s *PlaceOrderService) Side(side string) *PlaceOrderService {
	s.side = side
	return s
}

func (s *PlaceOrderService) OrdType(ordType string) *PlaceOrderService {
	s.ordType = ordType
	return s
}

func (s *PlaceOrderService) Px(px string) *PlaceOrderService {
	s.px = px
	return s
}

func (s *PlaceOrderService) PxUsd(pxUsd string) *PlaceOrderService {
	s.pxUsd = pxUsd
	return s
}

func (s *PlaceOrderService) PxVol(pxVol string) *PlaceOrderService {
	s.pxVol = pxVol
	return s
}

func (s *PlaceOrderService) PxType(pxType string) *PlaceOrderService {
	s.pxType = pxType
	return s
}

func (s *PlaceOrderService) Sz(sz string) *PlaceOrderService {
	s.sz = sz
	return s
}

var (
	errPlaceOrderMissingInstId  = errors.New("okx: place order requires instId")
	errPlaceOrderMissingTdMode  = errors.New("okx: place order requires tdMode")
	errPlaceOrderMissingSide    = errors.New("okx: place order requires side")
	errPlaceOrderMissingOrdType = errors.New("okx: place order requires ordType")
	errPlaceOrderMissingSz      = errors.New("okx: place order requires sz")
	errPlaceOrderMissingPx      = errors.New("okx: place order requires px/pxUsd/pxVol for limit order")
	errEmptyPlaceOrderResponse  = errors.New("okx: empty place order response")
)

type placeOrderRequest struct {
	InstId  string `json:"instId"`
	TdMode  string `json:"tdMode"`
	ClOrdId string `json:"clOrdId,omitempty"`
	Tag     string `json:"tag,omitempty"`
	Side    string `json:"side"`
	OrdType string `json:"ordType"`
	Px      string `json:"px,omitempty"`
	PxUsd   string `json:"pxUsd,omitempty"`
	PxVol   string `json:"pxVol,omitempty"`
	PxType  string `json:"pxType,omitempty"`
	Sz      string `json:"sz"`
}

// Do 下单（POST /api/v5/trade/order）。
func (s *PlaceOrderService) Do(ctx context.Context) (*TradeOrderAck, error) {
	if s.instId == "" {
		return nil, errPlaceOrderMissingInstId
	}
	if s.tdMode == "" {
		return nil, errPlaceOrderMissingTdMode
	}
	if s.side == "" {
		return nil, errPlaceOrderMissingSide
	}
	if s.ordType == "" {
		return nil, errPlaceOrderMissingOrdType
	}
	if s.sz == "" {
		return nil, errPlaceOrderMissingSz
	}
	if s.ordType == "limit" && s.px == "" && s.pxUsd == "" && s.pxVol == "" {
		return nil, errPlaceOrderMissingPx
	}

	req := placeOrderRequest{
		InstId:  s.instId,
		TdMode:  s.tdMode,
		ClOrdId: s.clOrdId,
		Tag:     s.tag,
		Side:    s.side,
		OrdType: s.ordType,
		Px:      s.px,
		PxUsd:   s.pxUsd,
		PxVol:   s.pxVol,
		PxType:  s.pxType,
		Sz:      s.sz,
	}

	var data []TradeOrderAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/trade/order", nil, req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyPlaceOrderResponse
	}
	if data[0].SCode != "" && data[0].SCode != "0" {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodPost,
			RequestPath: "/api/v5/trade/order",
			Code:        data[0].SCode,
			Message:     data[0].SMsg,
		}
	}
	return &data[0], nil
}
