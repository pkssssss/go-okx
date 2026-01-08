package okx

import (
	"context"
	"net/http"
)

// AccountBillsArchiveService 查询交易账户账单流水（近三个月）。
type AccountBillsArchiveService struct {
	c *Client
	q accountBillsQuery
}

// NewAccountBillsArchiveService 创建 AccountBillsArchiveService。
func (c *Client) NewAccountBillsArchiveService() *AccountBillsArchiveService {
	return &AccountBillsArchiveService{c: c}
}

// InstType 设置产品类型（SPOT/MARGIN/SWAP/FUTURES/OPTION）。
func (s *AccountBillsArchiveService) InstType(instType string) *AccountBillsArchiveService {
	s.q.instType = instType
	return s
}

// InstId 设置产品 ID，如 BTC-USDT。
func (s *AccountBillsArchiveService) InstId(instId string) *AccountBillsArchiveService {
	s.q.instId = instId
	return s
}

// Ccy 设置币种过滤（支持多币种，逗号分隔；最多 20 个）。
func (s *AccountBillsArchiveService) Ccy(ccy string) *AccountBillsArchiveService {
	s.q.ccy = ccy
	return s
}

// MgnMode 设置仓位类型（isolated/cross）。
func (s *AccountBillsArchiveService) MgnMode(mgnMode string) *AccountBillsArchiveService {
	s.q.mgnMode = mgnMode
	return s
}

// CtType 设置合约类型（linear/inverse），仅交割/永续有效。
func (s *AccountBillsArchiveService) CtType(ctType string) *AccountBillsArchiveService {
	s.q.ctType = ctType
	return s
}

// Type 设置账单类型（如 2=交易，8=资金费 等）。
func (s *AccountBillsArchiveService) Type(billType string) *AccountBillsArchiveService {
	s.q.billType = billType
	return s
}

// SubType 设置账单子类型。
func (s *AccountBillsArchiveService) SubType(subType string) *AccountBillsArchiveService {
	s.q.subType = subType
	return s
}

// After 请求此 id 之前（更旧的数据）的分页内容，传 billId。
func (s *AccountBillsArchiveService) After(after string) *AccountBillsArchiveService {
	s.q.after = after
	return s
}

// Before 请求此 id 之后（更新的数据）的分页内容，传 billId。
func (s *AccountBillsArchiveService) Before(before string) *AccountBillsArchiveService {
	s.q.before = before
	return s
}

// Begin 筛选开始时间（Unix 毫秒字符串）。
func (s *AccountBillsArchiveService) Begin(begin string) *AccountBillsArchiveService {
	s.q.begin = begin
	return s
}

// End 筛选结束时间（Unix 毫秒字符串）。
func (s *AccountBillsArchiveService) End(end string) *AccountBillsArchiveService {
	s.q.end = end
	return s
}

// Limit 设置返回条数（最大 100，默认 100）。
func (s *AccountBillsArchiveService) Limit(limit int) *AccountBillsArchiveService {
	s.q.limit = &limit
	return s
}

// Do 查询交易账户账单流水（GET /api/v5/account/bills-archive）。
func (s *AccountBillsArchiveService) Do(ctx context.Context) ([]AccountBill, error) {
	var data []AccountBill
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/bills-archive", s.q.values(), nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
