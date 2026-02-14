package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// AccountSwitchPrecheckUnmatchedInfo 表示账户模式切换不兼容信息。
type AccountSwitchPrecheckUnmatchedInfo struct {
	Type       string   `json:"type"`
	TotalAsset string   `json:"totalAsset"`
	PosList    []string `json:"posList"`
}

// AccountSwitchPrecheckPos 表示合约全仓仓位信息（切换后杠杆）。
type AccountSwitchPrecheckPos struct {
	PosId string `json:"posId"`
	Lever string `json:"lever"`
}

// AccountSwitchPrecheckPosTierCheck 表示梯度档位校验不通过的仓位信息。
type AccountSwitchPrecheckPosTierCheck struct {
	InstFamily string `json:"instFamily"`
	InstType   string `json:"instType"`
	Pos        string `json:"pos"`
	Lever      string `json:"lever"`
	MaxSz      string `json:"maxSz"`
}

// AccountSwitchPrecheckMarginDetail 表示币种维度保证金信息。
type AccountSwitchPrecheckMarginDetail struct {
	Ccy      string `json:"ccy"`
	AvailEq  string `json:"availEq"`
	MgnRatio string `json:"mgnRatio"`
}

// AccountSwitchPrecheckMargin 表示保证金相关信息。
type AccountSwitchPrecheckMargin struct {
	AcctAvailEq string                              `json:"acctAvailEq"`
	MgnRatio    string                              `json:"mgnRatio"`
	Details     []AccountSwitchPrecheckMarginDetail `json:"details"`
}

// AccountSwitchPrecheckResult 表示账户模式切换预检查结果。
type AccountSwitchPrecheckResult struct {
	SCode          string `json:"sCode"`
	CurAcctLv      string `json:"curAcctLv"`
	AcctLv         string `json:"acctLv"`
	RiskOffsetType string `json:"riskOffsetType"`

	UnmatchedInfoCheck []AccountSwitchPrecheckUnmatchedInfo `json:"unmatchedInfoCheck"`
	PosList            []AccountSwitchPrecheckPos           `json:"posList"`
	PosTierCheck       []AccountSwitchPrecheckPosTierCheck  `json:"posTierCheck"`

	MgnBf  *AccountSwitchPrecheckMargin `json:"mgnBf"`
	MgnAft *AccountSwitchPrecheckMargin `json:"mgnAft"`
}

// AccountSwitchPrecheckService 预检查账户模式切换。
type AccountSwitchPrecheckService struct {
	c *Client

	acctLv string
}

// NewAccountSwitchPrecheckService 创建 AccountSwitchPrecheckService。
func (c *Client) NewAccountSwitchPrecheckService() *AccountSwitchPrecheckService {
	return &AccountSwitchPrecheckService{c: c}
}

// AcctLv 设置目标账户模式（必填）。
func (s *AccountSwitchPrecheckService) AcctLv(acctLv string) *AccountSwitchPrecheckService {
	s.acctLv = acctLv
	return s
}

var (
	errAccountSwitchPrecheckMissingAcctLv = errors.New("okx: account switch precheck requires acctLv")
	errEmptyAccountSwitchPrecheck         = errors.New("okx: empty account switch precheck response")
)

// Do 预检查账户模式切换（GET /api/v5/account/set-account-switch-precheck）。
func (s *AccountSwitchPrecheckService) Do(ctx context.Context) (*AccountSwitchPrecheckResult, error) {
	if s.acctLv == "" {
		return nil, errAccountSwitchPrecheckMissingAcctLv
	}

	q := url.Values{}
	q.Set("acctLv", s.acctLv)

	var data []AccountSwitchPrecheckResult
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodGet, "/api/v5/account/set-account-switch-precheck", q, nil, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyAccountSwitchPrecheck
	}
	if data[0].SCode != "0" {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodGet,
			RequestPath: "/api/v5/account/set-account-switch-precheck",
			RequestID:   requestID,
			Code:        data[0].SCode,
			Message:     "account switch precheck failed",
		}
	}
	return &data[0], nil
}
