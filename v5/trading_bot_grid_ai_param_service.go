package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// TradingBotGridAIParamService 获取网格策略智能回测参数（公共）。
type TradingBotGridAIParamService struct {
	c *Client

	algoOrdType string
	instId      string
	direction   string
	duration    string
}

// NewTradingBotGridAIParamService 创建 TradingBotGridAIParamService。
func (c *Client) NewTradingBotGridAIParamService() *TradingBotGridAIParamService {
	return &TradingBotGridAIParamService{c: c}
}

// AlgoOrdType 设置策略订单类型（必填）。
func (s *TradingBotGridAIParamService) AlgoOrdType(algoOrdType string) *TradingBotGridAIParamService {
	s.algoOrdType = algoOrdType
	return s
}

// InstId 设置产品 ID（必填）。
func (s *TradingBotGridAIParamService) InstId(instId string) *TradingBotGridAIParamService {
	s.instId = instId
	return s
}

// Direction 设置合约网格类型（可选；合约网格必填）。
func (s *TradingBotGridAIParamService) Direction(direction string) *TradingBotGridAIParamService {
	s.direction = direction
	return s
}

// Duration 设置回测周期（可选）。
func (s *TradingBotGridAIParamService) Duration(duration string) *TradingBotGridAIParamService {
	s.duration = duration
	return s
}

var errTradingBotGridAIParamMissingRequired = errors.New("okx: tradingBot grid ai-param requires algoOrdType and instId")

// Do 获取网格策略智能回测参数（GET /api/v5/tradingBot/grid/ai-param）。
func (s *TradingBotGridAIParamService) Do(ctx context.Context) ([]TradingBotGridAIParam, error) {
	if s.algoOrdType == "" || s.instId == "" {
		return nil, errTradingBotGridAIParamMissingRequired
	}

	q := url.Values{}
	q.Set("algoOrdType", s.algoOrdType)
	q.Set("instId", s.instId)
	if s.direction != "" {
		q.Set("direction", s.direction)
	}
	if s.duration != "" {
		q.Set("duration", s.duration)
	}

	var data []TradingBotGridAIParam
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/tradingBot/grid/ai-param", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
