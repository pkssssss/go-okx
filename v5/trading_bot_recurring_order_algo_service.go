package okx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type tradingBotRecurringOrderAlgoRequest struct {
	StgyName string `json:"stgyName"`

	RecurringList []TradingBotRecurringListItem `json:"recurringList"`

	Period        string `json:"period"`
	RecurringDay  string `json:"recurringDay,omitempty"`
	RecurringHour string `json:"recurringHour,omitempty"`
	RecurringTime string `json:"recurringTime"`
	TimeZone      string `json:"timeZone"`

	Amt           string `json:"amt"`
	InvestmentCcy string `json:"investmentCcy"`
	TdMode        string `json:"tdMode"`

	AlgoClOrdId   string `json:"algoClOrdId,omitempty"`
	Tag           string `json:"tag,omitempty"`
	TradeQuoteCcy string `json:"tradeQuoteCcy,omitempty"`
}

// TradingBotRecurringOrderAlgoService 定投策略委托下单。
type TradingBotRecurringOrderAlgoService struct {
	c *Client
	r tradingBotRecurringOrderAlgoRequest
}

// NewTradingBotRecurringOrderAlgoService 创建 TradingBotRecurringOrderAlgoService。
func (c *Client) NewTradingBotRecurringOrderAlgoService() *TradingBotRecurringOrderAlgoService {
	return &TradingBotRecurringOrderAlgoService{c: c}
}

func (s *TradingBotRecurringOrderAlgoService) StgyName(stgyName string) *TradingBotRecurringOrderAlgoService {
	s.r.StgyName = stgyName
	return s
}

func (s *TradingBotRecurringOrderAlgoService) RecurringList(list []TradingBotRecurringListItem) *TradingBotRecurringOrderAlgoService {
	s.r.RecurringList = list
	return s
}

func (s *TradingBotRecurringOrderAlgoService) Period(period string) *TradingBotRecurringOrderAlgoService {
	s.r.Period = period
	return s
}

func (s *TradingBotRecurringOrderAlgoService) RecurringDay(recurringDay string) *TradingBotRecurringOrderAlgoService {
	s.r.RecurringDay = recurringDay
	return s
}

func (s *TradingBotRecurringOrderAlgoService) RecurringHour(recurringHour string) *TradingBotRecurringOrderAlgoService {
	s.r.RecurringHour = recurringHour
	return s
}

func (s *TradingBotRecurringOrderAlgoService) RecurringTime(recurringTime string) *TradingBotRecurringOrderAlgoService {
	s.r.RecurringTime = recurringTime
	return s
}

func (s *TradingBotRecurringOrderAlgoService) TimeZone(timeZone string) *TradingBotRecurringOrderAlgoService {
	s.r.TimeZone = timeZone
	return s
}

func (s *TradingBotRecurringOrderAlgoService) Amt(amt string) *TradingBotRecurringOrderAlgoService {
	s.r.Amt = amt
	return s
}

func (s *TradingBotRecurringOrderAlgoService) InvestmentCcy(investmentCcy string) *TradingBotRecurringOrderAlgoService {
	s.r.InvestmentCcy = investmentCcy
	return s
}

func (s *TradingBotRecurringOrderAlgoService) TdMode(tdMode string) *TradingBotRecurringOrderAlgoService {
	s.r.TdMode = tdMode
	return s
}

func (s *TradingBotRecurringOrderAlgoService) AlgoClOrdId(algoClOrdId string) *TradingBotRecurringOrderAlgoService {
	s.r.AlgoClOrdId = algoClOrdId
	return s
}

func (s *TradingBotRecurringOrderAlgoService) Tag(tag string) *TradingBotRecurringOrderAlgoService {
	s.r.Tag = tag
	return s
}

func (s *TradingBotRecurringOrderAlgoService) TradeQuoteCcy(tradeQuoteCcy string) *TradingBotRecurringOrderAlgoService {
	s.r.TradeQuoteCcy = tradeQuoteCcy
	return s
}

var (
	errTradingBotRecurringOrderAlgoMissingRequired = errors.New("okx: tradingBot recurring order-algo requires stgyName, recurringList, period, recurringTime, timeZone, amt, investmentCcy and tdMode")
	errTradingBotRecurringOrderAlgoInvalidList     = errors.New("okx: tradingBot recurring order-algo invalid recurringList")
	errEmptyTradingBotRecurringOrderAlgoResponse   = errors.New("okx: empty tradingBot recurring order-algo response")
	errInvalidTradingBotRecurringOrderAlgoResponse = errors.New("okx: invalid tradingBot recurring order-algo response")
)

// Do 定投策略委托下单（POST /api/v5/tradingBot/recurring/order-algo）。
func (s *TradingBotRecurringOrderAlgoService) Do(ctx context.Context) (*TradingBotOrderAck, error) {
	if s.r.StgyName == "" || s.r.Period == "" || s.r.RecurringTime == "" || s.r.TimeZone == "" || s.r.Amt == "" || s.r.InvestmentCcy == "" || s.r.TdMode == "" || len(s.r.RecurringList) == 0 {
		return nil, errTradingBotRecurringOrderAlgoMissingRequired
	}
	for _, it := range s.r.RecurringList {
		if it.Ccy == "" || it.Ratio == "" {
			return nil, errTradingBotRecurringOrderAlgoInvalidList
		}
	}

	var data []TradingBotOrderAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/tradingBot/recurring/order-algo", nil, s.r, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/tradingBot/recurring/order-algo", requestID, errEmptyTradingBotRecurringOrderAlgoResponse)
	}
	if len(data) != 1 {
		return nil, newInvalidDataAPIError(
			http.MethodPost,
			"/api/v5/tradingBot/recurring/order-algo",
			requestID,
			fmt.Errorf("%w: expected 1 ack, got %d", errInvalidTradingBotRecurringOrderAlgoResponse, len(data)),
		)
	}
	if data[0].SCode != "0" {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodPost,
			RequestPath: "/api/v5/tradingBot/recurring/order-algo",
			RequestID:   requestID,
			Code:        data[0].SCode,
			Message:     data[0].SMsg,
		}
	}
	return &data[0], nil
}
