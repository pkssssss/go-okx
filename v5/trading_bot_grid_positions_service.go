package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// TradingBotGridPositionsService 获取网格策略委托持仓。
type TradingBotGridPositionsService struct {
	c *Client

	algoOrdType string
	algoId      string
}

// NewTradingBotGridPositionsService 创建 TradingBotGridPositionsService。
func (c *Client) NewTradingBotGridPositionsService() *TradingBotGridPositionsService {
	return &TradingBotGridPositionsService{c: c}
}

func (s *TradingBotGridPositionsService) AlgoOrdType(algoOrdType string) *TradingBotGridPositionsService {
	s.algoOrdType = algoOrdType
	return s
}

func (s *TradingBotGridPositionsService) AlgoId(algoId string) *TradingBotGridPositionsService {
	s.algoId = algoId
	return s
}

var errTradingBotGridPositionsMissingRequired = errors.New("okx: tradingBot grid positions requires algoOrdType and algoId")

// Do 获取网格策略委托持仓（GET /api/v5/tradingBot/grid/positions）。
func (s *TradingBotGridPositionsService) Do(ctx context.Context) ([]TradingBotGridPosition, error) {
	if s.algoOrdType == "" || s.algoId == "" {
		return nil, errTradingBotGridPositionsMissingRequired
	}

	q := url.Values{}
	q.Set("algoOrdType", s.algoOrdType)
	q.Set("algoId", s.algoId)

	var data []TradingBotGridPosition
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/tradingBot/grid/positions", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
