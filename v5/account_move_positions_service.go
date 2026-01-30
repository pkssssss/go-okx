package okx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const accountMovePositionsMaxLegs = 30

// AccountMovePositionsLegFrom 表示源账户仓位信息。
type AccountMovePositionsLegFrom struct {
	PosId string `json:"posId"`
	Sz    string `json:"sz"`
	Side  string `json:"side"`
}

// AccountMovePositionsLegTo 表示目标账户移仓配置。
type AccountMovePositionsLegTo struct {
	TdMode  string `json:"tdMode,omitempty"`
	PosSide string `json:"posSide,omitempty"`
	Ccy     string `json:"ccy,omitempty"`
}

// AccountMovePositionsLeg 表示单笔移仓腿（leg）。
type AccountMovePositionsLeg struct {
	From AccountMovePositionsLegFrom `json:"from"`
	To   AccountMovePositionsLegTo   `json:"to"`
}

type accountMovePositionsRequest struct {
	FromAcct string                    `json:"fromAcct"`
	ToAcct   string                    `json:"toAcct"`
	Legs     []AccountMovePositionsLeg `json:"legs"`
	ClientId string                    `json:"clientId"`
}

// AccountMovePositionsLegResultFrom 表示返回中的源仓位信息。
type AccountMovePositionsLegResultFrom struct {
	PosId  string `json:"posId"`
	InstId string `json:"instId"`
	Px     string `json:"px"`
	Side   string `json:"side"`
	Sz     string `json:"sz"`
	SCode  string `json:"sCode"`
	SMsg   string `json:"sMsg"`
}

// AccountMovePositionsLegResultTo 表示返回中的目标仓位信息。
type AccountMovePositionsLegResultTo struct {
	InstId  string `json:"instId"`
	Px      string `json:"px"`
	Side    string `json:"side"`
	Sz      string `json:"sz"`
	TdMode  string `json:"tdMode"`
	PosSide string `json:"posSide"`
	Ccy     string `json:"ccy"`
	SCode   string `json:"sCode"`
	SMsg    string `json:"sMsg"`
}

// AccountMovePositionsLegResult 表示返回中的单笔移仓腿结果。
type AccountMovePositionsLegResult struct {
	From AccountMovePositionsLegResultFrom `json:"from"`
	To   AccountMovePositionsLegResultTo   `json:"to"`
}

// AccountMovePositionsAck 表示移仓返回项。
type AccountMovePositionsAck struct {
	ClientId  string                          `json:"clientId"`
	BlockTdId string                          `json:"blockTdId"`
	State     string                          `json:"state"`
	TS        int64                           `json:"ts,string"`
	FromAcct  string                          `json:"fromAcct"`
	ToAcct    string                          `json:"toAcct"`
	Legs      []AccountMovePositionsLegResult `json:"legs"`
}

// AccountMovePositionsService 移仓（子账户间仓位划转）。
type AccountMovePositionsService struct {
	c   *Client
	req accountMovePositionsRequest
}

// NewAccountMovePositionsService 创建 AccountMovePositionsService。
func (c *Client) NewAccountMovePositionsService() *AccountMovePositionsService {
	return &AccountMovePositionsService{c: c}
}

// FromAcct 设置源账户名（必填，"0" 代表母账户）。
func (s *AccountMovePositionsService) FromAcct(fromAcct string) *AccountMovePositionsService {
	s.req.FromAcct = fromAcct
	return s
}

// ToAcct 设置目标账户名（必填，"0" 代表母账户）。
func (s *AccountMovePositionsService) ToAcct(toAcct string) *AccountMovePositionsService {
	s.req.ToAcct = toAcct
	return s
}

// ClientId 设置客户自定义 ID（必填，1-32）。
func (s *AccountMovePositionsService) ClientId(clientId string) *AccountMovePositionsService {
	s.req.ClientId = clientId
	return s
}

// Legs 设置移仓 legs（必填，最多 30 个）。
func (s *AccountMovePositionsService) Legs(legs []AccountMovePositionsLeg) *AccountMovePositionsService {
	s.req.Legs = legs
	return s
}

var (
	errAccountMovePositionsMissingRequired = errors.New("okx: move positions requires fromAcct/toAcct/clientId and at least one leg")
	errAccountMovePositionsTooManyLegs     = errors.New("okx: move positions max 30 legs")
	errEmptyAccountMovePositions           = errors.New("okx: empty move positions response")
)

// Do 移仓（POST /api/v5/account/move-positions）。
func (s *AccountMovePositionsService) Do(ctx context.Context) (*AccountMovePositionsAck, error) {
	if s.req.FromAcct == "" || s.req.ToAcct == "" || s.req.ClientId == "" || len(s.req.Legs) == 0 {
		return nil, errAccountMovePositionsMissingRequired
	}
	if len(s.req.Legs) > accountMovePositionsMaxLegs {
		return nil, errAccountMovePositionsTooManyLegs
	}
	for i, leg := range s.req.Legs {
		if leg.From.PosId == "" {
			return nil, fmt.Errorf("okx: move positions legs[%d] missing from.posId", i)
		}
		if leg.From.Sz == "" {
			return nil, fmt.Errorf("okx: move positions legs[%d] missing from.sz", i)
		}
		if leg.From.Side == "" {
			return nil, fmt.Errorf("okx: move positions legs[%d] missing from.side", i)
		}
	}

	var data []AccountMovePositionsAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/account/move-positions", nil, s.req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyAccountMovePositions
	}
	ack := &data[0]
	if err := accountCheckMovePositionsAck(http.MethodPost, "/api/v5/account/move-positions", requestID, ack); err != nil {
		return ack, err
	}
	return ack, nil
}
