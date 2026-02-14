package okx

import (
	"context"
	"errors"
	"net/http"
)

// TradeClosePositionAck 表示市价仓位全平的返回项。
type TradeClosePositionAck struct {
	ClOrdId string `json:"clOrdId"`
	InstId  string `json:"instId"`
	PosSide string `json:"posSide"`
	Tag     string `json:"tag"`
}

// ClosePositionsService 市价仓位全平。
type ClosePositionsService struct {
	c *Client

	instId  string
	posSide string
	mgnMode string
	ccy     string

	autoCxl *bool
	clOrdId string
	tag     string
}

// NewClosePositionsService 创建 ClosePositionsService。
func (c *Client) NewClosePositionsService() *ClosePositionsService {
	return &ClosePositionsService{c: c}
}

// InstId 设置产品 ID（必填）。
func (s *ClosePositionsService) InstId(instId string) *ClosePositionsService {
	s.instId = instId
	return s
}

// PosSide 设置持仓方向（可选；开平仓模式下必填：long/short；买卖模式下可省略或填 net）。
func (s *ClosePositionsService) PosSide(posSide string) *ClosePositionsService {
	s.posSide = posSide
	return s
}

// MgnMode 设置保证金模式（必填：cross/isolated）。
func (s *ClosePositionsService) MgnMode(mgnMode string) *ClosePositionsService {
	s.mgnMode = mgnMode
	return s
}

// Ccy 设置保证金币种（可选；合约模式下的全仓币币杠杆平仓必填）。
func (s *ClosePositionsService) Ccy(ccy string) *ClosePositionsService {
	s.ccy = ccy
	return s
}

// AutoCxl 设置当市价全平时是否自动撤销平仓挂单（可选，默认 false）。
func (s *ClosePositionsService) AutoCxl(enable bool) *ClosePositionsService {
	s.autoCxl = &enable
	return s
}

// ClOrdId 设置客户自定义 ID（可选，1-32）。
func (s *ClosePositionsService) ClOrdId(clOrdId string) *ClosePositionsService {
	s.clOrdId = clOrdId
	return s
}

// Tag 设置订单标签（可选，1-16）。
func (s *ClosePositionsService) Tag(tag string) *ClosePositionsService {
	s.tag = tag
	return s
}

var (
	errClosePositionsMissingRequired = errors.New("okx: close positions requires instId/mgnMode")
	errEmptyClosePositionsResponse   = errors.New("okx: empty close positions response")
	errInvalidClosePositionsResponse = errors.New("okx: invalid close positions response")
)

type closePositionsRequest struct {
	InstId  string `json:"instId"`
	PosSide string `json:"posSide,omitempty"`
	MgnMode string `json:"mgnMode"`
	Ccy     string `json:"ccy,omitempty"`
	AutoCxl *bool  `json:"autoCxl,omitempty"`
	ClOrdId string `json:"clOrdId,omitempty"`
	Tag     string `json:"tag,omitempty"`
}

func validateClosePositionsAck(ack *TradeClosePositionAck, req closePositionsRequest) error {
	if ack == nil || ack.InstId == "" {
		return errInvalidClosePositionsResponse
	}
	if ack.InstId != req.InstId {
		return errInvalidClosePositionsResponse
	}
	if req.PosSide != "" && ack.PosSide != req.PosSide {
		return errInvalidClosePositionsResponse
	}
	if req.ClOrdId != "" && ack.ClOrdId != req.ClOrdId {
		return errInvalidClosePositionsResponse
	}
	if req.Tag != "" && ack.Tag != req.Tag {
		return errInvalidClosePositionsResponse
	}
	return nil
}

// Do 市价仓位全平（POST /api/v5/trade/close-position）。
func (s *ClosePositionsService) Do(ctx context.Context) ([]TradeClosePositionAck, error) {
	if s.instId == "" || s.mgnMode == "" {
		return nil, errClosePositionsMissingRequired
	}

	req := closePositionsRequest{
		InstId:  s.instId,
		PosSide: s.posSide,
		MgnMode: s.mgnMode,
		Ccy:     s.ccy,
		AutoCxl: s.autoCxl,
		ClOrdId: s.clOrdId,
		Tag:     s.tag,
	}

	var data []TradeClosePositionAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/trade/close-position", nil, req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/trade/close-position", requestID, errEmptyClosePositionsResponse)
	}
	for i := range data {
		if err := validateClosePositionsAck(&data[i], req); err != nil {
			return nil, err
		}
	}
	return data, nil
}
