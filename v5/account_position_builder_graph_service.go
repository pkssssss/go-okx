package okx

import (
	"context"
	"errors"
	"net/http"
)

// AccountPositionBuilderGraphMmrConfig 表示仓位创建器趋势图的 MMR 配置。
type AccountPositionBuilderGraphMmrConfig struct {
	AcctLv string `json:"acctLv,omitempty"`
	Lever  string `json:"lever,omitempty"`
}

type accountPositionBuilderGraphRequest struct {
	InclRealPosAndEq *bool                                 `json:"inclRealPosAndEq,omitempty"`
	SimPos           []AccountPositionBuilderSimPos        `json:"simPos,omitempty"`
	SimAsset         []AccountPositionBuilderSimAsset      `json:"simAsset,omitempty"`
	GreeksType       string                                `json:"greeksType,omitempty"`
	Type             string                                `json:"type"`
	MmrConfig        *AccountPositionBuilderGraphMmrConfig `json:"mmrConfig"`
}

// AccountPositionBuilderGraphMmrData 表示单个 shockFactor 点位的 MMR 数据。
type AccountPositionBuilderGraphMmrData struct {
	ShockFactor string `json:"shockFactor"`
	Mmr         string `json:"mmr"`
	MmrRatio    string `json:"mmrRatio"`
}

// AccountPositionBuilderGraphResult 表示仓位创建器趋势图返回项。
type AccountPositionBuilderGraphResult struct {
	Type    string                               `json:"type"`
	MmrData []AccountPositionBuilderGraphMmrData `json:"mmrData"`
}

// AccountPositionBuilderGraphService 仓位创建器趋势图。
type AccountPositionBuilderGraphService struct {
	c   *Client
	req accountPositionBuilderGraphRequest
}

// NewAccountPositionBuilderGraphService 创建 AccountPositionBuilderGraphService。
func (c *Client) NewAccountPositionBuilderGraphService() *AccountPositionBuilderGraphService {
	return &AccountPositionBuilderGraphService{c: c}
}

// InclRealPosAndEq 是否代入已有仓位和资产（可选，默认 true）。
func (s *AccountPositionBuilderGraphService) InclRealPosAndEq(incl bool) *AccountPositionBuilderGraphService {
	s.req.InclRealPosAndEq = &incl
	return s
}

// SimPos 设置模拟仓位列表（可选，最多 200 个）。
func (s *AccountPositionBuilderGraphService) SimPos(simPos []AccountPositionBuilderSimPos) *AccountPositionBuilderGraphService {
	s.req.SimPos = simPos
	return s
}

// SimAsset 设置模拟资产列表（可选，最多 200 个）。
func (s *AccountPositionBuilderGraphService) SimAsset(simAsset []AccountPositionBuilderSimAsset) *AccountPositionBuilderGraphService {
	s.req.SimAsset = simAsset
	return s
}

// GreeksType 设置希腊值类型（可选：BS/PA/CASH，默认 BS）。
func (s *AccountPositionBuilderGraphService) GreeksType(greeksType string) *AccountPositionBuilderGraphService {
	s.req.GreeksType = greeksType
	return s
}

// Type 设置趋势图类型（必填，目前仅支持 mmr）。
func (s *AccountPositionBuilderGraphService) Type(typ string) *AccountPositionBuilderGraphService {
	s.req.Type = typ
	return s
}

// MmrConfig 设置 MMR 配置（必填）。
func (s *AccountPositionBuilderGraphService) MmrConfig(cfg AccountPositionBuilderGraphMmrConfig) *AccountPositionBuilderGraphService {
	s.req.MmrConfig = &cfg
	return s
}

var errAccountPositionBuilderGraphMissingRequired = errors.New("okx: position builder graph requires type and mmrConfig")

// Do 仓位创建器趋势图（POST /api/v5/account/position-builder-graph）。
func (s *AccountPositionBuilderGraphService) Do(ctx context.Context) ([]AccountPositionBuilderGraphResult, error) {
	if s.req.Type == "" || s.req.MmrConfig == nil {
		return nil, errAccountPositionBuilderGraphMissingRequired
	}
	if err := validateAccountPositionBuilderSimPos(s.req.SimPos); err != nil {
		return nil, err
	}
	if err := validateAccountPositionBuilderSimAsset(s.req.SimAsset); err != nil {
		return nil, err
	}

	var data []AccountPositionBuilderGraphResult
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/account/position-builder-graph", nil, s.req, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
