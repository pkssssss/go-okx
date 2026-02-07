package okx

import (
	"context"
	"errors"
	"net/http"
)

type tradingBotGridMarginBalanceRequest struct {
	AlgoId  string `json:"algoId"`
	Type    string `json:"type"`
	Amt     string `json:"amt,omitempty"`
	Percent string `json:"percent,omitempty"`
}

// TradingBotGridMarginBalanceService 调整保证金。
type TradingBotGridMarginBalanceService struct {
	c *Client
	r tradingBotGridMarginBalanceRequest
}

// NewTradingBotGridMarginBalanceService 创建 TradingBotGridMarginBalanceService。
func (c *Client) NewTradingBotGridMarginBalanceService() *TradingBotGridMarginBalanceService {
	return &TradingBotGridMarginBalanceService{c: c}
}

func (s *TradingBotGridMarginBalanceService) AlgoId(algoId string) *TradingBotGridMarginBalanceService {
	s.r.AlgoId = algoId
	return s
}

func (s *TradingBotGridMarginBalanceService) Type(typ string) *TradingBotGridMarginBalanceService {
	s.r.Type = typ
	return s
}

func (s *TradingBotGridMarginBalanceService) Amt(amt string) *TradingBotGridMarginBalanceService {
	s.r.Amt = amt
	return s
}

func (s *TradingBotGridMarginBalanceService) Percent(percent string) *TradingBotGridMarginBalanceService {
	s.r.Percent = percent
	return s
}

var (
	errTradingBotGridMarginBalanceMissingRequired = errors.New("okx: tradingBot grid margin-balance requires algoId and type")
	errTradingBotGridMarginBalanceMissingAmt      = errors.New("okx: tradingBot grid margin-balance requires amt or percent")
	errEmptyTradingBotGridMarginBalanceResponse   = errors.New("okx: empty tradingBot grid margin-balance response")
)

// Do 调整保证金（POST /api/v5/tradingBot/grid/margin-balance）。
func (s *TradingBotGridMarginBalanceService) Do(ctx context.Context) (*TradingBotOrderAck, error) {
	if s.r.AlgoId == "" || s.r.Type == "" {
		return nil, errTradingBotGridMarginBalanceMissingRequired
	}
	if s.r.Amt == "" && s.r.Percent == "" {
		return nil, errTradingBotGridMarginBalanceMissingAmt
	}

	var data []TradingBotOrderAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/tradingBot/grid/margin-balance", nil, s.r, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyTradingBotGridMarginBalanceResponse
	}
	if data[0].SCode != "" && data[0].SCode != "0" {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodPost,
			RequestPath: "/api/v5/tradingBot/grid/margin-balance",
			Code:        data[0].SCode,
			Message:     data[0].SMsg,
		}
	}
	return &data[0], nil
}
