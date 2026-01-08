package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// OpenInterest 表示持仓总量（未平仓量）。
type OpenInterest struct {
	InstType string `json:"instType"`
	InstId   string `json:"instId"`

	OI    string `json:"oi"`
	OICcy string `json:"oiCcy"`
	OIUsd string `json:"oiUsd"`

	TS int64 `json:"ts,string"`
}

// PublicOpenInterestService 查询持仓总量（未平仓量）。
type PublicOpenInterestService struct {
	c *Client

	instType   string
	uly        string
	instFamily string
	instId     string
}

// NewPublicOpenInterestService 创建 PublicOpenInterestService。
func (c *Client) NewPublicOpenInterestService() *PublicOpenInterestService {
	return &PublicOpenInterestService{c: c}
}

// InstType 设置产品类型（SWAP/FUTURES/OPTION），必填。
func (s *PublicOpenInterestService) InstType(instType string) *PublicOpenInterestService {
	s.instType = instType
	return s
}

// Uly 设置标的指数（适用于交割/永续/期权）。
func (s *PublicOpenInterestService) Uly(uly string) *PublicOpenInterestService {
	s.uly = uly
	return s
}

// InstFamily 设置交易品种（适用于交割/永续/期权），如 BTC-USD。
func (s *PublicOpenInterestService) InstFamily(instFamily string) *PublicOpenInterestService {
	s.instFamily = instFamily
	return s
}

// InstId 设置产品 ID。
func (s *PublicOpenInterestService) InstId(instId string) *PublicOpenInterestService {
	s.instId = instId
	return s
}

var errPublicOpenInterestMissingInstType = errors.New("okx: public open interest requires instType")

// Do 查询持仓总量（GET /api/v5/public/open-interest）。
func (s *PublicOpenInterestService) Do(ctx context.Context) ([]OpenInterest, error) {
	if s.instType == "" {
		return nil, errPublicOpenInterestMissingInstType
	}

	q := url.Values{}
	q.Set("instType", s.instType)
	if s.uly != "" {
		q.Set("uly", s.uly)
	}
	if s.instFamily != "" {
		q.Set("instFamily", s.instFamily)
	}
	if s.instId != "" {
		q.Set("instId", s.instId)
	}

	var data []OpenInterest
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/public/open-interest", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
