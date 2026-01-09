package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// AccountPrecheckSetDeltaNeutralUnmatchedInfo 表示不匹配信息。
type AccountPrecheckSetDeltaNeutralUnmatchedInfo struct {
	PosList    []string `json:"posList"`
	OrdList    []string `json:"ordList"`
	DeltaLever string   `json:"deltaLever"`
	Type       string   `json:"type"`
}

// AccountPrecheckSetDeltaNeutralResult 表示设置 Delta 中性预检查返回项。
type AccountPrecheckSetDeltaNeutralResult struct {
	UnmatchedInfoCheck []AccountPrecheckSetDeltaNeutralUnmatchedInfo `json:"unmatchedInfoCheck"`
}

// AccountPrecheckSetDeltaNeutralService 设置 Delta 中性预检查。
type AccountPrecheckSetDeltaNeutralService struct {
	c *Client

	stgyType string
}

// NewAccountPrecheckSetDeltaNeutralService 创建 AccountPrecheckSetDeltaNeutralService。
func (c *Client) NewAccountPrecheckSetDeltaNeutralService() *AccountPrecheckSetDeltaNeutralService {
	return &AccountPrecheckSetDeltaNeutralService{c: c}
}

// StgyType 设置策略类型（必填：0 普通策略模式；1 delta 中性策略模式）。
func (s *AccountPrecheckSetDeltaNeutralService) StgyType(stgyType string) *AccountPrecheckSetDeltaNeutralService {
	s.stgyType = stgyType
	return s
}

var (
	errAccountPrecheckSetDeltaNeutralMissingStgyType = errors.New("okx: precheck set delta neutral requires stgyType")
	errEmptyAccountPrecheckSetDeltaNeutral           = errors.New("okx: empty precheck set delta neutral response")
)

// Do 设置 Delta 中性预检查（GET /api/v5/account/precheck-set-delta-neutral）。
func (s *AccountPrecheckSetDeltaNeutralService) Do(ctx context.Context) (*AccountPrecheckSetDeltaNeutralResult, error) {
	if s.stgyType == "" {
		return nil, errAccountPrecheckSetDeltaNeutralMissingStgyType
	}

	q := url.Values{}
	q.Set("stgyType", s.stgyType)

	var data []AccountPrecheckSetDeltaNeutralResult
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/precheck-set-delta-neutral", q, nil, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyAccountPrecheckSetDeltaNeutral
	}
	return &data[0], nil
}
