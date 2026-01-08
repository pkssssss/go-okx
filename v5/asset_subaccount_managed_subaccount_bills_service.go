package okx

import (
	"context"
	"net/http"
)

// AssetSubaccountManagedSubaccountBillsService 查询托管子账户转账记录（交易团队母账户）。
type AssetSubaccountManagedSubaccountBillsService struct {
	c *Client
	q assetSubaccountBillsQuery
}

// NewAssetSubaccountManagedSubaccountBillsService 创建 AssetSubaccountManagedSubaccountBillsService。
func (c *Client) NewAssetSubaccountManagedSubaccountBillsService() *AssetSubaccountManagedSubaccountBillsService {
	return &AssetSubaccountManagedSubaccountBillsService{c: c}
}

// Ccy 设置币种过滤。
func (s *AssetSubaccountManagedSubaccountBillsService) Ccy(ccy string) *AssetSubaccountManagedSubaccountBillsService {
	s.q.ccy = ccy
	return s
}

// Type 设置划转类型过滤（0=母账户转子账户，1=子账户转母账户）。
func (s *AssetSubaccountManagedSubaccountBillsService) Type(billType string) *AssetSubaccountManagedSubaccountBillsService {
	s.q.billType = billType
	return s
}

// SubAcct 设置子账户名称过滤。
func (s *AssetSubaccountManagedSubaccountBillsService) SubAcct(subAcct string) *AssetSubaccountManagedSubaccountBillsService {
	s.q.subAcct = subAcct
	return s
}

// SubUid 设置子账户 UID 过滤。
func (s *AssetSubaccountManagedSubaccountBillsService) SubUid(subUid string) *AssetSubaccountManagedSubaccountBillsService {
	s.q.subUid = subUid
	return s
}

// After 查询在此之前的内容（更旧的数据），值为 Unix 毫秒时间戳字符串。
func (s *AssetSubaccountManagedSubaccountBillsService) After(after string) *AssetSubaccountManagedSubaccountBillsService {
	s.q.after = after
	return s
}

// Before 查询在此之后的内容（更新的数据），值为 Unix 毫秒时间戳字符串。
func (s *AssetSubaccountManagedSubaccountBillsService) Before(before string) *AssetSubaccountManagedSubaccountBillsService {
	s.q.before = before
	return s
}

// Limit 设置返回条数（最大 100，默认 100）。
func (s *AssetSubaccountManagedSubaccountBillsService) Limit(limit int) *AssetSubaccountManagedSubaccountBillsService {
	s.q.limit = &limit
	return s
}

// Do 查询托管子账户转账记录（GET /api/v5/asset/subaccount/managed-subaccount-bills）。
func (s *AssetSubaccountManagedSubaccountBillsService) Do(ctx context.Context) ([]AssetSubaccountBill, error) {
	var data []AssetSubaccountBill
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/asset/subaccount/managed-subaccount-bills", s.q.values(), nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
