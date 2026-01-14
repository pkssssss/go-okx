package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// TradingBotPublicRSIBackTestingService 获取 RSI 回测（公共）。
type TradingBotPublicRSIBackTestingService struct {
	c *Client

	instId      string
	timeframe   string
	thold       string
	timePeriod  string
	triggerCond string
	duration    string
}

// NewTradingBotPublicRSIBackTestingService 创建 TradingBotPublicRSIBackTestingService。
func (c *Client) NewTradingBotPublicRSIBackTestingService() *TradingBotPublicRSIBackTestingService {
	return &TradingBotPublicRSIBackTestingService{c: c}
}

func (s *TradingBotPublicRSIBackTestingService) InstId(instId string) *TradingBotPublicRSIBackTestingService {
	s.instId = instId
	return s
}

func (s *TradingBotPublicRSIBackTestingService) Timeframe(timeframe string) *TradingBotPublicRSIBackTestingService {
	s.timeframe = timeframe
	return s
}

func (s *TradingBotPublicRSIBackTestingService) Thold(thold string) *TradingBotPublicRSIBackTestingService {
	s.thold = thold
	return s
}

func (s *TradingBotPublicRSIBackTestingService) TimePeriod(timePeriod string) *TradingBotPublicRSIBackTestingService {
	s.timePeriod = timePeriod
	return s
}

func (s *TradingBotPublicRSIBackTestingService) TriggerCond(triggerCond string) *TradingBotPublicRSIBackTestingService {
	s.triggerCond = triggerCond
	return s
}

func (s *TradingBotPublicRSIBackTestingService) Duration(duration string) *TradingBotPublicRSIBackTestingService {
	s.duration = duration
	return s
}

var (
	errTradingBotPublicRSIBackTestingMissingRequired = errors.New("okx: tradingBot public rsi-back-testing requires instId, timeframe, thold and timePeriod")
	errEmptyTradingBotPublicRSIBackTestingResponse   = errors.New("okx: empty tradingBot public rsi-back-testing response")
)

// Do 获取 RSI 回测（GET /api/v5/tradingBot/public/rsi-back-testing）。
func (s *TradingBotPublicRSIBackTestingService) Do(ctx context.Context) (*TradingBotPublicRSIBackTestingResult, error) {
	if s.instId == "" || s.timeframe == "" || s.thold == "" || s.timePeriod == "" {
		return nil, errTradingBotPublicRSIBackTestingMissingRequired
	}

	q := url.Values{}
	q.Set("instId", s.instId)
	q.Set("timeframe", s.timeframe)
	q.Set("thold", s.thold)
	q.Set("timePeriod", s.timePeriod)
	if s.triggerCond != "" {
		q.Set("triggerCond", s.triggerCond)
	}
	if s.duration != "" {
		q.Set("duration", s.duration)
	}

	var data []TradingBotPublicRSIBackTestingResult
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/tradingBot/public/rsi-back-testing", q, nil, false, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyTradingBotPublicRSIBackTestingResponse
	}
	return &data[0], nil
}
