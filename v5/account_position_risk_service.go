package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// AccountPositionRiskBalanceData 表示币种资产信息。
type AccountPositionRiskBalanceData struct {
	Ccy   string `json:"ccy"`
	Eq    string `json:"eq"`
	DisEq string `json:"disEq"`
}

// AccountPositionRiskPosData 表示持仓详细信息。
type AccountPositionRiskPosData struct {
	InstType string `json:"instType"`
	MgnMode  string `json:"mgnMode"`
	PosId    string `json:"posId"`
	InstId   string `json:"instId"`
	Pos      string `json:"pos"`

	BaseBal  string `json:"baseBal"`
	QuoteBal string `json:"quoteBal"`

	PosSide string `json:"posSide"`
	PosCcy  string `json:"posCcy"`
	Ccy     string `json:"ccy"`

	NotionalCcy string `json:"notionalCcy"`
	NotionalUsd string `json:"notionalUsd"`
}

// AccountPositionRisk 表示账户持仓风险（同一时间切片下的账户与持仓基础信息）。
type AccountPositionRisk struct {
	TS    UnixMilli `json:"ts"`
	AdjEq string    `json:"adjEq"`

	BalData []AccountPositionRiskBalanceData `json:"balData"`
	PosData []AccountPositionRiskPosData     `json:"posData"`
}

// AccountPositionRiskService 查看账户持仓风险。
type AccountPositionRiskService struct {
	c        *Client
	instType string
}

// NewAccountPositionRiskService 创建 AccountPositionRiskService。
func (c *Client) NewAccountPositionRiskService() *AccountPositionRiskService {
	return &AccountPositionRiskService{c: c}
}

// InstType 设置产品类型（MARGIN/SWAP/FUTURES/OPTION）。
func (s *AccountPositionRiskService) InstType(instType string) *AccountPositionRiskService {
	s.instType = instType
	return s
}

var errEmptyAccountPositionRisk = errors.New("okx: empty account position risk response")

// Do 查看账户持仓风险（GET /api/v5/account/account-position-risk）。
func (s *AccountPositionRiskService) Do(ctx context.Context) (*AccountPositionRisk, error) {
	var q url.Values
	if s.instType != "" {
		q = url.Values{}
		q.Set("instType", s.instType)
	}

	var data []AccountPositionRisk
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/account-position-risk", q, nil, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyAccountPositionRisk
	}
	return &data[0], nil
}
