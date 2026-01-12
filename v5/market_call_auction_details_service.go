package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// MarketCallAuctionDetails 表示集合竞价信息。
// 数值字段保持为 string（无损）。
type MarketCallAuctionDetails struct {
	InstId string `json:"instId"`

	EqPx        string `json:"eqPx"`
	MatchedSz   string `json:"matchedSz"`
	UnmatchedSz string `json:"unmatchedSz"`

	State string `json:"state"`

	AuctionEndTime int64 `json:"auctionEndTime,string"`
	TS             int64 `json:"ts,string"`
}

// MarketCallAuctionDetailsService 获取集合竞价相关信息。
type MarketCallAuctionDetailsService struct {
	c      *Client
	instId string
}

// NewMarketCallAuctionDetailsService 创建 MarketCallAuctionDetailsService。
func (c *Client) NewMarketCallAuctionDetailsService() *MarketCallAuctionDetailsService {
	return &MarketCallAuctionDetailsService{c: c}
}

// InstId 设置产品 ID（必填）。
func (s *MarketCallAuctionDetailsService) InstId(instId string) *MarketCallAuctionDetailsService {
	s.instId = instId
	return s
}

var (
	errMarketCallAuctionDetailsMissingInstId = errors.New("okx: market call auction details requires instId")
	errEmptyMarketCallAuctionDetailsResponse = errors.New("okx: empty market call auction details response")
)

// Do 获取集合竞价相关信息（GET /api/v5/market/call-auction-details）。
func (s *MarketCallAuctionDetailsService) Do(ctx context.Context) (*MarketCallAuctionDetails, error) {
	if s.instId == "" {
		return nil, errMarketCallAuctionDetailsMissingInstId
	}

	q := url.Values{}
	q.Set("instId", s.instId)

	var data []MarketCallAuctionDetails
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/market/call-auction-details", q, nil, false, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyMarketCallAuctionDetailsResponse
	}
	return &data[0], nil
}
