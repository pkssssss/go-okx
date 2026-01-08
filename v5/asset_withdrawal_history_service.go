package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// AssetWithdrawal 表示提币记录。
// 数值字段保持为 string（无损）。
type AssetWithdrawal struct {
	Ccy              string `json:"ccy"`
	Chain            string `json:"chain"`
	NonTradableAsset bool   `json:"nonTradableAsset"`

	Amt string `json:"amt"`
	TS  int64  `json:"ts,string"`

	From         string `json:"from"`
	AreaCodeFrom string `json:"areaCodeFrom"`
	To           string `json:"to"`
	AreaCodeTo   string `json:"areaCodeTo"`

	ToAddrType string            `json:"toAddrType"`
	Tag        string            `json:"tag"`
	PmtId      string            `json:"pmtId"`
	Memo       string            `json:"memo"`
	AddrEx     map[string]string `json:"addrEx"`

	TxId   string `json:"txId"`
	Fee    string `json:"fee"`
	FeeCcy string `json:"feeCcy"`

	State    string `json:"state"`
	WdId     string `json:"wdId"`
	ClientId string `json:"clientId"`
	Note     string `json:"note"`
}

type assetWithdrawalHistoryQuery struct {
	ccy      string
	wdId     string
	clientId string
	txId     string
	typ      string
	state    string
	after    string
	before   string
	limit    *int
}

func (q assetWithdrawalHistoryQuery) values() url.Values {
	v := url.Values{}
	if q.ccy != "" {
		v.Set("ccy", q.ccy)
	}
	if q.wdId != "" {
		v.Set("wdId", q.wdId)
	}
	if q.clientId != "" {
		v.Set("clientId", q.clientId)
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

// AssetWithdrawalHistoryService 获取提币记录。
type AssetWithdrawalHistoryService struct {
	c *Client
	q assetWithdrawalHistoryQuery
}

// NewAssetWithdrawalHistoryService 创建 AssetWithdrawalHistoryService。
func (c *Client) NewAssetWithdrawalHistoryService() *AssetWithdrawalHistoryService {
	return &AssetWithdrawalHistoryService{c: c}
}

// Ccy 设置币种过滤。
func (s *AssetWithdrawalHistoryService) Ccy(ccy string) *AssetWithdrawalHistoryService {
	s.q.ccy = ccy
	return s
}

// WdId 设置提币申请 ID。
func (s *AssetWithdrawalHistoryService) WdId(wdId string) *AssetWithdrawalHistoryService {
	s.q.wdId = wdId
	return s
}

// ClientId 设置客户自定义 ID（1-32）。
func (s *AssetWithdrawalHistoryService) ClientId(clientId string) *AssetWithdrawalHistoryService {
	s.q.clientId = clientId
	return s
}

// TxId 设置区块转账哈希记录。
func (s *AssetWithdrawalHistoryService) TxId(txId string) *AssetWithdrawalHistoryService {
	s.q.txId = txId
	return s
}

// Type 设置提币方式（3=内部转账，4=链上提币）。
func (s *AssetWithdrawalHistoryService) Type(typ string) *AssetWithdrawalHistoryService {
	s.q.typ = typ
	return s
}

// State 设置提币状态。
func (s *AssetWithdrawalHistoryService) State(state string) *AssetWithdrawalHistoryService {
	s.q.state = state
	return s
}

// After 查询在此之前的内容（更旧的数据），值为时间戳（Unix 毫秒字符串）。
func (s *AssetWithdrawalHistoryService) After(after string) *AssetWithdrawalHistoryService {
	s.q.after = after
	return s
}

// Before 查询在此之后的内容（更新的数据），值为时间戳（Unix 毫秒字符串）。
func (s *AssetWithdrawalHistoryService) Before(before string) *AssetWithdrawalHistoryService {
	s.q.before = before
	return s
}

// Limit 设置返回条数（最大 100，默认 100）。
func (s *AssetWithdrawalHistoryService) Limit(limit int) *AssetWithdrawalHistoryService {
	s.q.limit = &limit
	return s
}

// Do 获取提币记录（GET /api/v5/asset/withdrawal-history）。
func (s *AssetWithdrawalHistoryService) Do(ctx context.Context) ([]AssetWithdrawal, error) {
	var data []AssetWithdrawal
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/asset/withdrawal-history", s.q.values(), nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
