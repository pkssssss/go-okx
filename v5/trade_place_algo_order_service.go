package okx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type placeAlgoOrderRequest struct {
	InstId  string `json:"instId"`
	TdMode  string `json:"tdMode"`
	Side    string `json:"side"`
	OrdType string `json:"ordType"`

	PosSide string `json:"posSide,omitempty"`
	Ccy     string `json:"ccy,omitempty"`
	Tag     string `json:"tag,omitempty"`

	AlgoClOrdId string `json:"algoClOrdId,omitempty"`

	Sz            string `json:"sz,omitempty"`
	CloseFraction string `json:"closeFraction,omitempty"`

	ReduceOnly    *bool  `json:"reduceOnly,omitempty"`
	TgtCcy        string `json:"tgtCcy,omitempty"`
	TradeQuoteCcy string `json:"tradeQuoteCcy,omitempty"`

	CxlOnClosePos *bool `json:"cxlOnClosePos,omitempty"`

	TpTriggerPx     string `json:"tpTriggerPx,omitempty"`
	TpTriggerPxType string `json:"tpTriggerPxType,omitempty"`
	TpOrdPx         string `json:"tpOrdPx,omitempty"`
	TpOrdKind       string `json:"tpOrdKind,omitempty"`

	SlTriggerPx     string `json:"slTriggerPx,omitempty"`
	SlTriggerPxType string `json:"slTriggerPxType,omitempty"`
	SlOrdPx         string `json:"slOrdPx,omitempty"`
	SlOrdKind       string `json:"slOrdKind,omitempty"`

	TriggerPx     string `json:"triggerPx,omitempty"`
	TriggerPxType string `json:"triggerPxType,omitempty"`
	OrderPx       string `json:"orderPx,omitempty"`

	CallbackRatio  string `json:"callbackRatio,omitempty"`
	CallbackSpread string `json:"callbackSpread,omitempty"`
	ActivePx       string `json:"activePx,omitempty"`

	PxLimit      string `json:"pxLimit,omitempty"`
	SzLimit      string `json:"szLimit,omitempty"`
	TimeInterval string `json:"timeInterval,omitempty"`
	PxSpread     string `json:"pxSpread,omitempty"`

	AttachAlgoOrds []TradeAlgoOrderAttach `json:"attachAlgoOrds,omitempty"`
}

// PlaceAlgoOrderService 策略委托下单。
type PlaceAlgoOrderService struct {
	c   *Client
	req placeAlgoOrderRequest
}

// NewPlaceAlgoOrderService 创建 PlaceAlgoOrderService。
func (c *Client) NewPlaceAlgoOrderService() *PlaceAlgoOrderService {
	return &PlaceAlgoOrderService{c: c}
}

func (s *PlaceAlgoOrderService) InstId(instId string) *PlaceAlgoOrderService {
	s.req.InstId = instId
	return s
}

func (s *PlaceAlgoOrderService) TdMode(tdMode string) *PlaceAlgoOrderService {
	s.req.TdMode = tdMode
	return s
}

func (s *PlaceAlgoOrderService) Side(side string) *PlaceAlgoOrderService {
	s.req.Side = side
	return s
}

// OrdType 设置订单类型（必填，例如 conditional/oco/trigger/move_order_stop/twap）。
func (s *PlaceAlgoOrderService) OrdType(ordType string) *PlaceAlgoOrderService {
	s.req.OrdType = ordType
	return s
}

func (s *PlaceAlgoOrderService) PosSide(posSide string) *PlaceAlgoOrderService {
	s.req.PosSide = posSide
	return s
}

func (s *PlaceAlgoOrderService) Ccy(ccy string) *PlaceAlgoOrderService {
	s.req.Ccy = ccy
	return s
}

func (s *PlaceAlgoOrderService) Tag(tag string) *PlaceAlgoOrderService {
	s.req.Tag = tag
	return s
}

// AlgoClOrdId 设置客户自定义策略订单ID（可选）。
func (s *PlaceAlgoOrderService) AlgoClOrdId(algoClOrdId string) *PlaceAlgoOrderService {
	s.req.AlgoClOrdId = algoClOrdId
	return s
}

// Sz 设置委托数量（与 CloseFraction 二选一）。
func (s *PlaceAlgoOrderService) Sz(sz string) *PlaceAlgoOrderService {
	s.req.Sz = sz
	return s
}

// CloseFraction 设置平仓百分比（与 Sz 二选一；1 表示 100%）。
func (s *PlaceAlgoOrderService) CloseFraction(closeFraction string) *PlaceAlgoOrderService {
	s.req.CloseFraction = closeFraction
	return s
}

func (s *PlaceAlgoOrderService) ReduceOnly(enable bool) *PlaceAlgoOrderService {
	s.req.ReduceOnly = &enable
	return s
}

func (s *PlaceAlgoOrderService) TgtCcy(tgtCcy string) *PlaceAlgoOrderService {
	s.req.TgtCcy = tgtCcy
	return s
}

func (s *PlaceAlgoOrderService) TradeQuoteCcy(tradeQuoteCcy string) *PlaceAlgoOrderService {
	s.req.TradeQuoteCcy = tradeQuoteCcy
	return s
}

// CxlOnClosePos 设置止盈止损订单是否与仓位关联（可选；传 true 时要求 reduceOnly=true）。
func (s *PlaceAlgoOrderService) CxlOnClosePos(enable bool) *PlaceAlgoOrderService {
	s.req.CxlOnClosePos = &enable
	return s
}

func (s *PlaceAlgoOrderService) TpTriggerPx(px string) *PlaceAlgoOrderService {
	s.req.TpTriggerPx = px
	return s
}

func (s *PlaceAlgoOrderService) TpTriggerPxType(tpTriggerPxType string) *PlaceAlgoOrderService {
	s.req.TpTriggerPxType = tpTriggerPxType
	return s
}

func (s *PlaceAlgoOrderService) TpOrdPx(px string) *PlaceAlgoOrderService {
	s.req.TpOrdPx = px
	return s
}

func (s *PlaceAlgoOrderService) TpOrdKind(kind string) *PlaceAlgoOrderService {
	s.req.TpOrdKind = kind
	return s
}

func (s *PlaceAlgoOrderService) SlTriggerPx(px string) *PlaceAlgoOrderService {
	s.req.SlTriggerPx = px
	return s
}

func (s *PlaceAlgoOrderService) SlTriggerPxType(slTriggerPxType string) *PlaceAlgoOrderService {
	s.req.SlTriggerPxType = slTriggerPxType
	return s
}

func (s *PlaceAlgoOrderService) SlOrdPx(px string) *PlaceAlgoOrderService {
	s.req.SlOrdPx = px
	return s
}

func (s *PlaceAlgoOrderService) SlOrdKind(kind string) *PlaceAlgoOrderService {
	s.req.SlOrdKind = kind
	return s
}

// TriggerPx 设置计划委托触发价（ordType=trigger 等）。
func (s *PlaceAlgoOrderService) TriggerPx(px string) *PlaceAlgoOrderService {
	s.req.TriggerPx = px
	return s
}

func (s *PlaceAlgoOrderService) TriggerPxType(typ string) *PlaceAlgoOrderService {
	s.req.TriggerPxType = typ
	return s
}

// OrderPx 设置计划委托委托价（ordType=trigger 等；-1 表示市价）。
func (s *PlaceAlgoOrderService) OrderPx(px string) *PlaceAlgoOrderService {
	s.req.OrderPx = px
	return s
}

// CallbackRatio 设置移动止盈止损回调比例（ordType=move_order_stop）。
func (s *PlaceAlgoOrderService) CallbackRatio(callbackRatio string) *PlaceAlgoOrderService {
	s.req.CallbackRatio = callbackRatio
	return s
}

func (s *PlaceAlgoOrderService) CallbackSpread(callbackSpread string) *PlaceAlgoOrderService {
	s.req.CallbackSpread = callbackSpread
	return s
}

func (s *PlaceAlgoOrderService) ActivePx(activePx string) *PlaceAlgoOrderService {
	s.req.ActivePx = activePx
	return s
}

// PxLimit/SzLimit/TimeInterval/PxSpread 主要用于时间加权委托（ordType=twap）。
func (s *PlaceAlgoOrderService) PxLimit(pxLimit string) *PlaceAlgoOrderService {
	s.req.PxLimit = pxLimit
	return s
}

func (s *PlaceAlgoOrderService) SzLimit(szLimit string) *PlaceAlgoOrderService {
	s.req.SzLimit = szLimit
	return s
}

func (s *PlaceAlgoOrderService) TimeInterval(timeInterval string) *PlaceAlgoOrderService {
	s.req.TimeInterval = timeInterval
	return s
}

func (s *PlaceAlgoOrderService) PxSpread(pxSpread string) *PlaceAlgoOrderService {
	s.req.PxSpread = pxSpread
	return s
}

// AttachAlgoOrds 设置附带止盈止损信息（适用于计划委托等场景）。
func (s *PlaceAlgoOrderService) AttachAlgoOrds(attach []TradeAlgoOrderAttach) *PlaceAlgoOrderService {
	s.req.AttachAlgoOrds = attach
	return s
}

var (
	errPlaceAlgoOrderMissingInstId              = errors.New("okx: place algo order requires instId")
	errPlaceAlgoOrderMissingTdMode              = errors.New("okx: place algo order requires tdMode")
	errPlaceAlgoOrderMissingSide                = errors.New("okx: place algo order requires side")
	errPlaceAlgoOrderMissingOrdType             = errors.New("okx: place algo order requires ordType")
	errPlaceAlgoOrderMissingSzOrCloseFraction   = errors.New("okx: place algo order requires sz or closeFraction")
	errPlaceAlgoOrderSzAndCloseFractionConflict = errors.New("okx: place algo order requires at most one of sz/closeFraction")
	errEmptyPlaceAlgoOrderResponse              = errors.New("okx: empty place algo order response")
	errInvalidPlaceAlgoOrderResponse            = errors.New("okx: invalid place algo order response")
)

// Do 策略委托下单（POST /api/v5/trade/order-algo）。
func (s *PlaceAlgoOrderService) Do(ctx context.Context) (*TradeAlgoOrderAck, error) {
	if s.req.InstId == "" {
		return nil, errPlaceAlgoOrderMissingInstId
	}
	if s.req.TdMode == "" {
		return nil, errPlaceAlgoOrderMissingTdMode
	}
	if s.req.Side == "" {
		return nil, errPlaceAlgoOrderMissingSide
	}
	if s.req.OrdType == "" {
		return nil, errPlaceAlgoOrderMissingOrdType
	}
	if s.req.Sz == "" && s.req.CloseFraction == "" {
		return nil, errPlaceAlgoOrderMissingSzOrCloseFraction
	}
	if s.req.Sz != "" && s.req.CloseFraction != "" {
		return nil, errPlaceAlgoOrderSzAndCloseFractionConflict
	}

	var data []TradeAlgoOrderAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/trade/order-algo", nil, s.req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/trade/order-algo", requestID, errEmptyPlaceAlgoOrderResponse)
	}
	if len(data) != 1 {
		return nil, newInvalidDataAPIError(
			http.MethodPost,
			"/api/v5/trade/order-algo",
			requestID,
			fmt.Errorf("%w: expected 1 ack, got %d", errInvalidPlaceAlgoOrderResponse, len(data)),
		)
	}
	if data[0].SCode != "0" {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodPost,
			RequestPath: "/api/v5/trade/order-algo",
			RequestID:   requestID,
			Code:        data[0].SCode,
			Message:     data[0].SMsg,
		}
	}
	return &data[0], nil
}
