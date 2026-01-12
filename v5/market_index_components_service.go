package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// MarketIndexComponent 表示指数成分。
// 数值字段保持为 string（无损）。
type MarketIndexComponent struct {
	Exch   string `json:"exch"`
	Symbol string `json:"symbol"`
	SymPx  string `json:"symPx"`
	Wgt    string `json:"wgt"`
	CnvPx  string `json:"cnvPx"`
}

// MarketIndexComponents 表示指数成分数据。
// 数值字段保持为 string（无损）。
type MarketIndexComponents struct {
	Components []MarketIndexComponent `json:"components"`

	Last  string `json:"last"`
	Index string `json:"index"`
	TS    int64  `json:"ts,string"`
}

// MarketIndexComponentsService 获取指数成分数据。
type MarketIndexComponentsService struct {
	c     *Client
	index string
}

// NewMarketIndexComponentsService 创建 MarketIndexComponentsService。
func (c *Client) NewMarketIndexComponentsService() *MarketIndexComponentsService {
	return &MarketIndexComponentsService{c: c}
}

// Index 设置指数（必填），如 BTC-USD，与 uly 含义相同。
func (s *MarketIndexComponentsService) Index(index string) *MarketIndexComponentsService {
	s.index = index
	return s
}

var (
	errMarketIndexComponentsMissingIndex  = errors.New("okx: market index components requires index")
	errEmptyMarketIndexComponentsResponse = errors.New("okx: empty market index components response")
)

// Do 获取指数成分数据（GET /api/v5/market/index-components）。
func (s *MarketIndexComponentsService) Do(ctx context.Context) (*MarketIndexComponents, error) {
	if s.index == "" {
		return nil, errMarketIndexComponentsMissingIndex
	}

	q := url.Values{}
	q.Set("index", s.index)

	var data MarketIndexComponents
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/market/index-components", q, nil, false, &data); err != nil {
		return nil, err
	}
	if data.Index == "" && data.Last == "" && data.TS == 0 && len(data.Components) == 0 {
		return nil, errEmptyMarketIndexComponentsResponse
	}
	return &data, nil
}
