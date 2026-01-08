package okx

import (
	"context"
	"errors"
	"net/http"
)

// AccountConfig 表示账户配置信息。
//
// 字段按 OKX 返回保持为 string/bool/[]string；未包含字段会被忽略，后续可按需补齐。
type AccountConfig struct {
	Uid     string `json:"uid"`
	MainUid string `json:"mainUid"`

	AcctLv      string `json:"acctLv"`
	AcctStpMode string `json:"acctStpMode"`
	PosMode     string `json:"posMode"`

	AutoLoan   bool   `json:"autoLoan"`
	GreeksType string `json:"greeksType"`
	FeeType    string `json:"feeType"`

	Level    string `json:"level"`
	LevelTmp string `json:"levelTmp"`

	CtIsoMode  string `json:"ctIsoMode"`
	MgnIsoMode string `json:"mgnIsoMode"`

	SpotOffsetType string `json:"spotOffsetType"`
	StgyType       string `json:"stgyType"`

	RoleType        string   `json:"roleType"`
	TraderInsts     []string `json:"traderInsts"`
	SpotRoleType    string   `json:"spotRoleType"`
	SpotTraderInsts []string `json:"spotTraderInsts"`

	OpAuth string `json:"opAuth"`
	KycLv  string `json:"kycLv"`

	Label string `json:"label"`
	Ip    string `json:"ip"`
	Perm  string `json:"perm"`

	LiquidationGear string `json:"liquidationGear"`

	EnableSpotBorrow    bool `json:"enableSpotBorrow"`
	SpotBorrowAutoRepay bool `json:"spotBorrowAutoRepay"`

	Type string `json:"type"`

	SettleCcy     string   `json:"settleCcy"`
	SettleCcyList []string `json:"settleCcyList"`
}

// AccountConfigService 查看账户配置。
type AccountConfigService struct {
	c *Client
}

// NewAccountConfigService 创建 AccountConfigService。
func (c *Client) NewAccountConfigService() *AccountConfigService {
	return &AccountConfigService{c: c}
}

var errEmptyAccountConfig = errors.New("okx: empty account config response")

// Do 查看账户配置（GET /api/v5/account/config）。
func (s *AccountConfigService) Do(ctx context.Context) (*AccountConfig, error) {
	var data []AccountConfig
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/config", nil, nil, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyAccountConfig
	}
	return &data[0], nil
}
