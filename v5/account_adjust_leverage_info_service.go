package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// AccountAdjustLeverageInfo 表示杠杆倍数预估信息。
type AccountAdjustLeverageInfo struct {
	EstAvailQuoteTrans string `json:"estAvailQuoteTrans"`
	EstAvailTrans      string `json:"estAvailTrans"`
	EstLiqPx           string `json:"estLiqPx"`
	EstMaxAmt          string `json:"estMaxAmt"`
	EstMgn             string `json:"estMgn"`
	EstQuoteMaxAmt     string `json:"estQuoteMaxAmt"`
	EstQuoteMgn        string `json:"estQuoteMgn"`

	ExistOrd bool `json:"existOrd"`

	MaxLever string `json:"maxLever"`
	MinLever string `json:"minLever"`
}

// AccountAdjustLeverageInfoService 获取杠杆倍数预估信息。
type AccountAdjustLeverageInfoService struct {
	c *Client

	instType string
	mgnMode  string
	lever    string
	instId   string
	ccy      string
	posSide  string
}

// NewAccountAdjustLeverageInfoService 创建 AccountAdjustLeverageInfoService。
func (c *Client) NewAccountAdjustLeverageInfoService() *AccountAdjustLeverageInfoService {
	return &AccountAdjustLeverageInfoService{c: c}
}

// InstType 设置产品类型（必填：MARGIN/SWAP/FUTURES）。
func (s *AccountAdjustLeverageInfoService) InstType(instType string) *AccountAdjustLeverageInfoService {
	s.instType = instType
	return s
}

// MgnMode 设置保证金模式（必填：isolated/cross）。
func (s *AccountAdjustLeverageInfoService) MgnMode(mgnMode string) *AccountAdjustLeverageInfoService {
	s.mgnMode = mgnMode
	return s
}

// Lever 设置杠杆倍数（必填）。
func (s *AccountAdjustLeverageInfoService) Lever(lever string) *AccountAdjustLeverageInfoService {
	s.lever = lever
	return s
}

// InstId 设置产品 ID（可选）。
func (s *AccountAdjustLeverageInfoService) InstId(instId string) *AccountAdjustLeverageInfoService {
	s.instId = instId
	return s
}

// Ccy 设置保证金币种（可选）。
func (s *AccountAdjustLeverageInfoService) Ccy(ccy string) *AccountAdjustLeverageInfoService {
	s.ccy = ccy
	return s
}

// PosSide 设置持仓方向（可选：net/long/short）。
func (s *AccountAdjustLeverageInfoService) PosSide(posSide string) *AccountAdjustLeverageInfoService {
	s.posSide = posSide
	return s
}

var (
	errAccountAdjustLeverageInfoMissingRequired = errors.New("okx: adjust leverage info requires instType/mgnMode/lever")
	errEmptyAccountAdjustLeverageInfo           = errors.New("okx: empty adjust leverage info response")
)

// Do 获取杠杆倍数预估信息（GET /api/v5/account/adjust-leverage-info）。
func (s *AccountAdjustLeverageInfoService) Do(ctx context.Context) (*AccountAdjustLeverageInfo, error) {
	if s.instType == "" || s.mgnMode == "" || s.lever == "" {
		return nil, errAccountAdjustLeverageInfoMissingRequired
	}

	q := url.Values{}
	q.Set("instType", s.instType)
	q.Set("mgnMode", s.mgnMode)
	q.Set("lever", s.lever)
	if s.instId != "" {
		q.Set("instId", s.instId)
	}
	if s.ccy != "" {
		q.Set("ccy", s.ccy)
	}
	if s.posSide != "" {
		q.Set("posSide", s.posSide)
	}

	var data []AccountAdjustLeverageInfo
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/adjust-leverage-info", q, nil, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyAccountAdjustLeverageInfo
	}
	return &data[0], nil
}
