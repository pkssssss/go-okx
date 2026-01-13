package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// FiatBuySellHistoryService 获取买卖交易历史。
type FiatBuySellHistoryService struct {
	c *Client

	ordId   string
	clOrdId string
	state   string
	begin   string
	end     string
	limit   *int
}

// NewFiatBuySellHistoryService 创建 FiatBuySellHistoryService。
func (c *Client) NewFiatBuySellHistoryService() *FiatBuySellHistoryService {
	return &FiatBuySellHistoryService{c: c}
}

// OrdId 设置订单ID（可选）。
func (s *FiatBuySellHistoryService) OrdId(ordId string) *FiatBuySellHistoryService {
	s.ordId = ordId
	return s
}

// ClOrdId 设置用户自定义的订单标识（可选）。
func (s *FiatBuySellHistoryService) ClOrdId(clOrdId string) *FiatBuySellHistoryService {
	s.clOrdId = clOrdId
	return s
}

// State 设置交易状态（可选）：processing/completed/failed。
func (s *FiatBuySellHistoryService) State(state string) *FiatBuySellHistoryService {
	s.state = state
	return s
}

// Begin 设置开始时间（可选，Unix ms）。
func (s *FiatBuySellHistoryService) Begin(begin string) *FiatBuySellHistoryService {
	s.begin = begin
	return s
}

// End 设置结束时间（可选，Unix ms）。
func (s *FiatBuySellHistoryService) End(end string) *FiatBuySellHistoryService {
	s.end = end
	return s
}

// Limit 设置返回条数（可选，最大 100）。
func (s *FiatBuySellHistoryService) Limit(limit int) *FiatBuySellHistoryService {
	s.limit = &limit
	return s
}

// Do 获取买卖交易历史（GET /api/v5/fiat/buy-sell/history）。
func (s *FiatBuySellHistoryService) Do(ctx context.Context) ([]FiatBuySellOrder, error) {
	q := url.Values{}
	if s.ordId != "" {
		q.Set("ordId", s.ordId)
	}
	if s.clOrdId != "" {
		q.Set("clOrdId", s.clOrdId)
	}
	if s.state != "" {
		q.Set("state", s.state)
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

	var data []FiatBuySellOrder
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/fiat/buy-sell/history", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
