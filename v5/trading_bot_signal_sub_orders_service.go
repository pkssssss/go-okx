package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// TradingBotSignalSubOrdersService 获取信号策略子订单信息。
type TradingBotSignalSubOrdersService struct {
	c *Client

	algoId      string
	algoOrdType string
	state       string
	signalOrdId string

	after  string
	before string
	begin  string
	end    string
	limit  *int

	typ     string
	clOrdId string
}

// NewTradingBotSignalSubOrdersService 创建 TradingBotSignalSubOrdersService。
func (c *Client) NewTradingBotSignalSubOrdersService() *TradingBotSignalSubOrdersService {
	return &TradingBotSignalSubOrdersService{c: c}
}

func (s *TradingBotSignalSubOrdersService) AlgoId(algoId string) *TradingBotSignalSubOrdersService {
	s.algoId = algoId
	return s
}

func (s *TradingBotSignalSubOrdersService) AlgoOrdType(algoOrdType string) *TradingBotSignalSubOrdersService {
	s.algoOrdType = algoOrdType
	return s
}

// State 设置子订单状态（state 与 signalOrdId 必须传一个，若传两个，以 state 为主）。
func (s *TradingBotSignalSubOrdersService) State(state string) *TradingBotSignalSubOrdersService {
	s.state = state
	return s
}

func (s *TradingBotSignalSubOrdersService) SignalOrdId(signalOrdId string) *TradingBotSignalSubOrdersService {
	s.signalOrdId = signalOrdId
	return s
}

func (s *TradingBotSignalSubOrdersService) After(after string) *TradingBotSignalSubOrdersService {
	s.after = after
	return s
}

func (s *TradingBotSignalSubOrdersService) Before(before string) *TradingBotSignalSubOrdersService {
	s.before = before
	return s
}

func (s *TradingBotSignalSubOrdersService) Begin(begin string) *TradingBotSignalSubOrdersService {
	s.begin = begin
	return s
}

func (s *TradingBotSignalSubOrdersService) End(end string) *TradingBotSignalSubOrdersService {
	s.end = end
	return s
}

func (s *TradingBotSignalSubOrdersService) Limit(limit int) *TradingBotSignalSubOrdersService {
	s.limit = &limit
	return s
}

// Type 设置子订单类型（文档标记为即将废弃）。
func (s *TradingBotSignalSubOrdersService) Type(typ string) *TradingBotSignalSubOrdersService {
	s.typ = typ
	return s
}

// ClOrdId 设置子订单自定义订单ID（文档标记为即将废弃）。
func (s *TradingBotSignalSubOrdersService) ClOrdId(clOrdId string) *TradingBotSignalSubOrdersService {
	s.clOrdId = clOrdId
	return s
}

var (
	errTradingBotSignalSubOrdersMissingRequired           = errors.New("okx: tradingBot signal sub-orders requires algoId and algoOrdType")
	errTradingBotSignalSubOrdersMissingStateOrSignalOrdId = errors.New("okx: tradingBot signal sub-orders requires state or signalOrdId")
)

// Do 获取信号策略子订单信息（GET /api/v5/tradingBot/signal/sub-orders）。
func (s *TradingBotSignalSubOrdersService) Do(ctx context.Context) ([]TradingBotSignalSubOrder, error) {
	if s.algoId == "" || s.algoOrdType == "" {
		return nil, errTradingBotSignalSubOrdersMissingRequired
	}
	if s.state == "" && s.signalOrdId == "" {
		return nil, errTradingBotSignalSubOrdersMissingStateOrSignalOrdId
	}

	q := url.Values{}
	q.Set("algoId", s.algoId)
	q.Set("algoOrdType", s.algoOrdType)
	if s.state != "" {
		q.Set("state", s.state)
	} else {
		q.Set("signalOrdId", s.signalOrdId)
	}

	if s.after != "" {
		q.Set("after", s.after)
	}
	if s.before != "" {
		q.Set("before", s.before)
	}
	if s.begin != "" {
		q.Set("begin", s.begin)
	}
	if s.end != "" {
		q.Set("end", s.end)
	}
	if s.limit != nil {
		q.Set("limit", strconv.Itoa(*s.limit))
	}
	if s.typ != "" {
		q.Set("type", s.typ)
	}
	if s.clOrdId != "" {
		q.Set("clOrdId", s.clOrdId)
	}

	var data []TradingBotSignalSubOrder
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/tradingBot/signal/sub-orders", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
