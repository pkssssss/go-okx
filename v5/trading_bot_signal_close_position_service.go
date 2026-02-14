package okx

import (
	"context"
	"errors"
	"net/http"
)

type tradingBotSignalClosePositionRequest struct {
	AlgoId string `json:"algoId"`
	InstId string `json:"instId"`
}

// TradingBotSignalClosePositionService 市价仓位全平（平掉指定产品持仓）。
type TradingBotSignalClosePositionService struct {
	c *Client
	r tradingBotSignalClosePositionRequest
}

// NewTradingBotSignalClosePositionService 创建 TradingBotSignalClosePositionService。
func (c *Client) NewTradingBotSignalClosePositionService() *TradingBotSignalClosePositionService {
	return &TradingBotSignalClosePositionService{c: c}
}

func (s *TradingBotSignalClosePositionService) AlgoId(algoId string) *TradingBotSignalClosePositionService {
	s.r.AlgoId = algoId
	return s
}

func (s *TradingBotSignalClosePositionService) InstId(instId string) *TradingBotSignalClosePositionService {
	s.r.InstId = instId
	return s
}

var (
	errTradingBotSignalClosePositionMissingRequired = errors.New("okx: tradingBot signal close-position requires algoId and instId")
	errEmptyTradingBotSignalClosePositionResponse   = errors.New("okx: empty tradingBot signal close-position response")
)

// Do 市价仓位全平（POST /api/v5/tradingBot/signal/close-position）。
func (s *TradingBotSignalClosePositionService) Do(ctx context.Context) (*TradingBotAlgoIdAck, error) {
	if s.r.AlgoId == "" || s.r.InstId == "" {
		return nil, errTradingBotSignalClosePositionMissingRequired
	}

	var data []TradingBotAlgoIdAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/tradingBot/signal/close-position", nil, s.r, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/tradingBot/signal/close-position", requestID, errEmptyTradingBotSignalClosePositionResponse)
	}
	return &data[0], nil
}
