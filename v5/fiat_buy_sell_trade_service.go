package okx

import (
	"context"
	"errors"
	"net/http"
)

type fiatBuySellTradeRequest struct {
	ClOrdId       string `json:"clOrdId"`
	Side          string `json:"side"`
	FromCcy       string `json:"fromCcy"`
	ToCcy         string `json:"toCcy"`
	RfqAmt        string `json:"rfqAmt"`
	RfqCcy        string `json:"rfqCcy"`
	PaymentMethod string `json:"paymentMethod"`
	QuoteId       string `json:"quoteId"`
}

// FiatBuySellTradeService 买卖交易。
type FiatBuySellTradeService struct {
	c *Client

	clOrdId       string
	side          string
	fromCcy       string
	toCcy         string
	rfqAmt        string
	rfqCcy        string
	paymentMethod string
	quoteId       string
}

// NewFiatBuySellTradeService 创建 FiatBuySellTradeService。
func (c *Client) NewFiatBuySellTradeService() *FiatBuySellTradeService {
	return &FiatBuySellTradeService{c: c}
}

// ClOrdId 设置用户自定义订单标识（必填）。
func (s *FiatBuySellTradeService) ClOrdId(clOrdId string) *FiatBuySellTradeService {
	s.clOrdId = clOrdId
	return s
}

// Side 设置交易方向（必填）：buy/sell。
func (s *FiatBuySellTradeService) Side(side string) *FiatBuySellTradeService {
	s.side = side
	return s
}

// FromCcy 设置卖出币种（必填）。
func (s *FiatBuySellTradeService) FromCcy(fromCcy string) *FiatBuySellTradeService {
	s.fromCcy = fromCcy
	return s
}

// ToCcy 设置买入币种（必填）。
func (s *FiatBuySellTradeService) ToCcy(toCcy string) *FiatBuySellTradeService {
	s.toCcy = toCcy
	return s
}

// RfqAmt 设置询价数量（必填）。
func (s *FiatBuySellTradeService) RfqAmt(rfqAmt string) *FiatBuySellTradeService {
	s.rfqAmt = rfqAmt
	return s
}

// RfqCcy 设置询价币种（必填）。
func (s *FiatBuySellTradeService) RfqCcy(rfqCcy string) *FiatBuySellTradeService {
	s.rfqCcy = rfqCcy
	return s
}

// PaymentMethod 设置支付方式（必填）：balance。
func (s *FiatBuySellTradeService) PaymentMethod(paymentMethod string) *FiatBuySellTradeService {
	s.paymentMethod = paymentMethod
	return s
}

// QuoteId 设置报价ID（必填）。
func (s *FiatBuySellTradeService) QuoteId(quoteId string) *FiatBuySellTradeService {
	s.quoteId = quoteId
	return s
}

var (
	errFiatBuySellTradeMissingRequired = errors.New("okx: fiat buy-sell trade requires clOrdId, quoteId, side, fromCcy, toCcy, rfqAmt, rfqCcy, paymentMethod")
	errEmptyFiatBuySellTradeResponse   = errors.New("okx: empty fiat buy-sell trade response")
)

// Do 买卖交易（POST /api/v5/fiat/buy-sell/trade）。
func (s *FiatBuySellTradeService) Do(ctx context.Context) (*FiatBuySellOrder, error) {
	if s.clOrdId == "" || s.quoteId == "" || s.side == "" || s.fromCcy == "" || s.toCcy == "" || s.rfqAmt == "" || s.rfqCcy == "" || s.paymentMethod == "" {
		return nil, errFiatBuySellTradeMissingRequired
	}

	req := fiatBuySellTradeRequest{
		ClOrdId:       s.clOrdId,
		Side:          s.side,
		FromCcy:       s.fromCcy,
		ToCcy:         s.toCcy,
		RfqAmt:        s.rfqAmt,
		RfqCcy:        s.rfqCcy,
		PaymentMethod: s.paymentMethod,
		QuoteId:       s.quoteId,
	}

	var data []FiatBuySellOrder
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/fiat/buy-sell/trade", nil, req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyFiatBuySellTradeResponse
	}
	return &data[0], nil
}
