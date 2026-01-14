package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// TradingBotSignalEventHistoryService 获取信号策略历史事件。
type TradingBotSignalEventHistoryService struct {
	c *Client

	algoId string
	after  string
	before string
	limit  *int
}

// NewTradingBotSignalEventHistoryService 创建 TradingBotSignalEventHistoryService。
func (c *Client) NewTradingBotSignalEventHistoryService() *TradingBotSignalEventHistoryService {
	return &TradingBotSignalEventHistoryService{c: c}
}

func (s *TradingBotSignalEventHistoryService) AlgoId(algoId string) *TradingBotSignalEventHistoryService {
	s.algoId = algoId
	return s
}

func (s *TradingBotSignalEventHistoryService) After(after string) *TradingBotSignalEventHistoryService {
	s.after = after
	return s
}

func (s *TradingBotSignalEventHistoryService) Before(before string) *TradingBotSignalEventHistoryService {
	s.before = before
	return s
}

func (s *TradingBotSignalEventHistoryService) Limit(limit int) *TradingBotSignalEventHistoryService {
	s.limit = &limit
	return s
}

var errTradingBotSignalEventHistoryMissingAlgoId = errors.New("okx: tradingBot signal event-history requires algoId")

// Do 获取信号策略历史事件（GET /api/v5/tradingBot/signal/event-history）。
func (s *TradingBotSignalEventHistoryService) Do(ctx context.Context) ([]TradingBotSignalEventHistory, error) {
	if s.algoId == "" {
		return nil, errTradingBotSignalEventHistoryMissingAlgoId
	}

	q := url.Values{}
	q.Set("algoId", s.algoId)
	if s.after != "" {
		q.Set("after", s.after)
	}
	if s.before != "" {
		q.Set("before", s.before)
	}
	if s.limit != nil {
		q.Set("limit", strconv.Itoa(*s.limit))
	}

	var data []TradingBotSignalEventHistory
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/tradingBot/signal/event-history", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
