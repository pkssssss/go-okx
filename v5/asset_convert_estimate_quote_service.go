package okx

import (
	"context"
	"errors"
	"net/http"
)

type assetConvertEstimateQuoteRequest struct {
	BaseCcy  string `json:"baseCcy"`
	QuoteCcy string `json:"quoteCcy"`
	Side     string `json:"side"`
	RfqSz    string `json:"rfqSz"`
	RfqSzCcy string `json:"rfqSzCcy"`

	ClQReqId string `json:"clQReqId,omitempty"`
	Tag      string `json:"tag,omitempty"`
}

// AssetConvertQuote 表示闪兑报价信息。
// 数值字段保持为 string（无损）。
type AssetConvertQuote struct {
	QuoteTime string `json:"quoteTime"`
	TtlMs     string `json:"ttlMs"`

	ClQReqId string `json:"clQReqId"`
	QuoteId  string `json:"quoteId"`

	BaseCcy  string `json:"baseCcy"`
	QuoteCcy string `json:"quoteCcy"`
	Side     string `json:"side"`

	RfqSz    string `json:"rfqSz"`
	RfqSzCcy string `json:"rfqSzCcy"`

	OrigRfqSz string `json:"origRfqSz"`
	BaseSz    string `json:"baseSz"`
	QuoteSz   string `json:"quoteSz"`

	CnvtPx string `json:"cnvtPx"`
}

// AssetConvertEstimateQuoteService 闪兑预估询价。
type AssetConvertEstimateQuoteService struct {
	c   *Client
	req assetConvertEstimateQuoteRequest
}

// NewAssetConvertEstimateQuoteService 创建 AssetConvertEstimateQuoteService。
func (c *Client) NewAssetConvertEstimateQuoteService() *AssetConvertEstimateQuoteService {
	return &AssetConvertEstimateQuoteService{c: c}
}

// BaseCcy 设置交易货币币种（必填）。
func (s *AssetConvertEstimateQuoteService) BaseCcy(baseCcy string) *AssetConvertEstimateQuoteService {
	s.req.BaseCcy = baseCcy
	return s
}

// QuoteCcy 设置计价货币币种（必填）。
func (s *AssetConvertEstimateQuoteService) QuoteCcy(quoteCcy string) *AssetConvertEstimateQuoteService {
	s.req.QuoteCcy = quoteCcy
	return s
}

// Side 设置交易方向（必填：buy/sell，描述 baseCcy 方向）。
func (s *AssetConvertEstimateQuoteService) Side(side string) *AssetConvertEstimateQuoteService {
	s.req.Side = side
	return s
}

// RfqSz 设置询价数量（必填）。
func (s *AssetConvertEstimateQuoteService) RfqSz(rfqSz string) *AssetConvertEstimateQuoteService {
	s.req.RfqSz = rfqSz
	return s
}

// RfqSzCcy 设置询价币种（必填）。
func (s *AssetConvertEstimateQuoteService) RfqSzCcy(rfqSzCcy string) *AssetConvertEstimateQuoteService {
	s.req.RfqSzCcy = rfqSzCcy
	return s
}

// ClQReqId 设置客户端自定义的订单标识（可选：1-32）。
func (s *AssetConvertEstimateQuoteService) ClQReqId(clQReqId string) *AssetConvertEstimateQuoteService {
	s.req.ClQReqId = clQReqId
	return s
}

// Tag 设置订单标签（可选，适用于 broker）。
func (s *AssetConvertEstimateQuoteService) Tag(tag string) *AssetConvertEstimateQuoteService {
	s.req.Tag = tag
	return s
}

var errAssetConvertEstimateQuoteMissingRequired = errors.New("okx: convert estimate quote requires baseCcy/quoteCcy/side/rfqSz/rfqSzCcy")
var errEmptyAssetConvertEstimateQuote = errors.New("okx: empty convert estimate quote response")

// Do 闪兑预估询价（POST /api/v5/asset/convert/estimate-quote）。
func (s *AssetConvertEstimateQuoteService) Do(ctx context.Context) (*AssetConvertQuote, error) {
	if s.req.BaseCcy == "" || s.req.QuoteCcy == "" || s.req.Side == "" || s.req.RfqSz == "" || s.req.RfqSzCcy == "" {
		return nil, errAssetConvertEstimateQuoteMissingRequired
	}

	var data []AssetConvertQuote
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/asset/convert/estimate-quote", nil, s.req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyAssetConvertEstimateQuote
	}
	return &data[0], nil
}
