package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// AssetDeposit 表示充值记录。
// 数值字段保持为 string（无损）。
type AssetDeposit struct {
	Ccy   string `json:"ccy"`
	Chain string `json:"chain"`

	Amt          string `json:"amt"`
	From         string `json:"from"`
	AreaCodeFrom string `json:"areaCodeFrom"`
	To           string `json:"to"`
	TxId         string `json:"txId"`

	TS    int64  `json:"ts,string"`
	State string `json:"state"`

	DepId               string `json:"depId"`
	FromWdId            string `json:"fromWdId"`
	ActualDepBlkConfirm string `json:"actualDepBlkConfirm"`
}

type assetDepositHistoryQuery struct {
	ccy      string
	depId    string
	fromWdId string
	txId     string
	typ      string
	state    string
	after    string
	before   string
	limit    *int
}

func (q assetDepositHistoryQuery) values() url.Values {
	v := url.Values{}
	if q.ccy != "" {
		v.Set("ccy", q.ccy)
	}
	if q.depId != "" {
		v.Set("depId", q.depId)
	}
	if q.fromWdId != "" {
		v.Set("fromWdId", q.fromWdId)
	}
	if q.txId != "" {
		v.Set("txId", q.txId)
	}
	if q.typ != "" {
		v.Set("type", q.typ)
	}
	if q.state != "" {
		v.Set("state", q.state)
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

// AssetDepositHistoryService 获取充值记录（近 3 个月）。
type AssetDepositHistoryService struct {
	c *Client
	q assetDepositHistoryQuery
}

// NewAssetDepositHistoryService 创建 AssetDepositHistoryService。
func (c *Client) NewAssetDepositHistoryService() *AssetDepositHistoryService {
	return &AssetDepositHistoryService{c: c}
}

// Ccy 设置币种过滤。
func (s *AssetDepositHistoryService) Ccy(ccy string) *AssetDepositHistoryService {
	s.q.ccy = ccy
	return s
}

// DepId 设置充值记录 ID。
func (s *AssetDepositHistoryService) DepId(depId string) *AssetDepositHistoryService {
	s.q.depId = depId
	return s
}

// FromWdId 设置内部转账发起者提币申请 ID。
func (s *AssetDepositHistoryService) FromWdId(fromWdId string) *AssetDepositHistoryService {
	s.q.fromWdId = fromWdId
	return s
}

// TxId 设置区块转账哈希记录。
func (s *AssetDepositHistoryService) TxId(txId string) *AssetDepositHistoryService {
	s.q.txId = txId
	return s
}

// Type 设置充值方式（3=内部转账，4=链上充值）。
func (s *AssetDepositHistoryService) Type(typ string) *AssetDepositHistoryService {
	s.q.typ = typ
	return s
}

// State 设置充值状态。
func (s *AssetDepositHistoryService) State(state string) *AssetDepositHistoryService {
	s.q.state = state
	return s
}

// After 查询在此之前的内容（更旧的数据），值为时间戳（Unix 毫秒字符串）。
func (s *AssetDepositHistoryService) After(after string) *AssetDepositHistoryService {
	s.q.after = after
	return s
}

// Before 查询在此之后的内容（更新的数据），值为时间戳（Unix 毫秒字符串）。
func (s *AssetDepositHistoryService) Before(before string) *AssetDepositHistoryService {
	s.q.before = before
	return s
}

// Limit 设置返回条数（最大 100，默认 100）。
func (s *AssetDepositHistoryService) Limit(limit int) *AssetDepositHistoryService {
	s.q.limit = &limit
	return s
}

// Do 获取充值记录（GET /api/v5/asset/deposit-history）。
func (s *AssetDepositHistoryService) Do(ctx context.Context) ([]AssetDeposit, error) {
	var data []AssetDeposit
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/asset/deposit-history", s.q.values(), nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
