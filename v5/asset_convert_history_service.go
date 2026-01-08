package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type assetConvertHistoryQuery struct {
	clTReqId string
	after    string
	before   string
	limit    *int
	tag      string
}

func (q assetConvertHistoryQuery) values() url.Values {
	v := url.Values{}
	if q.clTReqId != "" {
		v.Set("clTReqId", q.clTReqId)
	}
	if q.after != "" {
		v.Set("after", q.after)
	}
	if q.before != "" {
		v.Set("before", q.before)
	}
	if q.limit != nil {
		v.Set("limit", strconv.Itoa(*q.limit))
	}
	if q.tag != "" {
		v.Set("tag", q.tag)
	}

	if len(v) == 0 {
		return nil
	}
	return v
}

// AssetConvertHistoryService 获取闪兑交易历史。
type AssetConvertHistoryService struct {
	c *Client
	q assetConvertHistoryQuery
}

// NewAssetConvertHistoryService 创建 AssetConvertHistoryService。
func (c *Client) NewAssetConvertHistoryService() *AssetConvertHistoryService {
	return &AssetConvertHistoryService{c: c}
}

// ClTReqId 设置用户自定义订单标识过滤（可选：1-32）。
func (s *AssetConvertHistoryService) ClTReqId(clTReqId string) *AssetConvertHistoryService {
	s.q.clTReqId = clTReqId
	return s
}

// After 查询在此之前的内容（更旧的数据），值为 Unix 毫秒时间戳字符串。
func (s *AssetConvertHistoryService) After(after string) *AssetConvertHistoryService {
	s.q.after = after
	return s
}

// Before 查询在此之后的内容（更新的数据），值为 Unix 毫秒时间戳字符串。
func (s *AssetConvertHistoryService) Before(before string) *AssetConvertHistoryService {
	s.q.before = before
	return s
}

// Limit 设置返回条数（最大 100，默认 100）。
func (s *AssetConvertHistoryService) Limit(limit int) *AssetConvertHistoryService {
	s.q.limit = &limit
	return s
}

// Tag 设置订单标签过滤（可选，适用于 broker）。
func (s *AssetConvertHistoryService) Tag(tag string) *AssetConvertHistoryService {
	s.q.tag = tag
	return s
}

// Do 获取闪兑交易历史（GET /api/v5/asset/convert/history）。
func (s *AssetConvertHistoryService) Do(ctx context.Context) ([]AssetConvertTrade, error) {
	var data []AssetConvertTrade
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/asset/convert/history", s.q.values(), nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
