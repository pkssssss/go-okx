package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// MarkPrice 表示标记价格。
type MarkPrice struct {
	InstType string `json:"instType"`
	InstId   string `json:"instId"`
	MarkPx   string `json:"markPx"`
	TS       int64  `json:"ts,string"`
}

// PublicMarkPriceService 查询标记价格。
type PublicMarkPriceService struct {
	c *Client

	instType   string
	uly        string
	instFamily string
	instId     string
}

// NewPublicMarkPriceService 创建 PublicMarkPriceService。
func (c *Client) NewPublicMarkPriceService() *PublicMarkPriceService {
	return &PublicMarkPriceService{c: c}
}

// InstType 设置产品类型（SWAP/FUTURES/OPTION），必填。
func (s *PublicMarkPriceService) InstType(instType string) *PublicMarkPriceService {
	s.instType = instType
	return s
}

// Uly 设置标的指数（适用于交割/永续/期权）。
func (s *PublicMarkPriceService) Uly(uly string) *PublicMarkPriceService {
	s.uly = uly
	return s
}

// InstFamily 设置交易品种（适用于交割/永续/期权），如 BTC-USD。
func (s *PublicMarkPriceService) InstFamily(instFamily string) *PublicMarkPriceService {
	s.instFamily = instFamily
	return s
}

// InstId 设置产品 ID。
func (s *PublicMarkPriceService) InstId(instId string) *PublicMarkPriceService {
	s.instId = instId
	return s
}

var errPublicMarkPriceMissingInstType = errors.New("okx: public mark price requires instType")

// Do 查询标记价格（GET /api/v5/public/mark-price）。
func (s *PublicMarkPriceService) Do(ctx context.Context) ([]MarkPrice, error) {
	if s.instType == "" {
		return nil, errPublicMarkPriceMissingInstType
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

	var data []MarkPrice
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/public/mark-price", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
