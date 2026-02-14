package okx

import (
	"context"
	"errors"
	"net/http"
)

// CopyTradingCloseSubpositionService 平仓带单。
type CopyTradingCloseSubpositionService struct {
	c *Client

	instType string
	subPosId string
	tag      string
	ordType  string
	px       string
}

// NewCopyTradingCloseSubpositionService 创建 CopyTradingCloseSubpositionService。
func (c *Client) NewCopyTradingCloseSubpositionService() *CopyTradingCloseSubpositionService {
	return &CopyTradingCloseSubpositionService{c: c}
}

// InstType 设置产品类型（默认 SWAP）。
func (s *CopyTradingCloseSubpositionService) InstType(instType string) *CopyTradingCloseSubpositionService {
	s.instType = instType
	return s
}

// SubPosId 设置带单仓位 ID（必填）。
func (s *CopyTradingCloseSubpositionService) SubPosId(subPosId string) *CopyTradingCloseSubpositionService {
	s.subPosId = subPosId
	return s
}

// Tag 设置订单标签。
func (s *CopyTradingCloseSubpositionService) Tag(tag string) *CopyTradingCloseSubpositionService {
	s.tag = tag
	return s
}

// OrdType 设置订单类型（market/limit；默认 market）。
func (s *CopyTradingCloseSubpositionService) OrdType(ordType string) *CopyTradingCloseSubpositionService {
	s.ordType = ordType
	return s
}

// Px 设置委托价格（仅适用于 limit；委托价为 0 代表撤销挂单）。
func (s *CopyTradingCloseSubpositionService) Px(px string) *CopyTradingCloseSubpositionService {
	s.px = px
	return s
}

var (
	errCopyTradingCloseSubpositionMissingSubPosId = errors.New("okx: copytrading close subposition requires subPosId")
	errEmptyCopyTradingCloseSubpositionResponse   = errors.New("okx: empty copytrading close subposition response")
)

type copyTradingCloseSubpositionRequest struct {
	InstType string `json:"instType,omitempty"`
	SubPosId string `json:"subPosId"`
	Tag      string `json:"tag,omitempty"`
	OrdType  string `json:"ordType,omitempty"`
	Px       string `json:"px,omitempty"`
}

// Do 平仓带单（POST /api/v5/copytrading/close-subposition）。
func (s *CopyTradingCloseSubpositionService) Do(ctx context.Context) (*CopyTradingSubPositionAck, error) {
	if s.subPosId == "" {
		return nil, errCopyTradingCloseSubpositionMissingSubPosId
	}

	req := copyTradingCloseSubpositionRequest{
		InstType: s.instType,
		SubPosId: s.subPosId,
		Tag:      s.tag,
		OrdType:  s.ordType,
		Px:       s.px,
	}

	var data []CopyTradingSubPositionAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/copytrading/close-subposition", nil, req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/copytrading/close-subposition", requestID, errEmptyCopyTradingCloseSubpositionResponse)
	}
	return &data[0], nil
}
