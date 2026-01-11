package okx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// RFQExecuteQuoteLeg 表示执行报价时的腿（用于部分执行）。
type RFQExecuteQuoteLeg struct {
	Sz     string `json:"sz"`
	InstId string `json:"instId"`
}

type rfqExecuteQuoteRequest struct {
	RfqId   string               `json:"rfqId"`
	QuoteId string               `json:"quoteId"`
	Legs    []RFQExecuteQuoteLeg `json:"legs,omitempty"`
}

// RFQExecuteQuoteService 执行报价。
type RFQExecuteQuoteService struct {
	c *Client

	rfqId   string
	quoteId string
	legs    []RFQExecuteQuoteLeg
}

// NewRFQExecuteQuoteService 创建 RFQExecuteQuoteService。
func (c *Client) NewRFQExecuteQuoteService() *RFQExecuteQuoteService {
	return &RFQExecuteQuoteService{c: c}
}

// RfqId 设置询价单 ID（必填）。
func (s *RFQExecuteQuoteService) RfqId(rfqId string) *RFQExecuteQuoteService {
	s.rfqId = rfqId
	return s
}

// QuoteId 设置报价单 ID（必填）。
func (s *RFQExecuteQuoteService) QuoteId(quoteId string) *RFQExecuteQuoteService {
	s.quoteId = quoteId
	return s
}

// Legs 设置执行腿（可选；用于部分执行）。
func (s *RFQExecuteQuoteService) Legs(legs []RFQExecuteQuoteLeg) *RFQExecuteQuoteService {
	s.legs = legs
	return s
}

var (
	errRFQExecuteQuoteMissingRequired = errors.New("okx: rfq execute quote requires rfqId and quoteId")
	errEmptyRFQExecuteQuoteResponse   = errors.New("okx: empty rfq execute quote response")
)

// Do 执行报价（POST /api/v5/rfq/execute-quote）。
func (s *RFQExecuteQuoteService) Do(ctx context.Context) (*StrucBlockTrade, error) {
	if s.rfqId == "" || s.quoteId == "" {
		return nil, errRFQExecuteQuoteMissingRequired
	}

	for i, leg := range s.legs {
		if leg.InstId == "" {
			return nil, fmt.Errorf("okx: rfq execute quote legs[%d] missing instId", i)
		}
		if leg.Sz == "" {
			return nil, fmt.Errorf("okx: rfq execute quote legs[%d] missing sz", i)
		}
	}

	req := rfqExecuteQuoteRequest{
		RfqId:   s.rfqId,
		QuoteId: s.quoteId,
		Legs:    s.legs,
	}

	var data []StrucBlockTrade
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/rfq/execute-quote", nil, req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyRFQExecuteQuoteResponse
	}
	return &data[0], nil
}
