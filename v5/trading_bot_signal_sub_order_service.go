package okx

import (
	"context"
	"errors"
	"net/http"
)

type tradingBotSignalSubOrderRequest struct {
	InstId string `json:"instId"`
	AlgoId string `json:"algoId"`

	Side    string `json:"side"`
	OrdType string `json:"ordType"`
	Sz      string `json:"sz"`
	Px      string `json:"px,omitempty"`

	ReduceOnly *bool `json:"reduceOnly,omitempty"`
}

// TradingBotSignalSubOrderService 信号策略下单（创建子订单）。
type TradingBotSignalSubOrderService struct {
	c *Client
	r tradingBotSignalSubOrderRequest
}

// NewTradingBotSignalSubOrderService 创建 TradingBotSignalSubOrderService。
func (c *Client) NewTradingBotSignalSubOrderService() *TradingBotSignalSubOrderService {
	return &TradingBotSignalSubOrderService{c: c}
}

func (s *TradingBotSignalSubOrderService) InstId(instId string) *TradingBotSignalSubOrderService {
	s.r.InstId = instId
	return s
}

func (s *TradingBotSignalSubOrderService) AlgoId(algoId string) *TradingBotSignalSubOrderService {
	s.r.AlgoId = algoId
	return s
}

func (s *TradingBotSignalSubOrderService) Side(side string) *TradingBotSignalSubOrderService {
	s.r.Side = side
	return s
}

func (s *TradingBotSignalSubOrderService) OrdType(ordType string) *TradingBotSignalSubOrderService {
	s.r.OrdType = ordType
	return s
}

func (s *TradingBotSignalSubOrderService) Sz(sz string) *TradingBotSignalSubOrderService {
	s.r.Sz = sz
	return s
}

func (s *TradingBotSignalSubOrderService) Px(px string) *TradingBotSignalSubOrderService {
	s.r.Px = px
	return s
}

func (s *TradingBotSignalSubOrderService) ReduceOnly(reduceOnly bool) *TradingBotSignalSubOrderService {
	s.r.ReduceOnly = &reduceOnly
	return s
}

var (
	errTradingBotSignalSubOrderMissingRequired = errors.New("okx: tradingBot signal sub-order requires algoId, instId, side, ordType and sz")
)

// Do 下单（POST /api/v5/tradingBot/signal/sub-order）。
func (s *TradingBotSignalSubOrderService) Do(ctx context.Context) error {
	if s.r.AlgoId == "" || s.r.InstId == "" || s.r.Side == "" || s.r.OrdType == "" || s.r.Sz == "" {
		return errTradingBotSignalSubOrderMissingRequired
	}

	return s.c.do(ctx, http.MethodPost, "/api/v5/tradingBot/signal/sub-order", nil, s.r, true, nil)
}
