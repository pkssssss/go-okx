package okx

import (
	"context"
	"errors"
	"net/http"
)

type tradingBotSignalCancelSubOrderRequest struct {
	AlgoId      string `json:"algoId"`
	InstId      string `json:"instId"`
	SignalOrdId string `json:"signalOrdId"`
}

// TradingBotSignalCancelSubOrderService 撤销信号策略子订单。
type TradingBotSignalCancelSubOrderService struct {
	c *Client
	r tradingBotSignalCancelSubOrderRequest
}

// NewTradingBotSignalCancelSubOrderService 创建 TradingBotSignalCancelSubOrderService。
func (c *Client) NewTradingBotSignalCancelSubOrderService() *TradingBotSignalCancelSubOrderService {
	return &TradingBotSignalCancelSubOrderService{c: c}
}

func (s *TradingBotSignalCancelSubOrderService) AlgoId(algoId string) *TradingBotSignalCancelSubOrderService {
	s.r.AlgoId = algoId
	return s
}

func (s *TradingBotSignalCancelSubOrderService) InstId(instId string) *TradingBotSignalCancelSubOrderService {
	s.r.InstId = instId
	return s
}

func (s *TradingBotSignalCancelSubOrderService) SignalOrdId(signalOrdId string) *TradingBotSignalCancelSubOrderService {
	s.r.SignalOrdId = signalOrdId
	return s
}

var (
	errTradingBotSignalCancelSubOrderMissingRequired = errors.New("okx: tradingBot signal cancel-sub-order requires algoId, instId and signalOrdId")
	errEmptyTradingBotSignalCancelSubOrderResponse   = errors.New("okx: empty tradingBot signal cancel-sub-order response")
)

// Do 撤单（POST /api/v5/tradingBot/signal/cancel-sub-order）。
func (s *TradingBotSignalCancelSubOrderService) Do(ctx context.Context) (*TradingBotSignalCancelSubOrderAck, error) {
	if s.r.AlgoId == "" || s.r.InstId == "" || s.r.SignalOrdId == "" {
		return nil, errTradingBotSignalCancelSubOrderMissingRequired
	}

	var data []TradingBotSignalCancelSubOrderAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/tradingBot/signal/cancel-sub-order", nil, s.r, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyTradingBotSignalCancelSubOrderResponse
	}
	if data[0].SCode != "0" {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodPost,
			RequestPath: "/api/v5/tradingBot/signal/cancel-sub-order",
			RequestID:   requestID,
			Code:        data[0].SCode,
			Message:     data[0].SMsg,
		}
	}
	return &data[0], nil
}
