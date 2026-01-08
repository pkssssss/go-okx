package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// FundingRate 表示资金费率信息。
type FundingRate struct {
	InstType string `json:"instType"`
	InstId   string `json:"instId"`

	FundingRate string `json:"fundingRate"`
	FundingTime int64  `json:"fundingTime,string"`

	NextFundingRate string `json:"nextFundingRate"`
	NextFundingTime int64  `json:"nextFundingTime,string"`

	MinFundingRate string `json:"minFundingRate"`
	MaxFundingRate string `json:"maxFundingRate"`

	InterestRate string `json:"interestRate"`
	Premium      string `json:"premium"`

	SettFundingRate string `json:"settFundingRate"`

	Method      string `json:"method"`
	SettState   string `json:"settState"`
	FormulaType string `json:"formulaType"`

	TS int64 `json:"ts,string"`
}

// PublicFundingRateService 查询资金费率信息。
type PublicFundingRateService struct {
	c *Client

	instId string
}

// NewPublicFundingRateService 创建 PublicFundingRateService。
func (c *Client) NewPublicFundingRateService() *PublicFundingRateService {
	return &PublicFundingRateService{c: c}
}

// InstId 设置产品 ID（必填，通常为 SWAP）。
func (s *PublicFundingRateService) InstId(instId string) *PublicFundingRateService {
	s.instId = instId
	return s
}

var (
	errPublicFundingRateMissingInstId = errors.New("okx: public funding rate requires instId")
	errEmptyFundingRateResponse       = errors.New("okx: empty funding rate response")
)

// Do 查询资金费率信息（GET /api/v5/public/funding-rate）。
func (s *PublicFundingRateService) Do(ctx context.Context) (*FundingRate, error) {
	if s.instId == "" {
		return nil, errPublicFundingRateMissingInstId
	}

	q := url.Values{}
	q.Set("instId", s.instId)

	var data []FundingRate
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/public/funding-rate", q, nil, false, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyFundingRateResponse
	}
	return &data[0], nil
}
