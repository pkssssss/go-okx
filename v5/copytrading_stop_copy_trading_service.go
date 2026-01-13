package okx

import (
	"context"
	"errors"
	"net/http"
)

// CopyTradingStopCopyTradingService 停止跟单。
type CopyTradingStopCopyTradingService struct {
	c *Client

	instType        string
	uniqueCode      string
	subPosCloseType string
}

// NewCopyTradingStopCopyTradingService 创建 CopyTradingStopCopyTradingService。
func (c *Client) NewCopyTradingStopCopyTradingService() *CopyTradingStopCopyTradingService {
	return &CopyTradingStopCopyTradingService{c: c}
}

// InstType 设置产品类型（默认 SWAP）。
func (s *CopyTradingStopCopyTradingService) InstType(instType string) *CopyTradingStopCopyTradingService {
	s.instType = instType
	return s
}

// UniqueCode 设置交易员唯一标识码（必填）。
func (s *CopyTradingStopCopyTradingService) UniqueCode(uniqueCode string) *CopyTradingStopCopyTradingService {
	s.uniqueCode = uniqueCode
	return s
}

// SubPosCloseType 设置剩余仓位处理方式（market_close/copy_close/manual_close；有相关条目时必填）。
func (s *CopyTradingStopCopyTradingService) SubPosCloseType(subPosCloseType string) *CopyTradingStopCopyTradingService {
	s.subPosCloseType = subPosCloseType
	return s
}

var (
	errCopyTradingStopCopyTradingMissingUniqueCode = errors.New("okx: copytrading stop copy trading requires uniqueCode")
	errEmptyCopyTradingStopCopyTradingResponse     = errors.New("okx: empty copytrading stop copy trading response")
)

type copyTradingStopCopyTradingRequest struct {
	InstType        string `json:"instType,omitempty"`
	UniqueCode      string `json:"uniqueCode"`
	SubPosCloseType string `json:"subPosCloseType,omitempty"`
}

// Do 停止跟单（POST /api/v5/copytrading/stop-copy-trading）。
func (s *CopyTradingStopCopyTradingService) Do(ctx context.Context) (*CopyTradingResult, error) {
	if s.uniqueCode == "" {
		return nil, errCopyTradingStopCopyTradingMissingUniqueCode
	}

	req := copyTradingStopCopyTradingRequest{
		InstType:        s.instType,
		UniqueCode:      s.uniqueCode,
		SubPosCloseType: s.subPosCloseType,
	}

	var data []CopyTradingResult
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/copytrading/stop-copy-trading", nil, req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyCopyTradingStopCopyTradingResponse
	}
	return &data[0], nil
}
