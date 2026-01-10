package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// DeliveryExerciseDetail 表示交割/行权记录的明细项。
type DeliveryExerciseDetail struct {
	InsId string `json:"insId"`
	Px    string `json:"px"`
	Type  string `json:"type"`
}

// DeliveryExerciseHistory 表示交割和行权记录。
//
// 说明：价格字段保持为 string（无损）。
type DeliveryExerciseHistory struct {
	TS      int64                    `json:"ts,string"`
	Details []DeliveryExerciseDetail `json:"details"`
}

// PublicDeliveryExerciseHistoryService 获取交割和行权记录。
type PublicDeliveryExerciseHistoryService struct {
	c *Client

	instType   string
	instFamily string
	after      string
	before     string
	limit      *int
}

// NewPublicDeliveryExerciseHistoryService 创建 PublicDeliveryExerciseHistoryService。
func (c *Client) NewPublicDeliveryExerciseHistoryService() *PublicDeliveryExerciseHistoryService {
	return &PublicDeliveryExerciseHistoryService{c: c}
}

// InstType 设置产品类型（必填：FUTURES/OPTION）。
func (s *PublicDeliveryExerciseHistoryService) InstType(instType string) *PublicDeliveryExerciseHistoryService {
	s.instType = instType
	return s
}

// InstFamily 设置交易品种（必填，如 BTC-USD）。
func (s *PublicDeliveryExerciseHistoryService) InstFamily(instFamily string) *PublicDeliveryExerciseHistoryService {
	s.instFamily = instFamily
	return s
}

// Uly 设置标的指数（兼容官方示例/SDK 参数名，等价于 InstFamily）。
func (s *PublicDeliveryExerciseHistoryService) Uly(uly string) *PublicDeliveryExerciseHistoryService {
	s.instFamily = uly
	return s
}

// After 设置请求此 ts 之前（更旧的数据）的分页内容。
func (s *PublicDeliveryExerciseHistoryService) After(after string) *PublicDeliveryExerciseHistoryService {
	s.after = after
	return s
}

// Before 设置请求此 ts 之后（更新的数据）的分页内容。
func (s *PublicDeliveryExerciseHistoryService) Before(before string) *PublicDeliveryExerciseHistoryService {
	s.before = before
	return s
}

// Limit 设置返回条数（最大 100，默认 100）。
func (s *PublicDeliveryExerciseHistoryService) Limit(limit int) *PublicDeliveryExerciseHistoryService {
	s.limit = &limit
	return s
}

var (
	errPublicDeliveryExerciseHistoryMissingInstType   = errors.New("okx: public delivery exercise history requires instType")
	errPublicDeliveryExerciseHistoryMissingInstFamily = errors.New("okx: public delivery exercise history requires instFamily")
)

// Do 获取交割和行权记录（GET /api/v5/public/delivery-exercise-history）。
func (s *PublicDeliveryExerciseHistoryService) Do(ctx context.Context) ([]DeliveryExerciseHistory, error) {
	if s.instType == "" {
		return nil, errPublicDeliveryExerciseHistoryMissingInstType
	}
	if s.instFamily == "" {
		return nil, errPublicDeliveryExerciseHistoryMissingInstFamily
	}

	q := url.Values{}
	q.Set("instType", s.instType)
	q.Set("instFamily", s.instFamily)
	if s.after != "" {
		q.Set("after", s.after)
	}
	if s.before != "" {
		q.Set("before", s.before)
	}
	if s.limit != nil {
		q.Set("limit", strconv.Itoa(*s.limit))
	}

	var data []DeliveryExerciseHistory
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/public/delivery-exercise-history", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
