package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type assetBillsHistoryQuery struct {
	ccy        string
	billType   string
	clientId   string
	after      string
	before     string
	limit      *int
	pagingType string
}

func (q assetBillsHistoryQuery) values() url.Values {
	v := url.Values{}
	if q.ccy != "" {
		v.Set("ccy", q.ccy)
	}
	if q.billType != "" {
		v.Set("type", q.billType)
	}
	if q.clientId != "" {
		v.Set("clientId", q.clientId)
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
	if q.pagingType != "" {
		v.Set("pagingType", q.pagingType)
	}

	if len(v) == 0 {
		return nil
	}
	return v
}

// AssetBillsHistoryService 获取资金流水全历史（可追溯至 2021-02-01）。
type AssetBillsHistoryService struct {
	c *Client
	q assetBillsHistoryQuery
}

// NewAssetBillsHistoryService 创建 AssetBillsHistoryService。
func (c *Client) NewAssetBillsHistoryService() *AssetBillsHistoryService {
	return &AssetBillsHistoryService{c: c}
}

// Ccy 设置币种过滤。
func (s *AssetBillsHistoryService) Ccy(ccy string) *AssetBillsHistoryService {
	s.q.ccy = ccy
	return s
}

// Type 设置账单类型（如 1=充值，2=提现，130=从交易账户转入，131=转出至交易账户 等）。
func (s *AssetBillsHistoryService) Type(billType string) *AssetBillsHistoryService {
	s.q.billType = billType
	return s
}

// ClientId 设置转账或提币的客户自定义 ID（1-32）。
func (s *AssetBillsHistoryService) ClientId(clientId string) *AssetBillsHistoryService {
	s.q.clientId = clientId
	return s
}

// After 查询在此之前的内容（更旧的数据）。
// pagingType=1 时为时间戳（Unix 毫秒字符串）；pagingType=2 时为账单记录 ID。
func (s *AssetBillsHistoryService) After(after string) *AssetBillsHistoryService {
	s.q.after = after
	return s
}

// Before 查询在此之后的内容（更新的数据），值为时间戳（Unix 毫秒字符串）。
func (s *AssetBillsHistoryService) Before(before string) *AssetBillsHistoryService {
	s.q.before = before
	return s
}

// Limit 设置返回条数（最大 100，默认 100）。
func (s *AssetBillsHistoryService) Limit(limit int) *AssetBillsHistoryService {
	s.q.limit = &limit
	return s
}

// PagingType 设置分页类型：1=按时间戳分页；2=按账单记录 ID 分页（默认 1）。
func (s *AssetBillsHistoryService) PagingType(pagingType string) *AssetBillsHistoryService {
	s.q.pagingType = pagingType
	return s
}

// Do 获取资金流水全历史（GET /api/v5/asset/bills-history）。
func (s *AssetBillsHistoryService) Do(ctx context.Context) ([]AssetBill, error) {
	var data []AssetBill
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/asset/bills-history", s.q.values(), nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
