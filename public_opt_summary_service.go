package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// OptSummary 表示期权行情概要（部分字段）。
//
// 说明：该端点字段较多且多为小数，SDK 侧统一保持 string（无损）。
type OptSummary struct {
	InstType string `json:"instType"`
	InstId   string `json:"instId"`
	Uly      string `json:"uly"`

	AskVol  string `json:"askVol"`
	BidVol  string `json:"bidVol"`
	MarkVol string `json:"markVol"`
	RealVol string `json:"realVol"`

	Delta string `json:"delta"`
	Gamma string `json:"gamma"`
	Theta string `json:"theta"`
	Vega  string `json:"vega"`

	VolLv    string `json:"volLv"`
	FwdPx    string `json:"fwdPx"`
	Distance string `json:"distance"`

	TS int64 `json:"ts,string"`
}

// PublicOptSummaryService 查询期权行情概要。
type PublicOptSummaryService struct {
	c *Client

	uly     string
	expTime *int
}

// NewPublicOptSummaryService 创建 PublicOptSummaryService。
func (c *Client) NewPublicOptSummaryService() *PublicOptSummaryService {
	return &PublicOptSummaryService{c: c}
}

// Uly 设置标的指数（必填），如 BTC-USD。
func (s *PublicOptSummaryService) Uly(uly string) *PublicOptSummaryService {
	s.uly = uly
	return s
}

// ExpTime 设置到期日（可选，格式 YYMMDD，如 260123）。
func (s *PublicOptSummaryService) ExpTime(expTime int) *PublicOptSummaryService {
	s.expTime = &expTime
	return s
}

var errPublicOptSummaryMissingUly = errors.New("okx: public opt summary requires uly")

// Do 查询期权行情概要（GET /api/v5/public/opt-summary）。
func (s *PublicOptSummaryService) Do(ctx context.Context) ([]OptSummary, error) {
	if s.uly == "" {
		return nil, errPublicOptSummaryMissingUly
	}

	q := url.Values{}
	q.Set("uly", s.uly)
	if s.expTime != nil {
		q.Set("expTime", strconv.Itoa(*s.expTime))
	}

	var data []OptSummary
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/public/opt-summary", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
