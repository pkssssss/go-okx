package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// AlgoOrdersHistoryService 获取历史策略委托单列表（最近 3 个月）。
type AlgoOrdersHistoryService struct {
	c *Client

	ordType  string
	state    string
	algoId   string
	instType string
	instId   string
	after    string
	before   string
	limit    *int
}

// NewAlgoOrdersHistoryService 创建 AlgoOrdersHistoryService。
func (c *Client) NewAlgoOrdersHistoryService() *AlgoOrdersHistoryService {
	return &AlgoOrdersHistoryService{c: c}
}

// OrdType 设置订单类型（必填）。
func (s *AlgoOrdersHistoryService) OrdType(ordType string) *AlgoOrdersHistoryService {
	s.ordType = ordType
	return s
}

// State 设置订单状态过滤（effective/canceled/order_failed；与 AlgoId 二选一）。
func (s *AlgoOrdersHistoryService) State(state string) *AlgoOrdersHistoryService {
	s.state = state
	return s
}

// AlgoId 设置策略委托单ID过滤（与 State 二选一）。
func (s *AlgoOrdersHistoryService) AlgoId(algoId string) *AlgoOrdersHistoryService {
	s.algoId = algoId
	return s
}

func (s *AlgoOrdersHistoryService) InstType(instType string) *AlgoOrdersHistoryService {
	s.instType = instType
	return s
}

func (s *AlgoOrdersHistoryService) InstId(instId string) *AlgoOrdersHistoryService {
	s.instId = instId
	return s
}

func (s *AlgoOrdersHistoryService) After(after string) *AlgoOrdersHistoryService {
	s.after = after
	return s
}

func (s *AlgoOrdersHistoryService) Before(before string) *AlgoOrdersHistoryService {
	s.before = before
	return s
}

func (s *AlgoOrdersHistoryService) Limit(limit int) *AlgoOrdersHistoryService {
	s.limit = &limit
	return s
}

var (
	errAlgoOrdersHistoryMissingOrdType         = errors.New("okx: orders algo history requires ordType")
	errAlgoOrdersHistoryMissingStateOrAlgoId   = errors.New("okx: orders algo history requires state or algoId")
	errAlgoOrdersHistoryStateAndAlgoIdConflict = errors.New("okx: orders algo history requires at most one of state/algoId")
)

// Do 获取历史策略委托单列表（GET /api/v5/trade/orders-algo-history）。
func (s *AlgoOrdersHistoryService) Do(ctx context.Context) ([]TradeAlgoOrder, error) {
	if s.ordType == "" {
		return nil, errAlgoOrdersHistoryMissingOrdType
	}
	if s.state == "" && s.algoId == "" {
		return nil, errAlgoOrdersHistoryMissingStateOrAlgoId
	}
	if s.state != "" && s.algoId != "" {
		return nil, errAlgoOrdersHistoryStateAndAlgoIdConflict
	}

	q := url.Values{}
	q.Set("ordType", s.ordType)
	if s.state != "" {
		q.Set("state", s.state)
	}
	if s.algoId != "" {
		q.Set("algoId", s.algoId)
	}
	if s.instType != "" {
		q.Set("instType", s.instType)
	}
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

	var data []TradeAlgoOrder
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/trade/orders-algo-history", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
