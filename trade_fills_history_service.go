package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// TradeFillsHistoryService 获取成交明细（近三个月）。
type TradeFillsHistoryService struct {
	c *Client

	instType   string
	instFamily string
	instId     string
	ordId      string
	subType    string
	after      string
	before     string
	begin      string
	end        string
	limit      *int
}

// NewTradeFillsHistoryService 创建 TradeFillsHistoryService。
func (c *Client) NewTradeFillsHistoryService() *TradeFillsHistoryService {
	return &TradeFillsHistoryService{c: c}
}

// InstType 设置产品类型（必填）。
func (s *TradeFillsHistoryService) InstType(instType string) *TradeFillsHistoryService {
	s.instType = instType
	return s
}

func (s *TradeFillsHistoryService) InstFamily(instFamily string) *TradeFillsHistoryService {
	s.instFamily = instFamily
	return s
}

func (s *TradeFillsHistoryService) InstId(instId string) *TradeFillsHistoryService {
	s.instId = instId
	return s
}

func (s *TradeFillsHistoryService) OrdId(ordId string) *TradeFillsHistoryService {
	s.ordId = ordId
	return s
}

// SubType 设置成交类型过滤（如 1: 买入, 2: 卖出...）。
func (s *TradeFillsHistoryService) SubType(subType string) *TradeFillsHistoryService {
	s.subType = subType
	return s
}

func (s *TradeFillsHistoryService) After(after string) *TradeFillsHistoryService {
	s.after = after
	return s
}

func (s *TradeFillsHistoryService) Before(before string) *TradeFillsHistoryService {
	s.before = before
	return s
}

// Begin 设置开始时间（Unix 毫秒时间戳字符串）。
func (s *TradeFillsHistoryService) Begin(begin string) *TradeFillsHistoryService {
	s.begin = begin
	return s
}

// End 设置结束时间（Unix 毫秒时间戳字符串）。
func (s *TradeFillsHistoryService) End(end string) *TradeFillsHistoryService {
	s.end = end
	return s
}

func (s *TradeFillsHistoryService) Limit(limit int) *TradeFillsHistoryService {
	s.limit = &limit
	return s
}

var errTradeFillsHistoryMissingInstType = errors.New("okx: fills history requires instType")

// Do 获取成交明细（近三个月）（GET /api/v5/trade/fills-history）。
func (s *TradeFillsHistoryService) Do(ctx context.Context) ([]TradeFill, error) {
	if s.instType == "" {
		return nil, errTradeFillsHistoryMissingInstType
	}

	q := url.Values{}
	q.Set("instType", s.instType)
	if s.instFamily != "" {
		q.Set("instFamily", s.instFamily)
	}
	if s.instId != "" {
		q.Set("instId", s.instId)
	}
	if s.ordId != "" {
		q.Set("ordId", s.ordId)
	}
	if s.subType != "" {
		q.Set("subType", s.subType)
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

	var data []TradeFill
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/trade/fills-history", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
