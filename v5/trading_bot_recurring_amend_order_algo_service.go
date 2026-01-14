package okx

import (
	"context"
	"errors"
	"net/http"
)

type tradingBotRecurringAmendOrderAlgoRequest struct {
	AlgoId   string `json:"algoId"`
	StgyName string `json:"stgyName"`
}

// TradingBotRecurringAmendOrderAlgoService 修改定投策略订单。
type TradingBotRecurringAmendOrderAlgoService struct {
	c *Client
	r tradingBotRecurringAmendOrderAlgoRequest
}

// NewTradingBotRecurringAmendOrderAlgoService 创建 TradingBotRecurringAmendOrderAlgoService。
func (c *Client) NewTradingBotRecurringAmendOrderAlgoService() *TradingBotRecurringAmendOrderAlgoService {
	return &TradingBotRecurringAmendOrderAlgoService{c: c}
}

func (s *TradingBotRecurringAmendOrderAlgoService) AlgoId(algoId string) *TradingBotRecurringAmendOrderAlgoService {
	s.r.AlgoId = algoId
	return s
}

func (s *TradingBotRecurringAmendOrderAlgoService) StgyName(stgyName string) *TradingBotRecurringAmendOrderAlgoService {
	s.r.StgyName = stgyName
	return s
}

var (
	errTradingBotRecurringAmendOrderAlgoMissingRequired = errors.New("okx: tradingBot recurring amend-order-algo requires algoId and stgyName")
	errEmptyTradingBotRecurringAmendOrderAlgoResponse   = errors.New("okx: empty tradingBot recurring amend-order-algo response")
)

// Do 修改定投策略订单（POST /api/v5/tradingBot/recurring/amend-order-algo）。
func (s *TradingBotRecurringAmendOrderAlgoService) Do(ctx context.Context) (*TradingBotOrderAck, error) {
	if s.r.AlgoId == "" || s.r.StgyName == "" {
		return nil, errTradingBotRecurringAmendOrderAlgoMissingRequired
	}

	var data []TradingBotOrderAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/tradingBot/recurring/amend-order-algo", nil, s.r, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyTradingBotRecurringAmendOrderAlgoResponse
	}
	if data[0].SCode != "" && data[0].SCode != "0" {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodPost,
			RequestPath: "/api/v5/tradingBot/recurring/amend-order-algo",
			Code:        data[0].SCode,
			Message:     data[0].SMsg,
		}
	}
	return &data[0], nil
}
