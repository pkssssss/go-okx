package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// TradingBotSignalSignalsService 查询所有信号。
type TradingBotSignalSignalsService struct {
	c *Client

	signalSourceType string
	signalChanId     string
	after            string
	before           string
	limit            *int
}

// NewTradingBotSignalSignalsService 创建 TradingBotSignalSignalsService。
func (c *Client) NewTradingBotSignalSignalsService() *TradingBotSignalSignalsService {
	return &TradingBotSignalSignalsService{c: c}
}

func (s *TradingBotSignalSignalsService) SignalSourceType(signalSourceType string) *TradingBotSignalSignalsService {
	s.signalSourceType = signalSourceType
	return s
}

func (s *TradingBotSignalSignalsService) SignalChanId(signalChanId string) *TradingBotSignalSignalsService {
	s.signalChanId = signalChanId
	return s
}

func (s *TradingBotSignalSignalsService) After(after string) *TradingBotSignalSignalsService {
	s.after = after
	return s
}

func (s *TradingBotSignalSignalsService) Before(before string) *TradingBotSignalSignalsService {
	s.before = before
	return s
}

func (s *TradingBotSignalSignalsService) Limit(limit int) *TradingBotSignalSignalsService {
	s.limit = &limit
	return s
}

var errTradingBotSignalSignalsMissingSignalSourceType = errors.New("okx: tradingBot signal signals requires signalSourceType")

// Do 查询所有信号（GET /api/v5/tradingBot/signal/signals）。
func (s *TradingBotSignalSignalsService) Do(ctx context.Context) ([]TradingBotSignal, error) {
	if s.signalSourceType == "" {
		return nil, errTradingBotSignalSignalsMissingSignalSourceType
	}

	q := url.Values{}
	q.Set("signalSourceType", s.signalSourceType)
	if s.signalChanId != "" {
		q.Set("signalChanId", s.signalChanId)
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

	var data []TradingBotSignal
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/tradingBot/signal/signals", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
