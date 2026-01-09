package okx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// AccountPositionsHistory 表示历史持仓信息（最近 3 个月有更新的仓位）。
//
// 数值与时间字段按 OKX 返回保持为 string/UnixMilli（无损），未包含字段会被忽略，后续可按需补齐。
type AccountPositionsHistory struct {
	InstType string `json:"instType"`
	InstId   string `json:"instId"`
	MgnMode  string `json:"mgnMode"`
	Type     string `json:"type"`

	CTime UnixMilli `json:"cTime"`
	UTime UnixMilli `json:"uTime"`

	PosId string `json:"posId"`

	OpenAvgPx      string `json:"openAvgPx"`
	NonSettleAvgPx string `json:"nonSettleAvgPx"`
	CloseAvgPx     string `json:"closeAvgPx"`
	OpenMaxPos     string `json:"openMaxPos"`
	CloseTotalPos  string `json:"closeTotalPos"`
	RealizedPnl    string `json:"realizedPnl"`
	SettledPnl     string `json:"settledPnl"`
	PnlRatio       string `json:"pnlRatio"`
	Fee            string `json:"fee"`
	FundingFee     string `json:"fundingFee"`
	LiqPenalty     string `json:"liqPenalty"`
	Pnl            string `json:"pnl"`
	PosSide        string `json:"posSide"`
	Lever          string `json:"lever"`
	Direction      string `json:"direction"`
	TriggerPx      string `json:"triggerPx"`
	Uly            string `json:"uly"`
	Ccy            string `json:"ccy"`
}

// AccountPositionsHistoryService 查看历史持仓信息。
type AccountPositionsHistoryService struct {
	c *Client

	instType string
	instId   string
	mgnMode  string
	typ      string
	posId    string

	after  string
	before string
	limit  *int
}

// NewAccountPositionsHistoryService 创建 AccountPositionsHistoryService。
func (c *Client) NewAccountPositionsHistoryService() *AccountPositionsHistoryService {
	return &AccountPositionsHistoryService{c: c}
}

// InstType 设置产品类型（MARGIN/SWAP/FUTURES/OPTION）。
func (s *AccountPositionsHistoryService) InstType(instType string) *AccountPositionsHistoryService {
	s.instType = instType
	return s
}

// InstId 设置交易产品 ID，如 BTC-USD-SWAP。
func (s *AccountPositionsHistoryService) InstId(instId string) *AccountPositionsHistoryService {
	s.instId = instId
	return s
}

// MgnMode 设置保证金模式（cross/isolated）。
func (s *AccountPositionsHistoryService) MgnMode(mgnMode string) *AccountPositionsHistoryService {
	s.mgnMode = mgnMode
	return s
}

// Type 设置最近一次平仓类型（可选，见 OKX 文档）。
func (s *AccountPositionsHistoryService) Type(typ string) *AccountPositionsHistoryService {
	s.typ = typ
	return s
}

// PosId 设置持仓 ID（可选）。
func (s *AccountPositionsHistoryService) PosId(posId string) *AccountPositionsHistoryService {
	s.posId = posId
	return s
}

// After 查询仓位更新 uTime 之前的内容（Unix 毫秒字符串）。
func (s *AccountPositionsHistoryService) After(after string) *AccountPositionsHistoryService {
	s.after = after
	return s
}

// Before 查询仓位更新 uTime 之后的内容（Unix 毫秒字符串）。
func (s *AccountPositionsHistoryService) Before(before string) *AccountPositionsHistoryService {
	s.before = before
	return s
}

// Limit 设置返回条数（最大 100，默认 100；uTime 相同的记录会一次性返回）。
func (s *AccountPositionsHistoryService) Limit(limit int) *AccountPositionsHistoryService {
	s.limit = &limit
	return s
}

// Do 查看历史持仓信息（GET /api/v5/account/positions-history）。
func (s *AccountPositionsHistoryService) Do(ctx context.Context) ([]AccountPositionsHistory, error) {
	q := url.Values{}
	if s.instType != "" {
		q.Set("instType", s.instType)
	}
	if s.instId != "" {
		q.Set("instId", s.instId)
	}
	if s.mgnMode != "" {
		q.Set("mgnMode", s.mgnMode)
	}
	if s.typ != "" {
		q.Set("type", s.typ)
	}
	if s.posId != "" {
		q.Set("posId", s.posId)
	}
	if s.after != "" {
		q.Set("after", s.after)
	}
	if s.before != "" {
		q.Set("before", s.before)
	}
	if s.limit != nil {
		q.Set("limit", strconv.Itoa(*s.limit))
	}
	if len(q) == 0 {
		q = nil
	}

	var data []AccountPositionsHistory
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/positions-history", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
