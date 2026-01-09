package okx

import (
	"context"
	"errors"
	"net/http"
)

type orderPrecheckRequest struct {
	InstId  string `json:"instId"`
	TdMode  string `json:"tdMode"`
	Side    string `json:"side"`
	OrdType string `json:"ordType"`
	Sz      string `json:"sz"`

	ClOrdId string `json:"clOrdId,omitempty"`
	PosSide string `json:"posSide,omitempty"`
	Px      string `json:"px,omitempty"`

	ReduceOnly *bool  `json:"reduceOnly,omitempty"`
	TgtCcy     string `json:"tgtCcy,omitempty"`

	AttachAlgoOrds []TradeAlgoOrderAttach `json:"attachAlgoOrds,omitempty"`
}

// OrderPrecheckService 订单预检查（仅适用于跨币种保证金模式和组合保证金模式）。
type OrderPrecheckService struct {
	c   *Client
	req orderPrecheckRequest
}

// NewOrderPrecheckService 创建 OrderPrecheckService。
func (c *Client) NewOrderPrecheckService() *OrderPrecheckService {
	return &OrderPrecheckService{c: c}
}

func (s *OrderPrecheckService) InstId(instId string) *OrderPrecheckService {
	s.req.InstId = instId
	return s
}

func (s *OrderPrecheckService) TdMode(tdMode string) *OrderPrecheckService {
	s.req.TdMode = tdMode
	return s
}

func (s *OrderPrecheckService) Side(side string) *OrderPrecheckService {
	s.req.Side = side
	return s
}

func (s *OrderPrecheckService) OrdType(ordType string) *OrderPrecheckService {
	s.req.OrdType = ordType
	return s
}

func (s *OrderPrecheckService) Sz(sz string) *OrderPrecheckService {
	s.req.Sz = sz
	return s
}

func (s *OrderPrecheckService) Px(px string) *OrderPrecheckService {
	s.req.Px = px
	return s
}

func (s *OrderPrecheckService) ClOrdId(clOrdId string) *OrderPrecheckService {
	s.req.ClOrdId = clOrdId
	return s
}

func (s *OrderPrecheckService) PosSide(posSide string) *OrderPrecheckService {
	s.req.PosSide = posSide
	return s
}

func (s *OrderPrecheckService) ReduceOnly(enable bool) *OrderPrecheckService {
	s.req.ReduceOnly = &enable
	return s
}

func (s *OrderPrecheckService) TgtCcy(tgtCcy string) *OrderPrecheckService {
	s.req.TgtCcy = tgtCcy
	return s
}

func (s *OrderPrecheckService) AttachAlgoOrds(attach []TradeAlgoOrderAttach) *OrderPrecheckService {
	s.req.AttachAlgoOrds = attach
	return s
}

var (
	errOrderPrecheckMissingInstId  = errors.New("okx: order precheck requires instId")
	errOrderPrecheckMissingTdMode  = errors.New("okx: order precheck requires tdMode")
	errOrderPrecheckMissingSide    = errors.New("okx: order precheck requires side")
	errOrderPrecheckMissingOrdType = errors.New("okx: order precheck requires ordType")
	errOrderPrecheckMissingSz      = errors.New("okx: order precheck requires sz")
	errOrderPrecheckMissingPx      = errors.New("okx: order precheck requires px for this ordType")
)

func orderPrecheckRequiresPrice(ordType string) bool {
	switch ordType {
	case "limit", "post_only", "fok", "ioc":
		return true
	default:
		return false
	}
}

// Do 订单预检查（POST /api/v5/trade/order-precheck）。
func (s *OrderPrecheckService) Do(ctx context.Context) ([]OrderPrecheckResult, error) {
	if s.req.InstId == "" {
		return nil, errOrderPrecheckMissingInstId
	}
	if s.req.TdMode == "" {
		return nil, errOrderPrecheckMissingTdMode
	}
	if s.req.Side == "" {
		return nil, errOrderPrecheckMissingSide
	}
	if s.req.OrdType == "" {
		return nil, errOrderPrecheckMissingOrdType
	}
	if s.req.Sz == "" {
		return nil, errOrderPrecheckMissingSz
	}
	if orderPrecheckRequiresPrice(s.req.OrdType) && s.req.Px == "" {
		return nil, errOrderPrecheckMissingPx
	}

	var data []OrderPrecheckResult
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/trade/order-precheck", nil, s.req, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
