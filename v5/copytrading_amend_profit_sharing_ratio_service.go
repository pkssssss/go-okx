package okx

import (
	"context"
	"errors"
	"net/http"
)

// CopyTradingAmendProfitSharingRatioService 修改分润比例。
type CopyTradingAmendProfitSharingRatioService struct {
	c *Client

	instType           string
	profitSharingRatio string
}

// NewCopyTradingAmendProfitSharingRatioService 创建 CopyTradingAmendProfitSharingRatioService。
func (c *Client) NewCopyTradingAmendProfitSharingRatioService() *CopyTradingAmendProfitSharingRatioService {
	return &CopyTradingAmendProfitSharingRatioService{c: c}
}

// InstType 设置产品类型（默认 SWAP）。
func (s *CopyTradingAmendProfitSharingRatioService) InstType(instType string) *CopyTradingAmendProfitSharingRatioService {
	s.instType = instType
	return s
}

// ProfitSharingRatio 设置分润比例（必填；0.1 代表 10%）。
func (s *CopyTradingAmendProfitSharingRatioService) ProfitSharingRatio(profitSharingRatio string) *CopyTradingAmendProfitSharingRatioService {
	s.profitSharingRatio = profitSharingRatio
	return s
}

var (
	errCopyTradingAmendProfitSharingRatioMissingProfitSharingRatio = errors.New("okx: copytrading amend profit sharing ratio requires profitSharingRatio")
	errEmptyCopyTradingAmendProfitSharingRatioResponse             = errors.New("okx: empty copytrading amend profit sharing ratio response")
)

type copyTradingAmendProfitSharingRatioRequest struct {
	InstType           string `json:"instType,omitempty"`
	ProfitSharingRatio string `json:"profitSharingRatio"`
}

// Do 修改分润比例（POST /api/v5/copytrading/amend-profit-sharing-ratio）。
func (s *CopyTradingAmendProfitSharingRatioService) Do(ctx context.Context) (*CopyTradingResult, error) {
	if s.profitSharingRatio == "" {
		return nil, errCopyTradingAmendProfitSharingRatioMissingProfitSharingRatio
	}

	req := copyTradingAmendProfitSharingRatioRequest{
		InstType:           s.instType,
		ProfitSharingRatio: s.profitSharingRatio,
	}

	var data []CopyTradingResult
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/copytrading/amend-profit-sharing-ratio", nil, req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyCopyTradingAmendProfitSharingRatioResponse
	}
	return &data[0], nil
}
