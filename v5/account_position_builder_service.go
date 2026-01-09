package okx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	accountPositionBuilderMaxSimPos   = 200
	accountPositionBuilderMaxSimAsset = 200
)

// AccountPositionBuilderSimPos 表示模拟仓位信息。
type AccountPositionBuilderSimPos struct {
	InstId string `json:"instId"`
	Pos    string `json:"pos"`
	AvgPx  string `json:"avgPx"`
	Lever  string `json:"lever,omitempty"`
}

// AccountPositionBuilderSimAsset 表示模拟资产信息。
type AccountPositionBuilderSimAsset struct {
	Ccy string `json:"ccy"`
	Amt string `json:"amt"`
}

type accountPositionBuilderRequest struct {
	AcctLv           string                           `json:"acctLv,omitempty"`
	InclRealPosAndEq *bool                            `json:"inclRealPosAndEq,omitempty"`
	Lever            string                           `json:"lever,omitempty"`
	SimPos           []AccountPositionBuilderSimPos   `json:"simPos,omitempty"`
	SimAsset         []AccountPositionBuilderSimAsset `json:"simAsset,omitempty"`
	GreeksType       string                           `json:"greeksType,omitempty"`
	IdxVol           string                           `json:"idxVol,omitempty"`
}

// AccountPositionBuilderResult 表示仓位创建器返回项（只强类型化常用字段，其余大字段用 RawMessage 承载）。
type AccountPositionBuilderResult struct {
	Eq          string `json:"eq"`
	TotalImr    string `json:"totalImr"`
	TotalMmr    string `json:"totalMmr"`
	BorrowMmr   string `json:"borrowMmr"`
	DerivMmr    string `json:"derivMmr"`
	MarginRatio string `json:"marginRatio"`
	Upl         string `json:"upl"`
	AcctLever   string `json:"acctLever"`
	TS          int64  `json:"ts,string"`

	Assets       json.RawMessage `json:"assets"`
	Positions    json.RawMessage `json:"positions"`
	RiskUnitData json.RawMessage `json:"riskUnitData"`
}

// AccountPositionBuilderService 仓位创建器（计算模拟/真实仓位的投资组合保证金信息）。
type AccountPositionBuilderService struct {
	c   *Client
	req accountPositionBuilderRequest
}

// NewAccountPositionBuilderService 创建 AccountPositionBuilderService。
func (c *Client) NewAccountPositionBuilderService() *AccountPositionBuilderService {
	return &AccountPositionBuilderService{c: c}
}

// AcctLv 切换至账户模式（可选：3 跨币种保证金；4 组合保证金）。
func (s *AccountPositionBuilderService) AcctLv(acctLv string) *AccountPositionBuilderService {
	s.req.AcctLv = acctLv
	return s
}

// InclRealPosAndEq 是否代入已有仓位和资产（可选，默认 true）。
func (s *AccountPositionBuilderService) InclRealPosAndEq(incl bool) *AccountPositionBuilderService {
	s.req.InclRealPosAndEq = &incl
	return s
}

// Lever 跨币种下整体的全仓合约杠杆数量（可选，默认 1；仅适用于跨币种保证金模式）。
func (s *AccountPositionBuilderService) Lever(lever string) *AccountPositionBuilderService {
	s.req.Lever = lever
	return s
}

// SimPos 设置模拟仓位列表（可选，最多 200 个）。
func (s *AccountPositionBuilderService) SimPos(simPos []AccountPositionBuilderSimPos) *AccountPositionBuilderService {
	s.req.SimPos = simPos
	return s
}

// SimAsset 设置模拟资产列表（可选，最多 200 个）。
func (s *AccountPositionBuilderService) SimAsset(simAsset []AccountPositionBuilderSimAsset) *AccountPositionBuilderService {
	s.req.SimAsset = simAsset
	return s
}

// GreeksType 设置希腊值类型（可选：BS/PA/CASH，默认 BS）。
func (s *AccountPositionBuilderService) GreeksType(greeksType string) *AccountPositionBuilderService {
	s.req.GreeksType = greeksType
	return s
}

// IdxVol 设置价格变动百分比（可选，小数形式，范围 -0.99~1，0.01 递增，默认 0）。
func (s *AccountPositionBuilderService) IdxVol(idxVol string) *AccountPositionBuilderService {
	s.req.IdxVol = idxVol
	return s
}

var (
	errAccountPositionBuilderTooManySimPos   = errors.New("okx: position builder simPos max 200")
	errAccountPositionBuilderTooManySimAsset = errors.New("okx: position builder simAsset max 200")
	errEmptyAccountPositionBuilder           = errors.New("okx: empty position builder response")
)

func validateAccountPositionBuilderSimPos(simPos []AccountPositionBuilderSimPos) error {
	if len(simPos) > accountPositionBuilderMaxSimPos {
		return errAccountPositionBuilderTooManySimPos
	}
	for i, p := range simPos {
		if p.InstId == "" {
			return fmt.Errorf("okx: position builder simPos[%d] missing instId", i)
		}
		if p.Pos == "" {
			return fmt.Errorf("okx: position builder simPos[%d] missing pos", i)
		}
		if p.AvgPx == "" {
			return fmt.Errorf("okx: position builder simPos[%d] missing avgPx", i)
		}
	}
	return nil
}

func validateAccountPositionBuilderSimAsset(simAsset []AccountPositionBuilderSimAsset) error {
	if len(simAsset) > accountPositionBuilderMaxSimAsset {
		return errAccountPositionBuilderTooManySimAsset
	}
	for i, a := range simAsset {
		if a.Ccy == "" {
			return fmt.Errorf("okx: position builder simAsset[%d] missing ccy", i)
		}
		if a.Amt == "" {
			return fmt.Errorf("okx: position builder simAsset[%d] missing amt", i)
		}
	}
	return nil
}

// Do 仓位创建器（POST /api/v5/account/position-builder）。
func (s *AccountPositionBuilderService) Do(ctx context.Context) (*AccountPositionBuilderResult, error) {
	if err := validateAccountPositionBuilderSimPos(s.req.SimPos); err != nil {
		return nil, err
	}
	if err := validateAccountPositionBuilderSimAsset(s.req.SimAsset); err != nil {
		return nil, err
	}

	var data []AccountPositionBuilderResult
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/account/position-builder", nil, s.req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyAccountPositionBuilder
	}
	return &data[0], nil
}
