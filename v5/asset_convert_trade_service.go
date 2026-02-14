package okx

import (
	"context"
	"errors"
	"net/http"
)

type assetConvertTradeRequest struct {
	QuoteId  string `json:"quoteId"`
	BaseCcy  string `json:"baseCcy"`
	QuoteCcy string `json:"quoteCcy"`
	Side     string `json:"side"`
	Sz       string `json:"sz"`
	SzCcy    string `json:"szCcy"`

	ClTReqId string `json:"clTReqId,omitempty"`
	Tag      string `json:"tag,omitempty"`
}

// AssetConvertTrade 表示闪兑成交信息。
// 数值字段保持为 string（无损）。
type AssetConvertTrade struct {
	TradeId  string `json:"tradeId"`
	QuoteId  string `json:"quoteId"`
	ClTReqId string `json:"clTReqId"`
	State    string `json:"state"`

	InstId   string `json:"instId"`
	BaseCcy  string `json:"baseCcy"`
	QuoteCcy string `json:"quoteCcy"`
	Side     string `json:"side"`

	FillPx      string `json:"fillPx"`
	FillBaseSz  string `json:"fillBaseSz"`
	FillQuoteSz string `json:"fillQuoteSz"`
	TS          int64  `json:"ts,string"`
}

// AssetConvertTradeService 闪兑交易。
type AssetConvertTradeService struct {
	c   *Client
	req assetConvertTradeRequest
}

// NewAssetConvertTradeService 创建 AssetConvertTradeService。
func (c *Client) NewAssetConvertTradeService() *AssetConvertTradeService {
	return &AssetConvertTradeService{c: c}
}

// QuoteId 设置报价 ID（必填）。
func (s *AssetConvertTradeService) QuoteId(quoteId string) *AssetConvertTradeService {
	s.req.QuoteId = quoteId
	return s
}

// BaseCcy 设置交易货币币种（必填）。
func (s *AssetConvertTradeService) BaseCcy(baseCcy string) *AssetConvertTradeService {
	s.req.BaseCcy = baseCcy
	return s
}

// QuoteCcy 设置计价货币币种（必填）。
func (s *AssetConvertTradeService) QuoteCcy(quoteCcy string) *AssetConvertTradeService {
	s.req.QuoteCcy = quoteCcy
	return s
}

// Side 设置交易方向（必填：buy/sell，描述 baseCcy 方向）。
func (s *AssetConvertTradeService) Side(side string) *AssetConvertTradeService {
	s.req.Side = side
	return s
}

// Sz 设置用户报价数量（必填）。
func (s *AssetConvertTradeService) Sz(sz string) *AssetConvertTradeService {
	s.req.Sz = sz
	return s
}

// SzCcy 设置用户报价币种（必填）。
func (s *AssetConvertTradeService) SzCcy(szCcy string) *AssetConvertTradeService {
	s.req.SzCcy = szCcy
	return s
}

// ClTReqId 设置用户自定义订单标识（可选：1-32）。
func (s *AssetConvertTradeService) ClTReqId(clTReqId string) *AssetConvertTradeService {
	s.req.ClTReqId = clTReqId
	return s
}

// Tag 设置订单标签（可选，适用于 broker）。
func (s *AssetConvertTradeService) Tag(tag string) *AssetConvertTradeService {
	s.req.Tag = tag
	return s
}

var errAssetConvertTradeMissingRequired = errors.New("okx: convert trade requires quoteId/baseCcy/quoteCcy/side/sz/szCcy")
var errEmptyAssetConvertTrade = errors.New("okx: empty convert trade response")

// Do 闪兑交易（POST /api/v5/asset/convert/trade）。
func (s *AssetConvertTradeService) Do(ctx context.Context) (*AssetConvertTrade, error) {
	if s.req.QuoteId == "" || s.req.BaseCcy == "" || s.req.QuoteCcy == "" || s.req.Side == "" || s.req.Sz == "" || s.req.SzCcy == "" {
		return nil, errAssetConvertTradeMissingRequired
	}

	var data []AssetConvertTrade
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/asset/convert/trade", nil, s.req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/asset/convert/trade", requestID, errEmptyAssetConvertTrade)
	}
	return &data[0], nil
}
