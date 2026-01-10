package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// EstimatedSettlementInfo 表示交割预估结算价格。
//
// 说明：estSettlePx 保持为 string（无损）。
type EstimatedSettlementInfo struct {
	InstId         string `json:"instId"`
	NextSettleTime int64  `json:"nextSettleTime,string"`
	EstSettlePx    string `json:"estSettlePx"`
	TS             int64  `json:"ts,string"`
}

// PublicEstimatedSettlementInfoService 获取交割预估结算价格（结算前一小时才有返回值）。
type PublicEstimatedSettlementInfoService struct {
	c *Client

	instId string
}

// NewPublicEstimatedSettlementInfoService 创建 PublicEstimatedSettlementInfoService。
func (c *Client) NewPublicEstimatedSettlementInfoService() *PublicEstimatedSettlementInfoService {
	return &PublicEstimatedSettlementInfoService{c: c}
}

// InstId 设置产品 ID（必填；仅适用于交割），如 XRP-USDT-250307。
func (s *PublicEstimatedSettlementInfoService) InstId(instId string) *PublicEstimatedSettlementInfoService {
	s.instId = instId
	return s
}

var errPublicEstimatedSettlementInfoMissingInstId = errors.New("okx: public estimated settlement info requires instId")

// Do 获取交割预估结算价格（GET /api/v5/public/estimated-settlement-info）。
func (s *PublicEstimatedSettlementInfoService) Do(ctx context.Context) ([]EstimatedSettlementInfo, error) {
	if s.instId == "" {
		return nil, errPublicEstimatedSettlementInfoMissingInstId
	}

	q := url.Values{}
	q.Set("instId", s.instId)

	var data []EstimatedSettlementInfo
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/public/estimated-settlement-info", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
