package okx

import (
	"context"
	"errors"
	"net/http"
)

type tradingBotSignalMarginBalanceRequest struct {
	AlgoId string `json:"algoId"`
	Type   string `json:"type"`
	Amt    string `json:"amt"`

	AllowReinvest *bool `json:"allowReinvest,omitempty"`
}

// TradingBotSignalMarginBalanceService 调整保证金。
type TradingBotSignalMarginBalanceService struct {
	c *Client
	r tradingBotSignalMarginBalanceRequest
}

// NewTradingBotSignalMarginBalanceService 创建 TradingBotSignalMarginBalanceService。
func (c *Client) NewTradingBotSignalMarginBalanceService() *TradingBotSignalMarginBalanceService {
	return &TradingBotSignalMarginBalanceService{c: c}
}

func (s *TradingBotSignalMarginBalanceService) AlgoId(algoId string) *TradingBotSignalMarginBalanceService {
	s.r.AlgoId = algoId
	return s
}

func (s *TradingBotSignalMarginBalanceService) Type(typ string) *TradingBotSignalMarginBalanceService {
	s.r.Type = typ
	return s
}

func (s *TradingBotSignalMarginBalanceService) Amt(amt string) *TradingBotSignalMarginBalanceService {
	s.r.Amt = amt
	return s
}

func (s *TradingBotSignalMarginBalanceService) AllowReinvest(allowReinvest bool) *TradingBotSignalMarginBalanceService {
	s.r.AllowReinvest = &allowReinvest
	return s
}

var (
	errTradingBotSignalMarginBalanceMissingRequired = errors.New("okx: tradingBot signal margin-balance requires algoId, type and amt")
	errEmptyTradingBotSignalMarginBalanceResponse   = errors.New("okx: empty tradingBot signal margin-balance response")
)

// Do 调整保证金（POST /api/v5/tradingBot/signal/margin-balance）。
func (s *TradingBotSignalMarginBalanceService) Do(ctx context.Context) (*TradingBotAlgoIdAck, error) {
	if s.r.AlgoId == "" || s.r.Type == "" || s.r.Amt == "" {
		return nil, errTradingBotSignalMarginBalanceMissingRequired
	}

	var data []TradingBotAlgoIdAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/tradingBot/signal/margin-balance", nil, s.r, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyTradingBotSignalMarginBalanceResponse
	}
	return &data[0], nil
}
