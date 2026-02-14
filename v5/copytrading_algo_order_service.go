package okx

import (
	"context"
	"errors"
	"net/http"
)

// CopyTradingAlgoOrderService 带单或跟单仓位止盈止损。
type CopyTradingAlgoOrderService struct {
	c *Client

	instType        string
	subPosId        string
	tpTriggerPx     string
	slTriggerPx     string
	tpOrdPx         string
	slOrdPx         string
	tpTriggerPxType string
	slTriggerPxType string
	tag             string
	subPosType      string
}

// NewCopyTradingAlgoOrderService 创建 CopyTradingAlgoOrderService。
func (c *Client) NewCopyTradingAlgoOrderService() *CopyTradingAlgoOrderService {
	return &CopyTradingAlgoOrderService{c: c}
}

// InstType 设置产品类型（默认 SWAP）。
func (s *CopyTradingAlgoOrderService) InstType(instType string) *CopyTradingAlgoOrderService {
	s.instType = instType
	return s
}

// SubPosId 设置带单或者跟单仓位 ID（必填）。
func (s *CopyTradingAlgoOrderService) SubPosId(subPosId string) *CopyTradingAlgoOrderService {
	s.subPosId = subPosId
	return s
}

// TpTriggerPx 设置止盈触发价（tpTriggerPx 和 slTriggerPx 至少填写一个；为 0 表示删除止盈）。
func (s *CopyTradingAlgoOrderService) TpTriggerPx(tpTriggerPx string) *CopyTradingAlgoOrderService {
	s.tpTriggerPx = tpTriggerPx
	return s
}

// SlTriggerPx 设置止损触发价（为 0 表示删除止损）。
func (s *CopyTradingAlgoOrderService) SlTriggerPx(slTriggerPx string) *CopyTradingAlgoOrderService {
	s.slTriggerPx = slTriggerPx
	return s
}

// TpOrdPx 设置止盈委托价（-1 为市价；仅适用于现货交易员）。
func (s *CopyTradingAlgoOrderService) TpOrdPx(tpOrdPx string) *CopyTradingAlgoOrderService {
	s.tpOrdPx = tpOrdPx
	return s
}

// SlOrdPx 设置止损委托价（-1 为市价；仅适用于现货交易员）。
func (s *CopyTradingAlgoOrderService) SlOrdPx(slOrdPx string) *CopyTradingAlgoOrderService {
	s.slOrdPx = slOrdPx
	return s
}

// TpTriggerPxType 设置止盈触发价类型（last/index/mark）。
func (s *CopyTradingAlgoOrderService) TpTriggerPxType(tpTriggerPxType string) *CopyTradingAlgoOrderService {
	s.tpTriggerPxType = tpTriggerPxType
	return s
}

// SlTriggerPxType 设置止损触发价类型（last/index/mark）。
func (s *CopyTradingAlgoOrderService) SlTriggerPxType(slTriggerPxType string) *CopyTradingAlgoOrderService {
	s.slTriggerPxType = slTriggerPxType
	return s
}

// Tag 设置订单标签。
func (s *CopyTradingAlgoOrderService) Tag(tag string) *CopyTradingAlgoOrderService {
	s.tag = tag
	return s
}

// SubPosType 设置数据类型（lead/copy；默认 lead）。
func (s *CopyTradingAlgoOrderService) SubPosType(subPosType string) *CopyTradingAlgoOrderService {
	s.subPosType = subPosType
	return s
}

var (
	errCopyTradingAlgoOrderMissingSubPosId  = errors.New("okx: copytrading algo order requires subPosId")
	errCopyTradingAlgoOrderMissingTriggerPx = errors.New("okx: copytrading algo order requires tpTriggerPx or slTriggerPx")
	errEmptyCopyTradingAlgoOrderResponse    = errors.New("okx: empty copytrading algo order response")
)

type copyTradingAlgoOrderRequest struct {
	InstType        string `json:"instType,omitempty"`
	SubPosId        string `json:"subPosId"`
	TpTriggerPx     string `json:"tpTriggerPx,omitempty"`
	SlTriggerPx     string `json:"slTriggerPx,omitempty"`
	TpOrdPx         string `json:"tpOrdPx,omitempty"`
	SlOrdPx         string `json:"slOrdPx,omitempty"`
	TpTriggerPxType string `json:"tpTriggerPxType,omitempty"`
	SlTriggerPxType string `json:"slTriggerPxType,omitempty"`
	Tag             string `json:"tag,omitempty"`
	SubPosType      string `json:"subPosType,omitempty"`
}

// Do 设置止盈止损（POST /api/v5/copytrading/algo-order）。
func (s *CopyTradingAlgoOrderService) Do(ctx context.Context) (*CopyTradingSubPositionAck, error) {
	if s.subPosId == "" {
		return nil, errCopyTradingAlgoOrderMissingSubPosId
	}
	if s.tpTriggerPx == "" && s.slTriggerPx == "" {
		return nil, errCopyTradingAlgoOrderMissingTriggerPx
	}

	req := copyTradingAlgoOrderRequest{
		InstType:        s.instType,
		SubPosId:        s.subPosId,
		TpTriggerPx:     s.tpTriggerPx,
		SlTriggerPx:     s.slTriggerPx,
		TpOrdPx:         s.tpOrdPx,
		SlOrdPx:         s.slOrdPx,
		TpTriggerPxType: s.tpTriggerPxType,
		SlTriggerPxType: s.slTriggerPxType,
		Tag:             s.tag,
		SubPosType:      s.subPosType,
	}

	var data []CopyTradingSubPositionAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/copytrading/algo-order", nil, req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/copytrading/algo-order", requestID, errEmptyCopyTradingAlgoOrderResponse)
	}
	return &data[0], nil
}
