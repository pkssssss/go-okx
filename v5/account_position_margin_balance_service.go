package okx

import (
	"context"
	"errors"
	"net/http"
)

type accountPositionMarginBalanceRequest struct {
	InstId  string `json:"instId"`
	PosSide string `json:"posSide"`
	Type    string `json:"type"`
	Amt     string `json:"amt"`
	Ccy     string `json:"ccy,omitempty"`
}

// AccountPositionMarginBalanceAck 表示调整保证金返回项。
type AccountPositionMarginBalanceAck struct {
	InstId   string `json:"instId"`
	PosSide  string `json:"posSide"`
	Type     string `json:"type"`
	Amt      string `json:"amt"`
	Leverage string `json:"leverage"`
	Ccy      string `json:"ccy"`
}

// AccountPositionMarginBalanceService 调整保证金（逐仓）。
type AccountPositionMarginBalanceService struct {
	c   *Client
	req accountPositionMarginBalanceRequest
}

// NewAccountPositionMarginBalanceService 创建 AccountPositionMarginBalanceService。
func (c *Client) NewAccountPositionMarginBalanceService() *AccountPositionMarginBalanceService {
	return &AccountPositionMarginBalanceService{c: c}
}

// InstId 设置产品 ID（必填）。
func (s *AccountPositionMarginBalanceService) InstId(instId string) *AccountPositionMarginBalanceService {
	s.req.InstId = instId
	return s
}

// PosSide 设置持仓方向（必填：long/short/net）。
func (s *AccountPositionMarginBalanceService) PosSide(posSide string) *AccountPositionMarginBalanceService {
	s.req.PosSide = posSide
	return s
}

// Type 设置增加/减少（必填：add/reduce）。
func (s *AccountPositionMarginBalanceService) Type(typ string) *AccountPositionMarginBalanceService {
	s.req.Type = typ
	return s
}

// Amt 设置保证金数量（必填）。
func (s *AccountPositionMarginBalanceService) Amt(amt string) *AccountPositionMarginBalanceService {
	s.req.Amt = amt
	return s
}

// Ccy 设置保证金币种（可选，仅逐仓杠杆仓位适用）。
func (s *AccountPositionMarginBalanceService) Ccy(ccy string) *AccountPositionMarginBalanceService {
	s.req.Ccy = ccy
	return s
}

var (
	errAccountPositionMarginBalanceMissingRequired = errors.New("okx: position margin balance requires instId/posSide/type/amt")
	errEmptyAccountPositionMarginBalance           = errors.New("okx: empty position margin balance response")
)

// Do 调整保证金（POST /api/v5/account/position/margin-balance）。
func (s *AccountPositionMarginBalanceService) Do(ctx context.Context) (*AccountPositionMarginBalanceAck, error) {
	if s.req.InstId == "" || s.req.PosSide == "" || s.req.Type == "" || s.req.Amt == "" {
		return nil, errAccountPositionMarginBalanceMissingRequired
	}

	var data []AccountPositionMarginBalanceAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/account/position/margin-balance", nil, s.req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyAccountPositionMarginBalance
	}
	return &data[0], nil
}
