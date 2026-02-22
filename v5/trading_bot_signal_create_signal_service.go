package okx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type tradingBotSignalCreateSignalRequest struct {
	SignalChanName string `json:"signalChanName"`
	SignalChanDesc string `json:"signalChanDesc,omitempty"`
}

// TradingBotSignalCreateSignalService 创建信号。
type TradingBotSignalCreateSignalService struct {
	c *Client
	r tradingBotSignalCreateSignalRequest
}

// NewTradingBotSignalCreateSignalService 创建 TradingBotSignalCreateSignalService。
func (c *Client) NewTradingBotSignalCreateSignalService() *TradingBotSignalCreateSignalService {
	return &TradingBotSignalCreateSignalService{c: c}
}

func (s *TradingBotSignalCreateSignalService) SignalChanName(name string) *TradingBotSignalCreateSignalService {
	s.r.SignalChanName = name
	return s
}

func (s *TradingBotSignalCreateSignalService) SignalChanDesc(desc string) *TradingBotSignalCreateSignalService {
	s.r.SignalChanDesc = desc
	return s
}

var (
	errTradingBotSignalCreateSignalMissingName     = errors.New("okx: tradingBot signal create-signal requires signalChanName")
	errEmptyTradingBotSignalCreateSignalResponse   = errors.New("okx: empty tradingBot signal create-signal response")
	errInvalidTradingBotSignalCreateSignalResponse = errors.New("okx: invalid tradingBot signal create-signal response")
)

// Do 创建信号（POST /api/v5/tradingBot/signal/create-signal）。
func (s *TradingBotSignalCreateSignalService) Do(ctx context.Context) (*TradingBotSignalCreateAck, error) {
	if s.r.SignalChanName == "" {
		return nil, errTradingBotSignalCreateSignalMissingName
	}

	var data []TradingBotSignalCreateAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/tradingBot/signal/create-signal", nil, s.r, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/tradingBot/signal/create-signal", requestID, errEmptyTradingBotSignalCreateSignalResponse)
	}
	if len(data) != 1 {
		return nil, newInvalidDataAPIError(
			http.MethodPost,
			"/api/v5/tradingBot/signal/create-signal",
			requestID,
			fmt.Errorf("%w: expected 1 ack, got %d", errInvalidTradingBotSignalCreateSignalResponse, len(data)),
		)
	}
	return &data[0], nil
}
