package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// TradingBotSignalPositionsService 获取信号策略持仓。
type TradingBotSignalPositionsService struct {
	c *Client

	algoOrdType string
	algoId      string
}

// NewTradingBotSignalPositionsService 创建 TradingBotSignalPositionsService。
func (c *Client) NewTradingBotSignalPositionsService() *TradingBotSignalPositionsService {
	return &TradingBotSignalPositionsService{c: c}
}

func (s *TradingBotSignalPositionsService) AlgoOrdType(algoOrdType string) *TradingBotSignalPositionsService {
	s.algoOrdType = algoOrdType
	return s
}

func (s *TradingBotSignalPositionsService) AlgoId(algoId string) *TradingBotSignalPositionsService {
	s.algoId = algoId
	return s
}

var errTradingBotSignalPositionsMissingRequired = errors.New("okx: tradingBot signal positions requires algoOrdType and algoId")

// Do 获取信号策略持仓（GET /api/v5/tradingBot/signal/positions）。
func (s *TradingBotSignalPositionsService) Do(ctx context.Context) ([]TradingBotSignalPosition, error) {
	if s.algoOrdType == "" || s.algoId == "" {
		return nil, errTradingBotSignalPositionsMissingRequired
	}

	q := url.Values{}
	q.Set("algoOrdType", s.algoOrdType)
	q.Set("algoId", s.algoId)

	var data []TradingBotSignalPosition
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/tradingBot/signal/positions", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
