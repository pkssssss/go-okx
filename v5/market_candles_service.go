package okx

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// Candle 表示 K 线数据。
//
// OKX 返回为数组：["ts","o","h","l","c","vol","volCcy","volCcyQuote","confirm"]
// 不同市场/接口可能字段个数不同；这里采用“最小必需字段 + 兼容可选字段”的解析策略。
type Candle struct {
	TS int64

	Open  string
	High  string
	Low   string
	Close string

	Vol         string
	VolCcy      string
	VolCcyQuote string
	Confirm     string
}

func (c *Candle) UnmarshalJSON(data []byte) error {
	*c = Candle{}

	var arr []string
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	if len(arr) < 6 {
		return errors.New("okx: invalid candle")
	}

	ts, err := strconv.ParseInt(arr[0], 10, 64)
	if err != nil {
		return errors.New("okx: invalid candle ts")
	}

	c.TS = ts
	c.Open = arr[1]
	c.High = arr[2]
	c.Low = arr[3]
	c.Close = arr[4]
	c.Vol = arr[5]
	if len(arr) > 6 {
		c.VolCcy = arr[6]
	}
	if len(arr) > 7 {
		c.VolCcyQuote = arr[7]
	}
	if len(arr) > 8 {
		c.Confirm = arr[8]
	}
	return nil
}

// MarketCandlesService 获取 K 线数据。
type MarketCandlesService struct {
	c *Client

	instId string
	bar    string
	after  string
	before string
	limit  *int
}

// NewMarketCandlesService 创建 MarketCandlesService。
func (c *Client) NewMarketCandlesService() *MarketCandlesService {
	return &MarketCandlesService{c: c}
}

// InstId 设置产品 ID（必填）。
func (s *MarketCandlesService) InstId(instId string) *MarketCandlesService {
	s.instId = instId
	return s
}

// Bar 设置 K 线粒度，如 1m/5m/1H/1D 等。
func (s *MarketCandlesService) Bar(bar string) *MarketCandlesService {
	s.bar = bar
	return s
}

// After 设置请求此时间戳之后的数据（毫秒字符串）。
func (s *MarketCandlesService) After(after string) *MarketCandlesService {
	s.after = after
	return s
}

// Before 设置请求此时间戳之前的数据（毫秒字符串）。
func (s *MarketCandlesService) Before(before string) *MarketCandlesService {
	s.before = before
	return s
}

// Limit 设置返回条数。
func (s *MarketCandlesService) Limit(limit int) *MarketCandlesService {
	s.limit = &limit
	return s
}

var errMarketCandlesMissingInstId = errors.New("okx: market candles requires instId")

// Do 获取 K 线数据（GET /api/v5/market/candles）。
func (s *MarketCandlesService) Do(ctx context.Context) ([]Candle, error) {
	if s.instId == "" {
		return nil, errMarketCandlesMissingInstId
	}

	q := url.Values{}
	q.Set("instId", s.instId)
	if s.bar != "" {
		q.Set("bar", s.bar)
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

	var data []Candle
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/market/candles", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
