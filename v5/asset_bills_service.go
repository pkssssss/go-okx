package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// AssetBill 表示资金账户账单流水（近一个月）。
// 数值字段保持为 string（无损）。
type AssetBill struct {
	BillId   string `json:"billId"`
	Ccy      string `json:"ccy"`
	ClientId string `json:"clientId"`
	BalChg   string `json:"balChg"`
	Bal      string `json:"bal"`
	Type     string `json:"type"`
	Notes    string `json:"notes"`
	TS       int64  `json:"ts,string"`
}

type assetBillsQuery struct {
	ccy      string
	billType string
	clientId string
	after    string
	before   string
	limit    *int
}

func (q assetBillsQuery) values() url.Values {
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

	if len(v) == 0 {
		return nil
	}
	return v
}

// AssetBillsService 获取资金流水（近一个月）。
type AssetBillsService struct {
	c *Client
	q assetBillsQuery
}

// NewAssetBillsService 创建 AssetBillsService。
func (c *Client) NewAssetBillsService() *AssetBillsService {
	return &AssetBillsService{c: c}
}

// Ccy 设置币种过滤。
func (s *AssetBillsService) Ccy(ccy string) *AssetBillsService {
	s.q.ccy = ccy
	return s
}

// Type 设置账单类型（如 1=充值，2=提现，130=从交易账户转入，131=转出至交易账户 等）。
func (s *AssetBillsService) Type(billType string) *AssetBillsService {
	s.q.billType = billType
	return s
}

// ClientId 设置转账或提币的客户自定义 ID（1-32）。
func (s *AssetBillsService) ClientId(clientId string) *AssetBillsService {
	s.q.clientId = clientId
	return s
}

// After 查询在此之前的内容（更旧的数据），值为时间戳（Unix 毫秒字符串）。
func (s *AssetBillsService) After(after string) *AssetBillsService {
	s.q.after = after
	return s
}

// Before 查询在此之后的内容（更新的数据），值为时间戳（Unix 毫秒字符串）。
func (s *AssetBillsService) Before(before string) *AssetBillsService {
	s.q.before = before
	return s
}

// Limit 设置返回条数（最大 100，默认 100）。
func (s *AssetBillsService) Limit(limit int) *AssetBillsService {
	s.q.limit = &limit
	return s
}

// Do 获取资金流水（GET /api/v5/asset/bills）。
func (s *AssetBillsService) Do(ctx context.Context) ([]AssetBill, error) {
	var data []AssetBill
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/asset/bills", s.q.values(), nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
