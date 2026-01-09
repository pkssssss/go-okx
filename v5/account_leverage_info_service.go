package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// AccountLeverageInfo 表示杠杆倍数信息。
type AccountLeverageInfo struct {
	InstId  string `json:"instId"`
	Ccy     string `json:"ccy"`
	MgnMode string `json:"mgnMode"`
	PosSide string `json:"posSide"`
	Lever   string `json:"lever"`
}

// AccountLeverageInfoService 获取杠杆倍数。
type AccountLeverageInfoService struct {
	c       *Client
	instId  string
	ccy     string
	mgnMode string
}

// NewAccountLeverageInfoService 创建 AccountLeverageInfoService。
func (c *Client) NewAccountLeverageInfoService() *AccountLeverageInfoService {
	return &AccountLeverageInfoService{c: c}
}

// InstId 设置产品 ID（支持多个 instId，逗号分隔；最多 20 个）。
func (s *AccountLeverageInfoService) InstId(instId string) *AccountLeverageInfoService {
	s.instId = instId
	return s
}

// Ccy 设置币种（用于币种维度杠杆，仅全仓币币杠杆适用；支持多个 ccy，逗号分隔；最多 20 个）。
func (s *AccountLeverageInfoService) Ccy(ccy string) *AccountLeverageInfoService {
	s.ccy = ccy
	return s
}

// MgnMode 设置保证金模式（必填：isolated/cross）。
func (s *AccountLeverageInfoService) MgnMode(mgnMode string) *AccountLeverageInfoService {
	s.mgnMode = mgnMode
	return s
}

var errAccountLeverageInfoMissingMgnMode = errors.New("okx: leverage info requires mgnMode")

// Do 获取杠杆倍数（GET /api/v5/account/leverage-info）。
func (s *AccountLeverageInfoService) Do(ctx context.Context) ([]AccountLeverageInfo, error) {
	if s.mgnMode == "" {
		return nil, errAccountLeverageInfoMissingMgnMode
	}

	q := url.Values{}
	q.Set("mgnMode", s.mgnMode)
	if s.instId != "" {
		q.Set("instId", s.instId)
	}
	if s.ccy != "" {
		q.Set("ccy", s.ccy)
	}

	var data []AccountLeverageInfo
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/leverage-info", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
