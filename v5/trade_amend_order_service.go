package okx

import (
	"context"
	"errors"
	"net/http"
)

// AmendOrderService 修改订单。
type AmendOrderService struct {
	c *Client

	instId      string
	ordId       string
	clOrdId     string
	reqId       string
	newSz       string
	newPx       string
	newPxUsd    string
	newPxVol    string
	pxAmendType string
	cxlOnFail   *bool

	expTimeHeader string
}

// NewAmendOrderService 创建 AmendOrderService。
func (c *Client) NewAmendOrderService() *AmendOrderService {
	return &AmendOrderService{c: c}
}

func (s *AmendOrderService) InstId(instId string) *AmendOrderService {
	s.instId = instId
	return s
}

func (s *AmendOrderService) OrdId(ordId string) *AmendOrderService {
	s.ordId = ordId
	return s
}

func (s *AmendOrderService) ClOrdId(clOrdId string) *AmendOrderService {
	s.clOrdId = clOrdId
	return s
}

// ReqId 设置用户自定义修改事件 ID（1-32 位）。
func (s *AmendOrderService) ReqId(reqId string) *AmendOrderService {
	s.reqId = reqId
	return s
}

func (s *AmendOrderService) NewSz(newSz string) *AmendOrderService {
	s.newSz = newSz
	return s
}

func (s *AmendOrderService) NewPx(newPx string) *AmendOrderService {
	s.newPx = newPx
	return s
}

func (s *AmendOrderService) NewPxUsd(newPxUsd string) *AmendOrderService {
	s.newPxUsd = newPxUsd
	return s
}

func (s *AmendOrderService) NewPxVol(newPxVol string) *AmendOrderService {
	s.newPxVol = newPxVol
	return s
}

func (s *AmendOrderService) PxAmendType(pxAmendType string) *AmendOrderService {
	s.pxAmendType = pxAmendType
	return s
}

// CxlOnFail 当改单失败时，是否自动撤单（默认 false）。
func (s *AmendOrderService) CxlOnFail(enable bool) *AmendOrderService {
	s.cxlOnFail = &enable
	return s
}

// ExpTime 设置 REST 请求头 expTime（请求有效截止时间，Unix 毫秒时间戳字符串）。
func (s *AmendOrderService) ExpTime(expTimeMillis string) *AmendOrderService {
	s.expTimeHeader = expTimeMillis
	return s
}

var (
	errAmendOrderMissingInstId = errors.New("okx: amend order requires instId")
	errAmendOrderMissingId     = errors.New("okx: amend order requires ordId or clOrdId")
	errAmendOrderTooManyId     = errors.New("okx: amend order requires exactly one of ordId or clOrdId")
	errAmendOrderMissingChange = errors.New("okx: amend order requires newSz or newPx/newPxUsd/newPxVol")
	errAmendOrderTooManyPx     = errors.New("okx: amend order requires at most one of newPx/newPxUsd/newPxVol")
	errEmptyAmendOrderResponse = errors.New("okx: empty amend order response")
)

type amendOrderRequest struct {
	InstId      string `json:"instId"`
	CxlOnFail   *bool  `json:"cxlOnFail,omitempty"`
	OrdId       string `json:"ordId,omitempty"`
	ClOrdId     string `json:"clOrdId,omitempty"`
	ReqId       string `json:"reqId,omitempty"`
	NewSz       string `json:"newSz,omitempty"`
	NewPx       string `json:"newPx,omitempty"`
	NewPxUsd    string `json:"newPxUsd,omitempty"`
	NewPxVol    string `json:"newPxVol,omitempty"`
	PxAmendType string `json:"pxAmendType,omitempty"`
}

// Do 修改订单（POST /api/v5/trade/amend-order）。
func (s *AmendOrderService) Do(ctx context.Context) (*TradeOrderAck, error) {
	if s.instId == "" {
		return nil, errAmendOrderMissingInstId
	}
	switch countNonEmptyStrings(s.ordId, s.clOrdId) {
	case 0:
		return nil, errAmendOrderMissingId
	case 2:
		return nil, errAmendOrderTooManyId
	}

	if countNonEmptyStrings(s.newPx, s.newPxUsd, s.newPxVol) > 1 {
		return nil, errAmendOrderTooManyPx
	}
	if s.newSz == "" && s.newPx == "" && s.newPxUsd == "" && s.newPxVol == "" {
		return nil, errAmendOrderMissingChange
	}

	req := amendOrderRequest{
		InstId:      s.instId,
		CxlOnFail:   s.cxlOnFail,
		OrdId:       s.ordId,
		ClOrdId:     s.clOrdId,
		ReqId:       s.reqId,
		NewSz:       s.newSz,
		NewPx:       s.newPx,
		NewPxUsd:    s.newPxUsd,
		NewPxVol:    s.newPxVol,
		PxAmendType: s.pxAmendType,
	}

	var data []TradeOrderAck
	var header http.Header
	if s.expTimeHeader != "" {
		header = make(http.Header)
		header.Set("expTime", s.expTimeHeader)
	}
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/trade/amend-order", nil, req, true, header, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodPost,
			RequestPath: "/api/v5/trade/amend-order",
			RequestID:   requestID,
			Code:        "0",
			Message:     errEmptyAmendOrderResponse.Error(),
		}
	}
	if data[0].SCode != "0" {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodPost,
			RequestPath: "/api/v5/trade/amend-order",
			RequestID:   requestID,
			Code:        data[0].SCode,
			Message:     data[0].SMsg,
		}
	}
	return &data[0], nil
}
