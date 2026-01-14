package okx

import (
	"context"
	"errors"
	"net/http"
)

type tradingBotGridClosePositionRequest struct {
	AlgoId   string `json:"algoId"`
	MktClose *bool  `json:"mktClose"`
	Sz       string `json:"sz,omitempty"`
	Px       string `json:"px,omitempty"`
}

// TradingBotGridClosePositionService 合约网格平仓。
type TradingBotGridClosePositionService struct {
	c *Client
	r tradingBotGridClosePositionRequest
}

// NewTradingBotGridClosePositionService 创建 TradingBotGridClosePositionService。
func (c *Client) NewTradingBotGridClosePositionService() *TradingBotGridClosePositionService {
	return &TradingBotGridClosePositionService{c: c}
}

func (s *TradingBotGridClosePositionService) AlgoId(algoId string) *TradingBotGridClosePositionService {
	s.r.AlgoId = algoId
	return s
}

func (s *TradingBotGridClosePositionService) MktClose(mktClose bool) *TradingBotGridClosePositionService {
	s.r.MktClose = &mktClose
	return s
}

func (s *TradingBotGridClosePositionService) Sz(sz string) *TradingBotGridClosePositionService {
	s.r.Sz = sz
	return s
}

func (s *TradingBotGridClosePositionService) Px(px string) *TradingBotGridClosePositionService {
	s.r.Px = px
	return s
}

var (
	errTradingBotGridClosePositionMissingRequired = errors.New("okx: tradingBot grid close-position requires algoId and mktClose")
	errTradingBotGridClosePositionMissingSzPx     = errors.New("okx: tradingBot grid close-position partial close requires sz and px")
	errEmptyTradingBotGridClosePositionResponse   = errors.New("okx: empty tradingBot grid close-position response")
)

// Do 合约网格平仓（POST /api/v5/tradingBot/grid/close-position）。
func (s *TradingBotGridClosePositionService) Do(ctx context.Context) (*TradingBotGridCloseOrderAck, error) {
	if s.r.AlgoId == "" || s.r.MktClose == nil {
		return nil, errTradingBotGridClosePositionMissingRequired
	}
	if !*s.r.MktClose && (s.r.Sz == "" || s.r.Px == "") {
		return nil, errTradingBotGridClosePositionMissingSzPx
	}

	var data []TradingBotGridCloseOrderAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/tradingBot/grid/close-position", nil, s.r, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyTradingBotGridClosePositionResponse
	}
	return &data[0], nil
}
