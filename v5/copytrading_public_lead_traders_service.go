package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// CopyTradingPublicLeadTradersService 获取交易员排名（公共）。
type CopyTradingPublicLeadTradersService struct {
	c *Client

	instType    string
	sortType    string
	state       string
	minLeadDays string
	minAssets   string
	maxAssets   string
	minAum      string
	maxAum      string
	dataVer     string
	page        *int
	limit       *int
}

// NewCopyTradingPublicLeadTradersService 创建 CopyTradingPublicLeadTradersService。
func (c *Client) NewCopyTradingPublicLeadTradersService() *CopyTradingPublicLeadTradersService {
	return &CopyTradingPublicLeadTradersService{c: c}
}

func (s *CopyTradingPublicLeadTradersService) InstType(instType string) *CopyTradingPublicLeadTradersService {
	s.instType = instType
	return s
}

func (s *CopyTradingPublicLeadTradersService) SortType(sortType string) *CopyTradingPublicLeadTradersService {
	s.sortType = sortType
	return s
}

func (s *CopyTradingPublicLeadTradersService) State(state string) *CopyTradingPublicLeadTradersService {
	s.state = state
	return s
}

func (s *CopyTradingPublicLeadTradersService) MinLeadDays(minLeadDays string) *CopyTradingPublicLeadTradersService {
	s.minLeadDays = minLeadDays
	return s
}

func (s *CopyTradingPublicLeadTradersService) MinAssets(minAssets string) *CopyTradingPublicLeadTradersService {
	s.minAssets = minAssets
	return s
}

func (s *CopyTradingPublicLeadTradersService) MaxAssets(maxAssets string) *CopyTradingPublicLeadTradersService {
	s.maxAssets = maxAssets
	return s
}

func (s *CopyTradingPublicLeadTradersService) MinAum(minAum string) *CopyTradingPublicLeadTradersService {
	s.minAum = minAum
	return s
}

func (s *CopyTradingPublicLeadTradersService) MaxAum(maxAum string) *CopyTradingPublicLeadTradersService {
	s.maxAum = maxAum
	return s
}

func (s *CopyTradingPublicLeadTradersService) DataVer(dataVer string) *CopyTradingPublicLeadTradersService {
	s.dataVer = dataVer
	return s
}

func (s *CopyTradingPublicLeadTradersService) Page(page int) *CopyTradingPublicLeadTradersService {
	s.page = &page
	return s
}

func (s *CopyTradingPublicLeadTradersService) Limit(limit int) *CopyTradingPublicLeadTradersService {
	s.limit = &limit
	return s
}

var errEmptyCopyTradingPublicLeadTradersResponse = errors.New("okx: empty copytrading public lead traders response")

// Do 获取交易员排名（GET /api/v5/copytrading/public-lead-traders）。
func (s *CopyTradingPublicLeadTradersService) Do(ctx context.Context) (*CopyTradingPublicLeadTraders, error) {
	v := url.Values{}
	if s.instType != "" {
		v.Set("instType", s.instType)
	}
	if s.sortType != "" {
		v.Set("sortType", s.sortType)
	}
	if s.state != "" {
		v.Set("state", s.state)
	}
	if s.minLeadDays != "" {
		v.Set("minLeadDays", s.minLeadDays)
	}
	if s.minAssets != "" {
		v.Set("minAssets", s.minAssets)
	}
	if s.maxAssets != "" {
		v.Set("maxAssets", s.maxAssets)
	}
	if s.minAum != "" {
		v.Set("minAum", s.minAum)
	}
	if s.maxAum != "" {
		v.Set("maxAum", s.maxAum)
	}
	if s.dataVer != "" {
		v.Set("dataVer", s.dataVer)
	}
	if s.page != nil {
		v.Set("page", strconv.Itoa(*s.page))
	}
	if s.limit != nil {
		v.Set("limit", strconv.Itoa(*s.limit))
	}
	if len(v) == 0 {
		v = nil
	}

	var data []CopyTradingPublicLeadTraders
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/copytrading/public-lead-traders", v, nil, false, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyCopyTradingPublicLeadTradersResponse
	}
	return &data[0], nil
}
