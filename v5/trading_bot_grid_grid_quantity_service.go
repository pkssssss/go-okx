package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// TradingBotGridGridQuantityService 获取最大网格数量（公共）。
type TradingBotGridGridQuantityService struct {
	c *Client

	instId      string
	runType     string
	algoOrdType string
	maxPx       string
	minPx       string
	lever       string
}

// NewTradingBotGridGridQuantityService 创建 TradingBotGridGridQuantityService。
func (c *Client) NewTradingBotGridGridQuantityService() *TradingBotGridGridQuantityService {
	return &TradingBotGridGridQuantityService{c: c}
}

func (s *TradingBotGridGridQuantityService) InstId(instId string) *TradingBotGridGridQuantityService {
	s.instId = instId
	return s
}

func (s *TradingBotGridGridQuantityService) RunType(runType string) *TradingBotGridGridQuantityService {
	s.runType = runType
	return s
}

func (s *TradingBotGridGridQuantityService) AlgoOrdType(algoOrdType string) *TradingBotGridGridQuantityService {
	s.algoOrdType = algoOrdType
	return s
}

func (s *TradingBotGridGridQuantityService) MaxPx(maxPx string) *TradingBotGridGridQuantityService {
	s.maxPx = maxPx
	return s
}

func (s *TradingBotGridGridQuantityService) MinPx(minPx string) *TradingBotGridGridQuantityService {
	s.minPx = minPx
	return s
}

// Lever 设置杠杆倍数（可选；合约网格时必填）。
func (s *TradingBotGridGridQuantityService) Lever(lever string) *TradingBotGridGridQuantityService {
	s.lever = lever
	return s
}

var (
	errTradingBotGridGridQuantityMissingRequired = errors.New("okx: tradingBot grid grid-quantity requires instId, runType, algoOrdType, maxPx and minPx")
	errEmptyTradingBotGridGridQuantityResponse   = errors.New("okx: empty tradingBot grid grid-quantity response")
)

// Do 获取最大网格数量（GET /api/v5/tradingBot/grid/grid-quantity）。
func (s *TradingBotGridGridQuantityService) Do(ctx context.Context) (*TradingBotGridMaxGridQty, error) {
	if s.instId == "" || s.runType == "" || s.algoOrdType == "" || s.maxPx == "" || s.minPx == "" {
		return nil, errTradingBotGridGridQuantityMissingRequired
	}

	q := url.Values{}
	q.Set("instId", s.instId)
	q.Set("runType", s.runType)
	q.Set("algoOrdType", s.algoOrdType)
	q.Set("maxPx", s.maxPx)
	q.Set("minPx", s.minPx)
	if s.lever != "" {
		q.Set("lever", s.lever)
	}

	var data []TradingBotGridMaxGridQty
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/tradingBot/grid/grid-quantity", q, nil, false, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyTradingBotGridGridQuantityResponse
	}
	return &data[0], nil
}
