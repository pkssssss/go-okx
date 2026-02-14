package okx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type rfqCreateQuoteRequest struct {
	RfqId     string     `json:"rfqId"`
	ClQuoteId string     `json:"clQuoteId,omitempty"`
	Tag       string     `json:"tag,omitempty"`
	QuoteSide string     `json:"quoteSide"`
	Anonymous *bool      `json:"anonymous,omitempty"`
	ExpiresIn string     `json:"expiresIn,omitempty"`
	Legs      []QuoteLeg `json:"legs"`
}

// RFQCreateQuoteService 创建报价单。
type RFQCreateQuoteService struct {
	c *Client

	rfqId     string
	clQuoteId string
	tag       string
	quoteSide string
	anonymous *bool
	expiresIn string
	legs      []QuoteLeg
}

// NewRFQCreateQuoteService 创建 RFQCreateQuoteService。
func (c *Client) NewRFQCreateQuoteService() *RFQCreateQuoteService {
	return &RFQCreateQuoteService{c: c}
}

// RfqId 设置询价单 ID（必填）。
func (s *RFQCreateQuoteService) RfqId(rfqId string) *RFQCreateQuoteService {
	s.rfqId = rfqId
	return s
}

// ClQuoteId 设置报价单自定义 ID（可选）。
func (s *RFQCreateQuoteService) ClQuoteId(clQuoteId string) *RFQCreateQuoteService {
	s.clQuoteId = clQuoteId
	return s
}

// Tag 设置报价单标签（可选）。
func (s *RFQCreateQuoteService) Tag(tag string) *RFQCreateQuoteService {
	s.tag = tag
	return s
}

// QuoteSide 设置报价单方向（必填：buy/sell）。
func (s *RFQCreateQuoteService) QuoteSide(quoteSide string) *RFQCreateQuoteService {
	s.quoteSide = quoteSide
	return s
}

// Anonymous 设置是否匿名报价（可选；默认 false）。
func (s *RFQCreateQuoteService) Anonymous(anonymous bool) *RFQCreateQuoteService {
	s.anonymous = &anonymous
	return s
}

// ExpiresIn 设置报价有效期（秒字符串，可选）。
func (s *RFQCreateQuoteService) ExpiresIn(expiresIn string) *RFQCreateQuoteService {
	s.expiresIn = expiresIn
	return s
}

// Legs 设置报价腿列表（必填，最多 15 条）。
func (s *RFQCreateQuoteService) Legs(legs []QuoteLeg) *RFQCreateQuoteService {
	s.legs = legs
	return s
}

var (
	errRFQCreateQuoteMissingRequired = errors.New("okx: rfq create quote requires rfqId/quoteSide/legs")
	errRFQCreateQuoteTooManyLegs     = errors.New("okx: rfq create quote max 15 legs")
	errEmptyRFQCreateQuoteResponse   = errors.New("okx: empty rfq create quote response")
)

// Do 创建报价单（POST /api/v5/rfq/create-quote）。
func (s *RFQCreateQuoteService) Do(ctx context.Context) (*Quote, error) {
	if s.rfqId == "" || s.quoteSide == "" || len(s.legs) == 0 {
		return nil, errRFQCreateQuoteMissingRequired
	}
	if len(s.legs) > rfqMaxLegs {
		return nil, errRFQCreateQuoteTooManyLegs
	}

	for i, leg := range s.legs {
		if leg.InstId == "" {
			return nil, fmt.Errorf("okx: rfq create quote legs[%d] missing instId", i)
		}
		if leg.Sz == "" {
			return nil, fmt.Errorf("okx: rfq create quote legs[%d] missing sz", i)
		}
		if leg.Px == "" {
			return nil, fmt.Errorf("okx: rfq create quote legs[%d] missing px", i)
		}
		if leg.Side == "" {
			return nil, fmt.Errorf("okx: rfq create quote legs[%d] missing side", i)
		}
	}

	req := rfqCreateQuoteRequest{
		RfqId:     s.rfqId,
		ClQuoteId: s.clQuoteId,
		Tag:       s.tag,
		QuoteSide: s.quoteSide,
		Anonymous: s.anonymous,
		ExpiresIn: s.expiresIn,
		Legs:      s.legs,
	}

	var data []Quote
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/rfq/create-quote", nil, req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/rfq/create-quote", requestID, errEmptyRFQCreateQuoteResponse)
	}
	return &data[0], nil
}
