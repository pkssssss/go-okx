package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// AlgoOrdersPendingService 获取未完成策略委托单列表。
type AlgoOrdersPendingService struct {
	c *Client

	algoId   string
	instType string
	instId   string
	ordType  string
	after    string
	before   string
	limit    *int
}

// NewAlgoOrdersPendingService 创建 AlgoOrdersPendingService。
func (c *Client) NewAlgoOrdersPendingService() *AlgoOrdersPendingService {
	return &AlgoOrdersPendingService{c: c}
}

func (s *AlgoOrdersPendingService) AlgoId(algoId string) *AlgoOrdersPendingService {
	s.algoId = algoId
	return s
}

func (s *AlgoOrdersPendingService) InstType(instType string) *AlgoOrdersPendingService {
	s.instType = instType
	return s
}

func (s *AlgoOrdersPendingService) InstId(instId string) *AlgoOrdersPendingService {
	s.instId = instId
	return s
}

// OrdType 设置订单类型（必填；支持 conditional/oco/trigger/move_order_stop/twap 等；conditional,oco 可逗号分隔）。
func (s *AlgoOrdersPendingService) OrdType(ordType string) *AlgoOrdersPendingService {
	s.ordType = ordType
	return s
}

func (s *AlgoOrdersPendingService) After(after string) *AlgoOrdersPendingService {
	s.after = after
	return s
}

func (s *AlgoOrdersPendingService) Before(before string) *AlgoOrdersPendingService {
	s.before = before
	return s
}

func (s *AlgoOrdersPendingService) Limit(limit int) *AlgoOrdersPendingService {
	s.limit = &limit
	return s
}

var errAlgoOrdersPendingMissingOrdType = errors.New("okx: orders algo pending requires ordType")

// Do 获取未完成策略委托单列表（GET /api/v5/trade/orders-algo-pending）。
func (s *AlgoOrdersPendingService) Do(ctx context.Context) ([]TradeAlgoOrder, error) {
	if s.ordType == "" {
		return nil, errAlgoOrdersPendingMissingOrdType
	}

	q := url.Values{}
	q.Set("ordType", s.ordType)
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
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/trade/orders-algo-pending", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
