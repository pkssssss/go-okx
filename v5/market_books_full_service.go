package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// MarketBooksFullService 获取产品完整深度。
type MarketBooksFullService struct {
	c      *Client
	instId string
	sz     *int
}

// NewMarketBooksFullService 创建 MarketBooksFullService。
func (c *Client) NewMarketBooksFullService() *MarketBooksFullService {
	return &MarketBooksFullService{c: c}
}

// InstId 设置产品 ID（必填）。
func (s *MarketBooksFullService) InstId(instId string) *MarketBooksFullService {
	s.instId = instId
	return s
}

// Sz 设置深度档位数量（最大 5000，即买卖深度共 10000 条）。
func (s *MarketBooksFullService) Sz(sz int) *MarketBooksFullService {
	s.sz = &sz
	return s
}

var (
	errMarketBooksFullMissingInstId = errors.New("okx: market books full requires instId")
	errEmptyMarketBooksFullResponse = errors.New("okx: empty market books full response")
)

// Do 获取产品完整深度（GET /api/v5/market/books-full）。
func (s *MarketBooksFullService) Do(ctx context.Context) (*OrderBook, error) {
	if s.instId == "" {
		return nil, errMarketBooksFullMissingInstId
	}

	q := url.Values{}
	q.Set("instId", s.instId)
	if s.sz != nil {
		q.Set("sz", strconv.Itoa(*s.sz))
	}

	var data []OrderBook
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/market/books-full", q, nil, false, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyMarketBooksFullResponse
	}
	return &data[0], nil
}
