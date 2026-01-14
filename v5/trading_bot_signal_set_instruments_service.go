package okx

import (
	"context"
	"errors"
	"net/http"
)

type tradingBotSignalSetInstrumentsRequest struct {
	AlgoId     string   `json:"algoId"`
	InstIds    []string `json:"instIds,omitempty"`
	IncludeAll *bool    `json:"includeAll"`
}

// TradingBotSignalSetInstrumentsService 设置币对。
type TradingBotSignalSetInstrumentsService struct {
	c *Client
	r tradingBotSignalSetInstrumentsRequest
}

// NewTradingBotSignalSetInstrumentsService 创建 TradingBotSignalSetInstrumentsService。
func (c *Client) NewTradingBotSignalSetInstrumentsService() *TradingBotSignalSetInstrumentsService {
	return &TradingBotSignalSetInstrumentsService{c: c}
}

func (s *TradingBotSignalSetInstrumentsService) AlgoId(algoId string) *TradingBotSignalSetInstrumentsService {
	s.r.AlgoId = algoId
	return s
}

func (s *TradingBotSignalSetInstrumentsService) InstIds(instIds []string) *TradingBotSignalSetInstrumentsService {
	s.r.InstIds = instIds
	return s
}

func (s *TradingBotSignalSetInstrumentsService) IncludeAll(includeAll bool) *TradingBotSignalSetInstrumentsService {
	s.r.IncludeAll = &includeAll
	return s
}

var (
	errTradingBotSignalSetInstrumentsMissingRequired = errors.New("okx: tradingBot signal set-instruments requires algoId and includeAll")
	errTradingBotSignalSetInstrumentsMissingInstIds  = errors.New("okx: tradingBot signal set-instruments requires instIds when includeAll is false")
	errEmptyTradingBotSignalSetInstrumentsResponse   = errors.New("okx: empty tradingBot signal set-instruments response")
)

// Do 设置币对（POST /api/v5/tradingBot/signal/set-instruments）。
func (s *TradingBotSignalSetInstrumentsService) Do(ctx context.Context) (*TradingBotAlgoIdAck, error) {
	if s.r.AlgoId == "" || s.r.IncludeAll == nil {
		return nil, errTradingBotSignalSetInstrumentsMissingRequired
	}
	if !*s.r.IncludeAll && len(s.r.InstIds) == 0 {
		return nil, errTradingBotSignalSetInstrumentsMissingInstIds
	}

	var data []TradingBotAlgoIdAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/tradingBot/signal/set-instruments", nil, s.r, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyTradingBotSignalSetInstrumentsResponse
	}
	return &data[0], nil
}
