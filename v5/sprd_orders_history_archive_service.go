package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// SprdOrdersHistoryArchiveService 获取历史订单记录（近三月）。
type SprdOrdersHistoryArchiveService struct {
	c *Client

	sprdId     string
	ordType    string
	state      string
	instType   string
	instFamily string

	beginId string
	endId   string

	begin string
	end   string

	limit *int
}

// NewSprdOrdersHistoryArchiveService 创建 SprdOrdersHistoryArchiveService。
func (c *Client) NewSprdOrdersHistoryArchiveService() *SprdOrdersHistoryArchiveService {
	return &SprdOrdersHistoryArchiveService{c: c}
}

func (s *SprdOrdersHistoryArchiveService) SprdId(sprdId string) *SprdOrdersHistoryArchiveService {
	s.sprdId = sprdId
	return s
}

func (s *SprdOrdersHistoryArchiveService) OrdType(ordType string) *SprdOrdersHistoryArchiveService {
	s.ordType = ordType
	return s
}

func (s *SprdOrdersHistoryArchiveService) State(state string) *SprdOrdersHistoryArchiveService {
	s.state = state
	return s
}

// InstType 设置产品类型（SPOT/FUTURES/SWAP）。
func (s *SprdOrdersHistoryArchiveService) InstType(instType string) *SprdOrdersHistoryArchiveService {
	s.instType = instType
	return s
}

// InstFamily 设置交易品种（如 BTC-USDT）。
func (s *SprdOrdersHistoryArchiveService) InstFamily(instFamily string) *SprdOrdersHistoryArchiveService {
	s.instFamily = instFamily
	return s
}

func (s *SprdOrdersHistoryArchiveService) BeginId(beginId string) *SprdOrdersHistoryArchiveService {
	s.beginId = beginId
	return s
}

func (s *SprdOrdersHistoryArchiveService) EndId(endId string) *SprdOrdersHistoryArchiveService {
	s.endId = endId
	return s
}

func (s *SprdOrdersHistoryArchiveService) Begin(begin string) *SprdOrdersHistoryArchiveService {
	s.begin = begin
	return s
}

func (s *SprdOrdersHistoryArchiveService) End(end string) *SprdOrdersHistoryArchiveService {
	s.end = end
	return s
}

func (s *SprdOrdersHistoryArchiveService) Limit(limit int) *SprdOrdersHistoryArchiveService {
	s.limit = &limit
	return s
}

// Do 获取历史订单记录（GET /api/v5/sprd/orders-history-archive）。
func (s *SprdOrdersHistoryArchiveService) Do(ctx context.Context) ([]SprdOrder, error) {
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
	if s.instType != "" {
		q.Set("instType", s.instType)
	}
	if s.instFamily != "" {
		q.Set("instFamily", s.instFamily)
	}
	if s.beginId != "" {
		q.Set("beginId", s.beginId)
	}
	if s.endId != "" {
		q.Set("endId", s.endId)
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

	var data []SprdOrder
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/sprd/orders-history-archive", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
