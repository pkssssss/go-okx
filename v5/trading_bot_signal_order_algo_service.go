package okx

import (
	"context"
	"errors"
	"net/http"
)

type tradingBotSignalOrderAlgoRequest struct {
	SignalChanId string `json:"signalChanId"`

	IncludeAll *bool    `json:"includeAll,omitempty"`
	InstIds    []string `json:"instIds,omitempty"`

	Lever      string `json:"lever"`
	InvestAmt  string `json:"investAmt"`
	SubOrdType string `json:"subOrdType"`

	Ratio string `json:"ratio,omitempty"`

	EntrySettingParam *TradingBotSignalEntrySettingParam `json:"entrySettingParam,omitempty"`
	ExitSettingParam  *TradingBotSignalExitSettingParam  `json:"exitSettingParam,omitempty"`
}

// TradingBotSignalOrderAlgoService 创建信号策略。
type TradingBotSignalOrderAlgoService struct {
	c *Client
	r tradingBotSignalOrderAlgoRequest
}

// NewTradingBotSignalOrderAlgoService 创建 TradingBotSignalOrderAlgoService。
func (c *Client) NewTradingBotSignalOrderAlgoService() *TradingBotSignalOrderAlgoService {
	return &TradingBotSignalOrderAlgoService{c: c}
}

func (s *TradingBotSignalOrderAlgoService) SignalChanId(signalChanId string) *TradingBotSignalOrderAlgoService {
	s.r.SignalChanId = signalChanId
	return s
}

func (s *TradingBotSignalOrderAlgoService) IncludeAll(includeAll bool) *TradingBotSignalOrderAlgoService {
	s.r.IncludeAll = &includeAll
	return s
}

func (s *TradingBotSignalOrderAlgoService) InstIds(instIds []string) *TradingBotSignalOrderAlgoService {
	s.r.InstIds = instIds
	return s
}

func (s *TradingBotSignalOrderAlgoService) Lever(lever string) *TradingBotSignalOrderAlgoService {
	s.r.Lever = lever
	return s
}

func (s *TradingBotSignalOrderAlgoService) InvestAmt(investAmt string) *TradingBotSignalOrderAlgoService {
	s.r.InvestAmt = investAmt
	return s
}

func (s *TradingBotSignalOrderAlgoService) SubOrdType(subOrdType string) *TradingBotSignalOrderAlgoService {
	s.r.SubOrdType = subOrdType
	return s
}

func (s *TradingBotSignalOrderAlgoService) Ratio(ratio string) *TradingBotSignalOrderAlgoService {
	s.r.Ratio = ratio
	return s
}

func (s *TradingBotSignalOrderAlgoService) EntrySettingParam(p TradingBotSignalEntrySettingParam) *TradingBotSignalOrderAlgoService {
	s.r.EntrySettingParam = &p
	return s
}

func (s *TradingBotSignalOrderAlgoService) ExitSettingParam(p TradingBotSignalExitSettingParam) *TradingBotSignalOrderAlgoService {
	s.r.ExitSettingParam = &p
	return s
}

var (
	errTradingBotSignalOrderAlgoMissingRequired = errors.New("okx: tradingBot signal order-algo requires signalChanId, lever, investAmt and subOrdType")
	errTradingBotSignalOrderAlgoMissingInstIds  = errors.New("okx: tradingBot signal order-algo requires instIds when includeAll is false")
	errTradingBotSignalOrderAlgoMissingTpSlType = errors.New("okx: tradingBot signal order-algo exitSettingParam requires tpSlType")
	errEmptyTradingBotSignalOrderAlgoResponse   = errors.New("okx: empty tradingBot signal order-algo response")
)

// Do 创建信号策略（POST /api/v5/tradingBot/signal/order-algo）。
func (s *TradingBotSignalOrderAlgoService) Do(ctx context.Context) (*TradingBotOrderAck, error) {
	if s.r.SignalChanId == "" || s.r.Lever == "" || s.r.InvestAmt == "" || s.r.SubOrdType == "" {
		return nil, errTradingBotSignalOrderAlgoMissingRequired
	}

	includeAll := false
	if s.r.IncludeAll != nil {
		includeAll = *s.r.IncludeAll
	}
	if !includeAll && len(s.r.InstIds) == 0 {
		return nil, errTradingBotSignalOrderAlgoMissingInstIds
	}
	if s.r.ExitSettingParam != nil && s.r.ExitSettingParam.TpSlType == "" {
		return nil, errTradingBotSignalOrderAlgoMissingTpSlType
	}

	var data []TradingBotOrderAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/tradingBot/signal/order-algo", nil, s.r, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyTradingBotSignalOrderAlgoResponse
	}
	if data[0].SCode != "" && data[0].SCode != "0" {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodPost,
			RequestPath: "/api/v5/tradingBot/signal/order-algo",
			RequestID:   requestID,
			Code:        data[0].SCode,
			Message:     data[0].SMsg,
		}
	}
	return &data[0], nil
}
