package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// PublicPriceLimitService 查询单个交易产品的最高买价和最低卖价。
type PublicPriceLimitService struct {
	c *Client

	instId string
}

// NewPublicPriceLimitService 创建 PublicPriceLimitService。
func (c *Client) NewPublicPriceLimitService() *PublicPriceLimitService {
	return &PublicPriceLimitService{c: c}
}

// InstId 设置产品ID，如 BTC-USDT-SWAP（必填）。
func (s *PublicPriceLimitService) InstId(instId string) *PublicPriceLimitService {
	s.instId = instId
	return s
}

var errPublicPriceLimitMissingInstId = errors.New("okx: public price limit requires instId")

// Do 查询限价（GET /api/v5/public/price-limit）。
func (s *PublicPriceLimitService) Do(ctx context.Context) ([]PriceLimit, error) {
	if s.instId == "" {
		return nil, errPublicPriceLimitMissingInstId
	}

	q := url.Values{}
	q.Set("instId", s.instId)

	var data []PriceLimit
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/public/price-limit", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
