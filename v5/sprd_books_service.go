package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// SprdBooksService 获取 Spread 产品深度（公共）。
type SprdBooksService struct {
	c      *Client
	sprdId string
	sz     *int
}

// NewSprdBooksService 创建 SprdBooksService。
func (c *Client) NewSprdBooksService() *SprdBooksService {
	return &SprdBooksService{c: c}
}

// SprdId 设置 Spread ID。
func (s *SprdBooksService) SprdId(sprdId string) *SprdBooksService {
	s.sprdId = sprdId
	return s
}

// Sz 设置深度档位数量（最大 400；默认 5）。
func (s *SprdBooksService) Sz(sz int) *SprdBooksService {
	s.sz = &sz
	return s
}

var (
	errSprdBooksMissingSprdId = errors.New("okx: sprd books requires sprdId")
	errEmptySprdBooksResponse = errors.New("okx: empty sprd books response")
)

// Do 获取 Spread 产品深度（GET /api/v5/sprd/books）。
func (s *SprdBooksService) Do(ctx context.Context) (*OrderBook, error) {
	if s.sprdId == "" {
		return nil, errSprdBooksMissingSprdId
	}

	q := url.Values{}
	q.Set("sprdId", s.sprdId)
	if s.sz != nil {
		q.Set("sz", strconv.Itoa(*s.sz))
	}

	var data []OrderBook
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/sprd/books", q, nil, false, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptySprdBooksResponse
	}
	return &data[0], nil
}
