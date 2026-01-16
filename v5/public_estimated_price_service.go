package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// EstimatedPrice 表示预估交割/行权价格。
//
// 说明：settlePx 保持为 string（无损）。
type EstimatedPrice struct {
	InstType   string `json:"instType"`
	InstId     string `json:"instId"`
	SettlePx   string `json:"settlePx"`
	SettleType string `json:"settleType"`
	TS         int64  `json:"ts,string"`
}

// PublicEstimatedPriceService 获取预估交割/行权价格（交割/行权前一小时才有返回值）。
type PublicEstimatedPriceService struct {
	c *Client

	instId string
}

// NewPublicEstimatedPriceService 创建 PublicEstimatedPriceService。
func (c *Client) NewPublicEstimatedPriceService() *PublicEstimatedPriceService {
	return &PublicEstimatedPriceService{c: c}
}

// InstId 设置产品 ID（必填；仅适用于交割/期权），如 BTC-USD-200214。
func (s *PublicEstimatedPriceService) InstId(instId string) *PublicEstimatedPriceService {
	s.instId = instId
	return s
}

var errPublicEstimatedPriceMissingInstId = errors.New("okx: public estimated price requires instId")

// Do 获取预估交割/行权价格（GET /api/v5/public/estimated-price）。
func (s *PublicEstimatedPriceService) Do(ctx context.Context) ([]EstimatedPrice, error) {
	if s.instId == "" {
		return nil, errPublicEstimatedPriceMissingInstId
	}

	q := url.Values{}
	q.Set("instId", s.instId)

	var data []EstimatedPrice
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/public/estimated-price", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
