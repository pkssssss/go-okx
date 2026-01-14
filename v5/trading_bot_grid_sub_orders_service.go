package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// TradingBotGridSubOrdersService 获取网格策略委托子订单信息。
type TradingBotGridSubOrdersService struct {
	c *Client

	algoId      string
	algoOrdType string
	typ         string
	groupId     string
	after       string
	before      string
	limit       *int
}

// NewTradingBotGridSubOrdersService 创建 TradingBotGridSubOrdersService。
func (c *Client) NewTradingBotGridSubOrdersService() *TradingBotGridSubOrdersService {
	return &TradingBotGridSubOrdersService{c: c}
}

func (s *TradingBotGridSubOrdersService) AlgoId(algoId string) *TradingBotGridSubOrdersService {
	s.algoId = algoId
	return s
}

func (s *TradingBotGridSubOrdersService) AlgoOrdType(algoOrdType string) *TradingBotGridSubOrdersService {
	s.algoOrdType = algoOrdType
	return s
}

// Type 设置子订单状态（必填：live/filled）。
func (s *TradingBotGridSubOrdersService) Type(typ string) *TradingBotGridSubOrdersService {
	s.typ = typ
	return s
}

func (s *TradingBotGridSubOrdersService) GroupId(groupId string) *TradingBotGridSubOrdersService {
	s.groupId = groupId
	return s
}

func (s *TradingBotGridSubOrdersService) After(after string) *TradingBotGridSubOrdersService {
	s.after = after
	return s
}

func (s *TradingBotGridSubOrdersService) Before(before string) *TradingBotGridSubOrdersService {
	s.before = before
	return s
}

func (s *TradingBotGridSubOrdersService) Limit(limit int) *TradingBotGridSubOrdersService {
	s.limit = &limit
	return s
}

var errTradingBotGridSubOrdersMissingRequired = errors.New("okx: tradingBot grid sub-orders requires algoId, algoOrdType and type")

// Do 获取网格策略委托子订单信息（GET /api/v5/tradingBot/grid/sub-orders）。
func (s *TradingBotGridSubOrdersService) Do(ctx context.Context) ([]TradingBotGridSubOrder, error) {
	if s.algoId == "" || s.algoOrdType == "" || s.typ == "" {
		return nil, errTradingBotGridSubOrdersMissingRequired
	}

	q := url.Values{}
	q.Set("algoId", s.algoId)
	q.Set("algoOrdType", s.algoOrdType)
	q.Set("type", s.typ)
	if s.groupId != "" {
		q.Set("groupId", s.groupId)
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

	var data []TradingBotGridSubOrder
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/tradingBot/grid/sub-orders", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
