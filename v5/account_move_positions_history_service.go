package okx

import (
	"context"
	"net/http"
	"net/url"
)

// AccountMovePositionsHistoryLegFrom 表示移仓历史中的源仓位信息。
type AccountMovePositionsHistoryLegFrom struct {
	PosId  string `json:"posId"`
	InstId string `json:"instId"`
	Px     string `json:"px"`
	Side   string `json:"side"`
	Sz     string `json:"sz"`
}

// AccountMovePositionsHistoryLegTo 表示移仓历史中的目标仓位信息。
type AccountMovePositionsHistoryLegTo struct {
	InstId  string `json:"instId"`
	Px      string `json:"px"`
	Side    string `json:"side"`
	Sz      string `json:"sz"`
	TdMode  string `json:"tdMode"`
	PosSide string `json:"posSide"`
	Ccy     string `json:"ccy"`
}

// AccountMovePositionsHistoryLeg 表示移仓历史中的单笔移仓腿。
type AccountMovePositionsHistoryLeg struct {
	From AccountMovePositionsHistoryLegFrom `json:"from"`
	To   AccountMovePositionsHistoryLegTo   `json:"to"`
}

// AccountMovePositionsHistoryItem 表示移仓历史记录。
type AccountMovePositionsHistoryItem struct {
	ClientId  string                           `json:"clientId"`
	BlockTdId string                           `json:"blockTdId"`
	State     string                           `json:"state"`
	TS        int64                            `json:"ts,string"`
	FromAcct  string                           `json:"fromAcct"`
	ToAcct    string                           `json:"toAcct"`
	Legs      []AccountMovePositionsHistoryLeg `json:"legs"`
}

// AccountMovePositionsHistoryService 获取移仓历史（过去 3 天）。
type AccountMovePositionsHistoryService struct {
	c *Client

	blockTdId string
	clientId  string
	beginTs   string
	endTs     string
	limit     string
	state     string
}

// NewAccountMovePositionsHistoryService 创建 AccountMovePositionsHistoryService。
func (c *Client) NewAccountMovePositionsHistoryService() *AccountMovePositionsHistoryService {
	return &AccountMovePositionsHistoryService{c: c}
}

// BlockTdId 设置大宗交易 ID。
func (s *AccountMovePositionsHistoryService) BlockTdId(blockTdId string) *AccountMovePositionsHistoryService {
	s.blockTdId = blockTdId
	return s
}

// ClientId 设置客户自定义 ID。
func (s *AccountMovePositionsHistoryService) ClientId(clientId string) *AccountMovePositionsHistoryService {
	s.clientId = clientId
	return s
}

// BeginTs 设置开始时间戳（毫秒）。
func (s *AccountMovePositionsHistoryService) BeginTs(beginTs string) *AccountMovePositionsHistoryService {
	s.beginTs = beginTs
	return s
}

// EndTs 设置结束时间戳（毫秒）。
func (s *AccountMovePositionsHistoryService) EndTs(endTs string) *AccountMovePositionsHistoryService {
	s.endTs = endTs
	return s
}

// Limit 设置返回数量（最大 100，默认 100）。
func (s *AccountMovePositionsHistoryService) Limit(limit string) *AccountMovePositionsHistoryService {
	s.limit = limit
	return s
}

// State 设置移仓状态过滤（filled/pending）。
func (s *AccountMovePositionsHistoryService) State(state string) *AccountMovePositionsHistoryService {
	s.state = state
	return s
}

// Do 获取移仓历史（GET /api/v5/account/move-positions-history）。
func (s *AccountMovePositionsHistoryService) Do(ctx context.Context) ([]AccountMovePositionsHistoryItem, error) {
	var q url.Values
	if s.blockTdId != "" || s.clientId != "" || s.beginTs != "" || s.endTs != "" || s.limit != "" || s.state != "" {
		q = url.Values{}
		if s.blockTdId != "" {
			q.Set("blockTdId", s.blockTdId)
		}
		if s.clientId != "" {
			q.Set("clientId", s.clientId)
		}
		if s.beginTs != "" {
			q.Set("beginTs", s.beginTs)
		}
		if s.endTs != "" {
			q.Set("endTs", s.endTs)
		}
		if s.limit != "" {
			q.Set("limit", s.limit)
		}
		if s.state != "" {
			q.Set("state", s.state)
		}
	}

	var data []AccountMovePositionsHistoryItem
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/move-positions-history", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
