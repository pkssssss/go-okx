package okx

import (
	"context"
	"errors"
	"net/http"
)

type tradingBotGridWithdrawIncomeRequest struct {
	AlgoId string `json:"algoId"`
}

// TradingBotGridWithdrawIncomeService 现货网格提取利润。
type TradingBotGridWithdrawIncomeService struct {
	c *Client
	r tradingBotGridWithdrawIncomeRequest
}

// NewTradingBotGridWithdrawIncomeService 创建 TradingBotGridWithdrawIncomeService。
func (c *Client) NewTradingBotGridWithdrawIncomeService() *TradingBotGridWithdrawIncomeService {
	return &TradingBotGridWithdrawIncomeService{c: c}
}

func (s *TradingBotGridWithdrawIncomeService) AlgoId(algoId string) *TradingBotGridWithdrawIncomeService {
	s.r.AlgoId = algoId
	return s
}

var (
	errTradingBotGridWithdrawIncomeMissingAlgoId = errors.New("okx: tradingBot grid withdraw-income requires algoId")
	errEmptyTradingBotGridWithdrawIncomeResponse = errors.New("okx: empty tradingBot grid withdraw-income response")
)

// Do 现货网格提取利润（POST /api/v5/tradingBot/grid/withdraw-income）。
func (s *TradingBotGridWithdrawIncomeService) Do(ctx context.Context) (*TradingBotGridWithdrawIncomeAck, error) {
	if s.r.AlgoId == "" {
		return nil, errTradingBotGridWithdrawIncomeMissingAlgoId
	}

	var data []TradingBotGridWithdrawIncomeAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/tradingBot/grid/withdraw-income", nil, s.r, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyTradingBotGridWithdrawIncomeResponse
	}
	return &data[0], nil
}
