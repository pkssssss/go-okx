package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// AssetSubaccountBill 表示子账户转账记录项（含托管子账户）。
// 数值字段保持为 string（无损）。
type AssetSubaccountBill struct {
	BillId  string `json:"billId"`
	Type    string `json:"type"`
	Ccy     string `json:"ccy"`
	Amt     string `json:"amt"`
	SubAcct string `json:"subAcct"`
	SubUid  string `json:"subUid"`
	TS      int64  `json:"ts,string"`
}

type assetSubaccountBillsQuery struct {
	ccy      string
	billType string
	subAcct  string
	subUid   string
	after    string
	before   string
	limit    *int
}

func (q assetSubaccountBillsQuery) values() url.Values {
	v := url.Values{}
	if q.ccy != "" {
		v.Set("ccy", q.ccy)
	}
	if q.billType != "" {
		v.Set("type", q.billType)
	}
	if q.subAcct != "" {
		v.Set("subAcct", q.subAcct)
	}
	if q.subUid != "" {
		v.Set("subUid", q.subUid)
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

// AssetSubaccountBillsService 查询子账户转账记录（母账户）。
type AssetSubaccountBillsService struct {
	c *Client
	q assetSubaccountBillsQuery
}

// NewAssetSubaccountBillsService 创建 AssetSubaccountBillsService。
func (c *Client) NewAssetSubaccountBillsService() *AssetSubaccountBillsService {
	return &AssetSubaccountBillsService{c: c}
}

// Ccy 设置币种过滤。
func (s *AssetSubaccountBillsService) Ccy(ccy string) *AssetSubaccountBillsService {
	s.q.ccy = ccy
	return s
}

// Type 设置划转类型过滤（0=母账户转子账户，1=子账户转母账户）。
func (s *AssetSubaccountBillsService) Type(billType string) *AssetSubaccountBillsService {
	s.q.billType = billType
	return s
}

// SubAcct 设置子账户名称过滤。
func (s *AssetSubaccountBillsService) SubAcct(subAcct string) *AssetSubaccountBillsService {
	s.q.subAcct = subAcct
	return s
}

// After 查询在此之前的内容（更旧的数据），值为 Unix 毫秒时间戳字符串。
func (s *AssetSubaccountBillsService) After(after string) *AssetSubaccountBillsService {
	s.q.after = after
	return s
}

// Before 查询在此之后的内容（更新的数据），值为 Unix 毫秒时间戳字符串。
func (s *AssetSubaccountBillsService) Before(before string) *AssetSubaccountBillsService {
	s.q.before = before
	return s
}

// Limit 设置返回条数（最大 100，默认 100）。
func (s *AssetSubaccountBillsService) Limit(limit int) *AssetSubaccountBillsService {
	s.q.limit = &limit
	return s
}

// Do 查询子账户转账记录（GET /api/v5/asset/subaccount/bills）。
func (s *AssetSubaccountBillsService) Do(ctx context.Context) ([]AssetSubaccountBill, error) {
	var data []AssetSubaccountBill
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/asset/subaccount/bills", s.q.values(), nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
