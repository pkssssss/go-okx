package okx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type amendAlgoOrderRequest struct {
	InstId string `json:"instId"`

	AlgoId      string `json:"algoId,omitempty"`
	AlgoClOrdId string `json:"algoClOrdId,omitempty"`

	CxlOnFail *bool  `json:"cxlOnFail,omitempty"`
	ReqId     string `json:"reqId,omitempty"`
	NewSz     string `json:"newSz,omitempty"`

	// 止盈止损
	NewTpTriggerPx     string `json:"newTpTriggerPx,omitempty"`
	NewTpOrdPx         string `json:"newTpOrdPx,omitempty"`
	NewSlTriggerPx     string `json:"newSlTriggerPx,omitempty"`
	NewSlOrdPx         string `json:"newSlOrdPx,omitempty"`
	NewTpTriggerPxType string `json:"newTpTriggerPxType,omitempty"`
	NewSlTriggerPxType string `json:"newSlTriggerPxType,omitempty"`

	// 计划委托
	NewTriggerPx     string                      `json:"newTriggerPx,omitempty"`
	NewOrdPx         string                      `json:"newOrdPx,omitempty"`
	NewTriggerPxType string                      `json:"newTriggerPxType,omitempty"`
	AttachAlgoOrds   []TradeAlgoOrderAttachAmend `json:"attachAlgoOrds,omitempty"`
}

// AmendAlgoOrderService 修改策略委托订单（仅部分类型支持）。
type AmendAlgoOrderService struct {
	c   *Client
	req amendAlgoOrderRequest
}

// NewAmendAlgoOrderService 创建 AmendAlgoOrderService。
func (c *Client) NewAmendAlgoOrderService() *AmendAlgoOrderService {
	return &AmendAlgoOrderService{c: c}
}

func (s *AmendAlgoOrderService) InstId(instId string) *AmendAlgoOrderService {
	s.req.InstId = instId
	return s
}

func (s *AmendAlgoOrderService) AlgoId(algoId string) *AmendAlgoOrderService {
	s.req.AlgoId = algoId
	return s
}

func (s *AmendAlgoOrderService) AlgoClOrdId(algoClOrdId string) *AmendAlgoOrderService {
	s.req.AlgoClOrdId = algoClOrdId
	return s
}

func (s *AmendAlgoOrderService) CxlOnFail(enable bool) *AmendAlgoOrderService {
	s.req.CxlOnFail = &enable
	return s
}

func (s *AmendAlgoOrderService) ReqId(reqId string) *AmendAlgoOrderService {
	s.req.ReqId = reqId
	return s
}

func (s *AmendAlgoOrderService) NewSz(newSz string) *AmendAlgoOrderService {
	s.req.NewSz = newSz
	return s
}

func (s *AmendAlgoOrderService) NewTpTriggerPx(px string) *AmendAlgoOrderService {
	s.req.NewTpTriggerPx = px
	return s
}

func (s *AmendAlgoOrderService) NewTpOrdPx(px string) *AmendAlgoOrderService {
	s.req.NewTpOrdPx = px
	return s
}

func (s *AmendAlgoOrderService) NewSlTriggerPx(px string) *AmendAlgoOrderService {
	s.req.NewSlTriggerPx = px
	return s
}

func (s *AmendAlgoOrderService) NewSlOrdPx(px string) *AmendAlgoOrderService {
	s.req.NewSlOrdPx = px
	return s
}

func (s *AmendAlgoOrderService) NewTpTriggerPxType(typ string) *AmendAlgoOrderService {
	s.req.NewTpTriggerPxType = typ
	return s
}

func (s *AmendAlgoOrderService) NewSlTriggerPxType(typ string) *AmendAlgoOrderService {
	s.req.NewSlTriggerPxType = typ
	return s
}

func (s *AmendAlgoOrderService) NewTriggerPx(px string) *AmendAlgoOrderService {
	s.req.NewTriggerPx = px
	return s
}

func (s *AmendAlgoOrderService) NewOrdPx(px string) *AmendAlgoOrderService {
	s.req.NewOrdPx = px
	return s
}

func (s *AmendAlgoOrderService) NewTriggerPxType(typ string) *AmendAlgoOrderService {
	s.req.NewTriggerPxType = typ
	return s
}

func (s *AmendAlgoOrderService) AttachAlgoOrds(attach []TradeAlgoOrderAttachAmend) *AmendAlgoOrderService {
	s.req.AttachAlgoOrds = attach
	return s
}

var (
	errAmendAlgoOrderMissingInstId    = errors.New("okx: amend algos requires instId")
	errAmendAlgoOrderMissingId        = errors.New("okx: amend algos requires algoId or algoClOrdId")
	errAmendAlgoOrderMissingAnyChange = errors.New("okx: amend algos requires at least one change")
	errEmptyAmendAlgoOrderResponse    = errors.New("okx: empty amend algos response")
	errInvalidAmendAlgoOrderResponse  = errors.New("okx: invalid amend algos response")
)

func (s *AmendAlgoOrderService) hasAnyChange() bool {
	if s.req.NewSz != "" ||
		s.req.NewTpTriggerPx != "" || s.req.NewTpOrdPx != "" ||
		s.req.NewSlTriggerPx != "" || s.req.NewSlOrdPx != "" ||
		s.req.NewTpTriggerPxType != "" || s.req.NewSlTriggerPxType != "" ||
		s.req.NewTriggerPx != "" || s.req.NewOrdPx != "" || s.req.NewTriggerPxType != "" ||
		len(s.req.AttachAlgoOrds) > 0 {
		return true
	}
	return false
}

// Do 修改策略委托订单（POST /api/v5/trade/amend-algos）。
func (s *AmendAlgoOrderService) Do(ctx context.Context) (*TradeAlgoOrderAck, error) {
	if s.req.InstId == "" {
		return nil, errAmendAlgoOrderMissingInstId
	}
	if s.req.AlgoId == "" && s.req.AlgoClOrdId == "" {
		return nil, errAmendAlgoOrderMissingId
	}
	if !s.hasAnyChange() {
		return nil, errAmendAlgoOrderMissingAnyChange
	}

	req := s.req
	if req.AlgoId != "" {
		req.AlgoClOrdId = ""
	}

	var data []TradeAlgoOrderAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/trade/amend-algos", nil, req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/trade/amend-algos", requestID, errEmptyAmendAlgoOrderResponse)
	}
	if len(data) != 1 {
		return nil, newInvalidDataAPIError(
			http.MethodPost,
			"/api/v5/trade/amend-algos",
			requestID,
			fmt.Errorf("%w: expected 1 ack, got %d", errInvalidAmendAlgoOrderResponse, len(data)),
		)
	}
	if data[0].SCode != "0" {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodPost,
			RequestPath: "/api/v5/trade/amend-algos",
			RequestID:   requestID,
			Code:        data[0].SCode,
			Message:     data[0].SMsg,
		}
	}
	return &data[0], nil
}
