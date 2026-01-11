package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// SprdOrdersPendingService 获取未成交订单列表。
type SprdOrdersPendingService struct {
	c *Client

	sprdId  string
	ordType string
	state   string
	beginId string
	endId   string
	limit   *int
}

// NewSprdOrdersPendingService 创建 SprdOrdersPendingService。
func (c *Client) NewSprdOrdersPendingService() *SprdOrdersPendingService {
	return &SprdOrdersPendingService{c: c}
}

// SprdId 设置 spread ID（如 BTC-USDT_BTC-USDT-SWAP）。
func (s *SprdOrdersPendingService) SprdId(sprdId string) *SprdOrdersPendingService {
	s.sprdId = sprdId
	return s
}

// OrdType 设置订单类型（market/limit/post_only/ioc）。
func (s *SprdOrdersPendingService) OrdType(ordType string) *SprdOrdersPendingService {
	s.ordType = ordType
	return s
}

// State 设置订单状态（live/partially_filled）。
func (s *SprdOrdersPendingService) State(state string) *SprdOrdersPendingService {
	s.state = state
	return s
}

// BeginId 设置请求起始订单 ID（不包含 beginId）。
func (s *SprdOrdersPendingService) BeginId(beginId string) *SprdOrdersPendingService {
	s.beginId = beginId
	return s
}

// EndId 设置请求结束订单 ID（不包含 endId）。
func (s *SprdOrdersPendingService) EndId(endId string) *SprdOrdersPendingService {
	s.endId = endId
	return s
}

// Limit 设置返回条数（最大 100，默认 100）。
func (s *SprdOrdersPendingService) Limit(limit int) *SprdOrdersPendingService {
	s.limit = &limit
	return s
}

// Do 获取未成交订单列表（GET /api/v5/sprd/orders-pending）。
func (s *SprdOrdersPendingService) Do(ctx context.Context) ([]SprdOrder, error) {
	q := url.Values{}
	if s.sprdId != "" {
		q.Set("sprdId", s.sprdId)
	}
	if s.ordType != "" {
		q.Set("ordType", s.ordType)
	}
	if s.state != "" {
		q.Set("state", s.state)
	}
	if s.beginId != "" {
		q.Set("beginId", s.beginId)
	}
	if s.endId != "" {
		q.Set("endId", s.endId)
	}
	if s.limit != nil {
		q.Set("limit", strconv.Itoa(*s.limit))
	}

	var data []SprdOrder
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/sprd/orders-pending", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
