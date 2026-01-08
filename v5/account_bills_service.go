package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// AccountBill 表示交易账户账单流水（近七天/近三个月）。
// 数值字段保持为 string（无损）。
type AccountBill struct {
	InstType string `json:"instType"`
	InstId   string `json:"instId"`

	BillId  string `json:"billId"`
	Type    string `json:"type"`
	SubType string `json:"subType"`

	Ccy    string `json:"ccy"`
	Bal    string `json:"bal"`
	BalChg string `json:"balChg"`

	OrdId   string `json:"ordId"`
	ClOrdId string `json:"clOrdId"`
	TradeId string `json:"tradeId"`
	Tag     string `json:"tag"`

	Px       string `json:"px"`
	Sz       string `json:"sz"`
	ExecType string `json:"execType"`

	Fee      string `json:"fee"`
	Interest string `json:"interest"`
	Pnl      string `json:"pnl"`

	MgnMode string `json:"mgnMode"`

	FillTime    string `json:"fillTime"`
	FillIdxPx   string `json:"fillIdxPx"`
	FillMarkPx  string `json:"fillMarkPx"`
	FillMarkVol string `json:"fillMarkVol"`
	FillPxUsd   string `json:"fillPxUsd"`
	FillPxVol   string `json:"fillPxVol"`
	FillFwdPx   string `json:"fillFwdPx"`

	PosBal    string `json:"posBal"`
	PosBalChg string `json:"posBalChg"`

	From string `json:"from"`
	To   string `json:"to"`

	Notes string `json:"notes"`

	TS int64 `json:"ts,string"`
}

type accountBillsQuery struct {
	instType string
	instId   string
	ccy      string
	mgnMode  string
	ctType   string
	billType string
	subType  string

	after  string
	before string
	begin  string
	end    string
	limit  *int
}

func (q accountBillsQuery) values() url.Values {
	v := url.Values{}
	if q.instType != "" {
		v.Set("instType", q.instType)
	}
	if q.instId != "" {
		v.Set("instId", q.instId)
	}
	if q.ccy != "" {
		v.Set("ccy", q.ccy)
	}
	if q.mgnMode != "" {
		v.Set("mgnMode", q.mgnMode)
	}
	if q.ctType != "" {
		v.Set("ctType", q.ctType)
	}
	if q.billType != "" {
		v.Set("type", q.billType)
	}
	if q.subType != "" {
		v.Set("subType", q.subType)
	}

	if q.after != "" {
		v.Set("after", q.after)
	}
	if q.before != "" {
		v.Set("before", q.before)
	}
	if q.begin != "" {
		v.Set("begin", q.begin)
	}
	if q.end != "" {
		v.Set("end", q.end)
	}
	if q.limit != nil {
		v.Set("limit", strconv.Itoa(*q.limit))
	}

	if len(v) == 0 {
		return nil
	}
	return v
}

// AccountBillsService 查询交易账户账单流水（近七天）。
type AccountBillsService struct {
	c *Client
	q accountBillsQuery
}

// NewAccountBillsService 创建 AccountBillsService。
func (c *Client) NewAccountBillsService() *AccountBillsService {
	return &AccountBillsService{c: c}
}

// InstType 设置产品类型（SPOT/MARGIN/SWAP/FUTURES/OPTION）。
func (s *AccountBillsService) InstType(instType string) *AccountBillsService {
	s.q.instType = instType
	return s
}

// InstId 设置产品 ID，如 BTC-USDT。
func (s *AccountBillsService) InstId(instId string) *AccountBillsService {
	s.q.instId = instId
	return s
}

// Ccy 设置币种过滤（支持多币种，逗号分隔；最多 20 个）。
func (s *AccountBillsService) Ccy(ccy string) *AccountBillsService {
	s.q.ccy = ccy
	return s
}

// MgnMode 设置仓位类型（isolated/cross）。
func (s *AccountBillsService) MgnMode(mgnMode string) *AccountBillsService {
	s.q.mgnMode = mgnMode
	return s
}

// CtType 设置合约类型（linear/inverse），仅交割/永续有效。
func (s *AccountBillsService) CtType(ctType string) *AccountBillsService {
	s.q.ctType = ctType
	return s
}

// Type 设置账单类型（如 2=交易，8=资金费 等）。
func (s *AccountBillsService) Type(billType string) *AccountBillsService {
	s.q.billType = billType
	return s
}

// SubType 设置账单子类型。
func (s *AccountBillsService) SubType(subType string) *AccountBillsService {
	s.q.subType = subType
	return s
}

// After 请求此 id 之前（更旧的数据）的分页内容，传 billId。
func (s *AccountBillsService) After(after string) *AccountBillsService {
	s.q.after = after
	return s
}

// Before 请求此 id 之后（更新的数据）的分页内容，传 billId。
func (s *AccountBillsService) Before(before string) *AccountBillsService {
	s.q.before = before
	return s
}

// Begin 筛选开始时间（Unix 毫秒字符串）。
func (s *AccountBillsService) Begin(begin string) *AccountBillsService {
	s.q.begin = begin
	return s
}

// End 筛选结束时间（Unix 毫秒字符串）。
func (s *AccountBillsService) End(end string) *AccountBillsService {
	s.q.end = end
	return s
}

// Limit 设置返回条数（最大 100，默认 100）。
func (s *AccountBillsService) Limit(limit int) *AccountBillsService {
	s.q.limit = &limit
	return s
}

// Do 查询交易账户账单流水（GET /api/v5/account/bills）。
func (s *AccountBillsService) Do(ctx context.Context) ([]AccountBill, error) {
	var data []AccountBill
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/bills", s.q.values(), nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
