package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// TradingBotSignalPositionsHistoryService 查看历史持仓信息（最近 3 个月有更新的仓位）。
type TradingBotSignalPositionsHistoryService struct {
	c *Client

	algoId string
	instId string
	after  string
	before string
	limit  *int
}

// NewTradingBotSignalPositionsHistoryService 创建 TradingBotSignalPositionsHistoryService。
func (c *Client) NewTradingBotSignalPositionsHistoryService() *TradingBotSignalPositionsHistoryService {
	return &TradingBotSignalPositionsHistoryService{c: c}
}

func (s *TradingBotSignalPositionsHistoryService) AlgoId(algoId string) *TradingBotSignalPositionsHistoryService {
	s.algoId = algoId
	return s
}

func (s *TradingBotSignalPositionsHistoryService) InstId(instId string) *TradingBotSignalPositionsHistoryService {
	s.instId = instId
	return s
}

func (s *TradingBotSignalPositionsHistoryService) After(after string) *TradingBotSignalPositionsHistoryService {
	s.after = after
	return s
}

func (s *TradingBotSignalPositionsHistoryService) Before(before string) *TradingBotSignalPositionsHistoryService {
	s.before = before
	return s
}

func (s *TradingBotSignalPositionsHistoryService) Limit(limit int) *TradingBotSignalPositionsHistoryService {
	s.limit = &limit
	return s
}

var errTradingBotSignalPositionsHistoryMissingAlgoId = errors.New("okx: tradingBot signal positions-history requires algoId")

// Do 查看历史持仓信息（GET /api/v5/tradingBot/signal/positions-history）。
func (s *TradingBotSignalPositionsHistoryService) Do(ctx context.Context) ([]TradingBotSignalPositionsHistory, error) {
	if s.algoId == "" {
		return nil, errTradingBotSignalPositionsHistoryMissingAlgoId
	}

	q := url.Values{}
	q.Set("algoId", s.algoId)
	if s.instId != "" {
		q.Set("instId", s.instId)
	}
	if s.after != "" {
		q.Set("after", s.after)
	}
	if s.before != "" {
		q.Set("before", s.before)
	}
	if s.limit != nil {
		q.Set("limit", strconv.Itoa(*s.limit))
	}

	var data []TradingBotSignalPositionsHistory
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/tradingBot/signal/positions-history", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
