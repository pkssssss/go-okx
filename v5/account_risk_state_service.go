package okx

import (
	"context"
	"errors"
	"net/http"
)

// AccountRiskState 表示账户特定风险状态（仅适用于 PM 账户）。
type AccountRiskState struct {
	AtRisk    bool      `json:"atRisk"`
	AtRiskIdx []string  `json:"atRiskIdx"`
	AtRiskMgn []string  `json:"atRiskMgn"`
	TS        UnixMilli `json:"ts"`
}

// AccountRiskStateService 查看账户特定风险状态。
type AccountRiskStateService struct {
	c *Client
}

// NewAccountRiskStateService 创建 AccountRiskStateService。
func (c *Client) NewAccountRiskStateService() *AccountRiskStateService {
	return &AccountRiskStateService{c: c}
}

var errEmptyAccountRiskState = errors.New("okx: empty account risk state response")

// Do 查看账户特定风险状态（GET /api/v5/account/risk-state）。
func (s *AccountRiskStateService) Do(ctx context.Context) (*AccountRiskState, error) {
	var data []AccountRiskState
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/risk-state", nil, nil, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyAccountRiskState
	}
	return &data[0], nil
}
