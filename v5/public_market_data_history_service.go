package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// MarketDataHistoryFile 表示历史市场数据文件信息。
//
// 说明：时间戳按 Unix 毫秒解析为 int64；文件大小字段保持为 string（无损）。
type MarketDataHistoryFile struct {
	Filename string `json:"filename"`
	DateTs   int64  `json:"dateTs,string"`
	SizeMB   string `json:"sizeMB"`
	URL      string `json:"url"`
}

// MarketDataHistoryGroup 表示按产品与日期聚合后的数据组。
//
// 说明：时间戳按 Unix 毫秒解析为 int64；大小字段保持为 string（无损）。
type MarketDataHistoryGroup struct {
	InstId     string `json:"instId"`
	InstFamily string `json:"instFamily"`
	InstType   string `json:"instType"`

	DateRangeStart int64 `json:"dateRangeStart,string"`
	DateRangeEnd   int64 `json:"dateRangeEnd,string"`

	GroupSizeMB  string                  `json:"groupSizeMB"`
	GroupDetails []MarketDataHistoryFile `json:"groupDetails"`
}

// MarketDataHistory 表示历史市场数据响应。
type MarketDataHistory struct {
	DateAggrType string                   `json:"dateAggrType"`
	Details      []MarketDataHistoryGroup `json:"details"`
	TotalSizeMB  string                   `json:"totalSizeMB"`
	TS           int64                    `json:"ts,string"`
}

// PublicMarketDataHistoryService 获取历史市场数据。
type PublicMarketDataHistoryService struct {
	c *Client

	module         string
	instType       string
	instIdList     string
	instFamilyList string
	dateAggrType   string
	begin          string
	end            string
}

// NewPublicMarketDataHistoryService 创建 PublicMarketDataHistoryService。
func (c *Client) NewPublicMarketDataHistoryService() *PublicMarketDataHistoryService {
	return &PublicMarketDataHistoryService{c: c}
}

// Module 设置数据模块类型（必填：1/2/3/4/5/6）。
func (s *PublicMarketDataHistoryService) Module(module string) *PublicMarketDataHistoryService {
	s.module = module
	return s
}

// InstType 设置产品类型（必填：SPOT/FUTURES/SWAP/OPTION）。
func (s *PublicMarketDataHistoryService) InstType(instType string) *PublicMarketDataHistoryService {
	s.instType = instType
	return s
}

// InstIdList 设置产品 ID 列表（仅适用于 instType=SPOT；英文逗号分隔）。
func (s *PublicMarketDataHistoryService) InstIdList(instIdList string) *PublicMarketDataHistoryService {
	s.instIdList = instIdList
	return s
}

// InstFamilyList 设置交易品种列表（仅适用于 instType!=SPOT；英文逗号分隔）。
func (s *PublicMarketDataHistoryService) InstFamilyList(instFamilyList string) *PublicMarketDataHistoryService {
	s.instFamilyList = instFamilyList
	return s
}

// DateAggrType 设置日期聚合类型（必填：daily/monthly）。
func (s *PublicMarketDataHistoryService) DateAggrType(dateAggrType string) *PublicMarketDataHistoryService {
	s.dateAggrType = dateAggrType
	return s
}

// Begin 设置开始时间戳（必填；Unix 毫秒）。
func (s *PublicMarketDataHistoryService) Begin(begin string) *PublicMarketDataHistoryService {
	s.begin = begin
	return s
}

// End 设置结束时间戳（必填；Unix 毫秒）。
func (s *PublicMarketDataHistoryService) End(end string) *PublicMarketDataHistoryService {
	s.end = end
	return s
}

var (
	errPublicMarketDataHistoryMissingRequired       = errors.New("okx: public market data history requires module, instType, dateAggrType, begin, and end")
	errPublicMarketDataHistoryMissingInstIdList     = errors.New("okx: public market data history requires instIdList for SPOT")
	errPublicMarketDataHistoryMissingInstFamilyList = errors.New("okx: public market data history requires instFamilyList for non-SPOT")
)

// Do 获取历史市场数据（GET /api/v5/public/market-data-history）。
func (s *PublicMarketDataHistoryService) Do(ctx context.Context) ([]MarketDataHistory, error) {
	if s.module == "" || s.instType == "" || s.dateAggrType == "" || s.begin == "" || s.end == "" {
		return nil, errPublicMarketDataHistoryMissingRequired
	}
	if s.instType == "SPOT" {
		if s.instIdList == "" {
			return nil, errPublicMarketDataHistoryMissingInstIdList
		}
	} else {
		if s.instFamilyList == "" {
			return nil, errPublicMarketDataHistoryMissingInstFamilyList
		}
	}

	q := url.Values{}
	q.Set("module", s.module)
	q.Set("instType", s.instType)
	if s.instIdList != "" {
		q.Set("instIdList", s.instIdList)
	}
	if s.instFamilyList != "" {
		q.Set("instFamilyList", s.instFamilyList)
	}
	q.Set("dateAggrType", s.dateAggrType)
	q.Set("begin", s.begin)
	q.Set("end", s.end)

	var data []MarketDataHistory
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/public/market-data-history", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
