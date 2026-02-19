package okx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// PlaceOrderService 下单。
type PlaceOrderService struct {
	c *Client

	instId  string
	tdMode  string
	ccy     string
	clOrdId string
	tag     string
	side    string
	posSide string
	ordType string
	px      string
	pxUsd   string
	pxVol   string
	sz      string

	reduceOnly    *bool
	tgtCcy        string
	banAmend      *bool
	pxAmendType   string
	tradeQuoteCcy string
	stpMode       string
	expTimeHeader string
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

func (s *PlaceOrderService) Ccy(ccy string) *PlaceOrderService {
	s.ccy = ccy
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

func (s *PlaceOrderService) PosSide(posSide string) *PlaceOrderService {
	s.posSide = posSide
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

func (s *PlaceOrderService) Sz(sz string) *PlaceOrderService {
	s.sz = sz
	return s
}

func (s *PlaceOrderService) ReduceOnly(enable bool) *PlaceOrderService {
	s.reduceOnly = &enable
	return s
}

func (s *PlaceOrderService) TgtCcy(tgtCcy string) *PlaceOrderService {
	s.tgtCcy = tgtCcy
	return s
}

func (s *PlaceOrderService) BanAmend(enable bool) *PlaceOrderService {
	s.banAmend = &enable
	return s
}

func (s *PlaceOrderService) PxAmendType(pxAmendType string) *PlaceOrderService {
	s.pxAmendType = pxAmendType
	return s
}

func (s *PlaceOrderService) TradeQuoteCcy(tradeQuoteCcy string) *PlaceOrderService {
	s.tradeQuoteCcy = tradeQuoteCcy
	return s
}

func (s *PlaceOrderService) StpMode(stpMode string) *PlaceOrderService {
	s.stpMode = stpMode
	return s
}

// ExpTime 设置 REST 请求头 expTime（请求有效截止时间，Unix 毫秒时间戳字符串）。
func (s *PlaceOrderService) ExpTime(expTimeMillis string) *PlaceOrderService {
	s.expTimeHeader = expTimeMillis
	return s
}

var (
	errPlaceOrderMissingInstId  = errors.New("okx: place order requires instId")
	errPlaceOrderMissingTdMode  = errors.New("okx: place order requires tdMode")
	errPlaceOrderMissingSide    = errors.New("okx: place order requires side")
	errPlaceOrderMissingOrdType = errors.New("okx: place order requires ordType")
	errPlaceOrderMissingSz      = errors.New("okx: place order requires sz")
	errPlaceOrderTooManyPx      = errors.New("okx: place order requires at most one of px/pxUsd/pxVol")
	errEmptyPlaceOrderResponse  = errors.New("okx: empty place order response")
	errInvalidPlaceOrderResp    = errors.New("okx: invalid place order response")
)

type placeOrderRequest struct {
	InstId  string `json:"instId"`
	TdMode  string `json:"tdMode"`
	Ccy     string `json:"ccy,omitempty"`
	ClOrdId string `json:"clOrdId,omitempty"`
	Tag     string `json:"tag,omitempty"`
	Side    string `json:"side"`
	PosSide string `json:"posSide,omitempty"`
	OrdType string `json:"ordType"`
	Px      string `json:"px,omitempty"`
	PxUsd   string `json:"pxUsd,omitempty"`
	PxVol   string `json:"pxVol,omitempty"`
	Sz      string `json:"sz"`

	ReduceOnly    *bool  `json:"reduceOnly,omitempty"`
	TgtCcy        string `json:"tgtCcy,omitempty"`
	BanAmend      *bool  `json:"banAmend,omitempty"`
	PxAmendType   string `json:"pxAmendType,omitempty"`
	TradeQuoteCcy string `json:"tradeQuoteCcy,omitempty"`
	StpMode       string `json:"stpMode,omitempty"`
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

	if countNonEmptyStrings(s.px, s.pxUsd, s.pxVol) > 1 {
		return nil, errPlaceOrderTooManyPx
	}

	req := placeOrderRequest{
		InstId:  s.instId,
		TdMode:  s.tdMode,
		Ccy:     s.ccy,
		ClOrdId: s.clOrdId,
		Tag:     s.tag,
		Side:    s.side,
		PosSide: s.posSide,
		OrdType: s.ordType,
		Px:      s.px,
		PxUsd:   s.pxUsd,
		PxVol:   s.pxVol,
		Sz:      s.sz,

		ReduceOnly:    s.reduceOnly,
		TgtCcy:        s.tgtCcy,
		BanAmend:      s.banAmend,
		PxAmendType:   s.pxAmendType,
		TradeQuoteCcy: s.tradeQuoteCcy,
		StpMode:       s.stpMode,
	}

	var data []TradeOrderAck
	var header http.Header
	if s.expTimeHeader != "" {
		header = make(http.Header)
		header.Set("expTime", s.expTimeHeader)
	}
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/trade/order", nil, req, true, header, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodPost,
			RequestPath: "/api/v5/trade/order",
			RequestID:   requestID,
			Code:        "0",
			Message:     errEmptyPlaceOrderResponse.Error(),
		}
	}
	if len(data) != 1 {
		return nil, newInvalidDataAPIError(
			http.MethodPost,
			"/api/v5/trade/order",
			requestID,
			fmt.Errorf("%w: expected 1 ack, got %d", errInvalidPlaceOrderResp, len(data)),
		)
	}
	if data[0].SCode != "0" {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodPost,
			RequestPath: "/api/v5/trade/order",
			RequestID:   requestID,
			Code:        data[0].SCode,
			Message:     data[0].SMsg,
		}
	}
	return &data[0], nil
}
