package okx

import (
	"context"
	"errors"
	"net/http"
)

type fiatBuySellQuoteRequest struct {
	Side    string `json:"side"`
	FromCcy string `json:"fromCcy"`
	ToCcy   string `json:"toCcy"`
	RfqAmt  string `json:"rfqAmt"`
	RfqCcy  string `json:"rfqCcy"`
}

// FiatBuySellQuoteService 获取买卖交易报价。
type FiatBuySellQuoteService struct {
	c *Client

	side    string
	fromCcy string
	toCcy   string
	rfqAmt  string
	rfqCcy  string
}

// NewFiatBuySellQuoteService 创建 FiatBuySellQuoteService。
func (c *Client) NewFiatBuySellQuoteService() *FiatBuySellQuoteService {
	return &FiatBuySellQuoteService{c: c}
}

// Side 设置交易方向（必填）：buy/sell。
func (s *FiatBuySellQuoteService) Side(side string) *FiatBuySellQuoteService {
	s.side = side
	return s
}

// FromCcy 设置卖出币种（必填）。
func (s *FiatBuySellQuoteService) FromCcy(fromCcy string) *FiatBuySellQuoteService {
	s.fromCcy = fromCcy
	return s
}

// ToCcy 设置买入币种（必填）。
func (s *FiatBuySellQuoteService) ToCcy(toCcy string) *FiatBuySellQuoteService {
	s.toCcy = toCcy
	return s
}

// RfqAmt 设置询价数量（必填）。
func (s *FiatBuySellQuoteService) RfqAmt(rfqAmt string) *FiatBuySellQuoteService {
	s.rfqAmt = rfqAmt
	return s
}

// RfqCcy 设置询价币种（必填）。
func (s *FiatBuySellQuoteService) RfqCcy(rfqCcy string) *FiatBuySellQuoteService {
	s.rfqCcy = rfqCcy
	return s
}

var (
	errFiatBuySellQuoteMissingRequired = errors.New("okx: fiat buy-sell quote requires side, fromCcy, toCcy, rfqAmt, rfqCcy")
	errEmptyFiatBuySellQuoteResponse   = errors.New("okx: empty fiat buy-sell quote response")
)

// Do 获取买卖交易报价（POST /api/v5/fiat/buy-sell/quote）。
func (s *FiatBuySellQuoteService) Do(ctx context.Context) (*FiatBuySellQuote, error) {
	if s.side == "" || s.fromCcy == "" || s.toCcy == "" || s.rfqAmt == "" || s.rfqCcy == "" {
		return nil, errFiatBuySellQuoteMissingRequired
	}

	req := fiatBuySellQuoteRequest{
		Side:    s.side,
		FromCcy: s.fromCcy,
		ToCcy:   s.toCcy,
		RfqAmt:  s.rfqAmt,
		RfqCcy:  s.rfqCcy,
	}

	var data []FiatBuySellQuote
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/fiat/buy-sell/quote", nil, req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/fiat/buy-sell/quote", requestID, errEmptyFiatBuySellQuoteResponse)
	}
	return &data[0], nil
}
