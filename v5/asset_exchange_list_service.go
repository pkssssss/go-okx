package okx

import (
	"context"
	"net/http"
)

// AssetExchange 表示交易所信息（用于接收方信息 rcvrInfo.exchId 等）。
type AssetExchange struct {
	ExchName string `json:"exchName"`
	ExchId   string `json:"exchId"`
}

// AssetExchangeListService 获取交易所列表（公共）。
type AssetExchangeListService struct {
	c *Client
}

// NewAssetExchangeListService 创建 AssetExchangeListService。
func (c *Client) NewAssetExchangeListService() *AssetExchangeListService {
	return &AssetExchangeListService{c: c}
}

// Do 获取交易所列表（GET /api/v5/asset/exchange-list）。
func (s *AssetExchangeListService) Do(ctx context.Context) ([]AssetExchange, error) {
	var data []AssetExchange
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/asset/exchange-list", nil, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
