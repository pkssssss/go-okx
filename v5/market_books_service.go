package okx

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// OrderBookLevel 表示盘口档位。
//
// OKX 返回为数组：
// - market/books: ["px","sz","0","numOrders"]（第3位字段已弃用，始终为 "0"）
// - market/books-full: ["px","sz","numOrders"]
type OrderBookLevel struct {
	Px        string
	Sz        string
	LiqOrd    string
	NumOrders string
}

func (l *OrderBookLevel) UnmarshalJSON(data []byte) error {
	*l = OrderBookLevel{}

	var arr []string
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	if len(arr) < 3 {
		return errors.New("okx: invalid order book level")
	}

	l.Px = arr[0]
	l.Sz = arr[1]
	if len(arr) >= 4 {
		l.LiqOrd = arr[2]
		l.NumOrders = arr[3]
		return nil
	}

	l.NumOrders = arr[2]
	return nil
}

// OrderBook 表示产品深度。
type OrderBook struct {
	Asks []OrderBookLevel `json:"asks"`
	Bids []OrderBookLevel `json:"bids"`
	TS   int64            `json:"ts,string"`
}

// MarketBooksService 获取产品深度。
type MarketBooksService struct {
	c      *Client
	instId string
	sz     *int
}

// NewMarketBooksService 创建 MarketBooksService。
func (c *Client) NewMarketBooksService() *MarketBooksService {
	return &MarketBooksService{c: c}
}

// InstId 设置产品 ID。
func (s *MarketBooksService) InstId(instId string) *MarketBooksService {
	s.instId = instId
	return s
}

// Sz 设置深度档位数量（最大 400）。
func (s *MarketBooksService) Sz(sz int) *MarketBooksService {
	s.sz = &sz
	return s
}

var (
	errMarketBooksMissingInstId = errors.New("okx: market books requires instId")
	errEmptyMarketBooksResponse = errors.New("okx: empty market books response")
)

// Do 获取产品深度（GET /api/v5/market/books）。
func (s *MarketBooksService) Do(ctx context.Context) (*OrderBook, error) {
	if s.instId == "" {
		return nil, errMarketBooksMissingInstId
	}

	q := url.Values{}
	q.Set("instId", s.instId)
	if s.sz != nil {
		q.Set("sz", strconv.Itoa(*s.sz))
	}

	var data []OrderBook
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/market/books", q, nil, false, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyMarketBooksResponse
	}
	return &data[0], nil
}
