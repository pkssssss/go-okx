package okx

import (
	"context"
	"errors"
	"net/http"
)

// MarketPlatform24Volume 表示平台 24 小时总成交量。
// 数值字段保持为 string（无损）。
type MarketPlatform24Volume struct {
	VolUsd string `json:"volUsd"`
	VolCny string `json:"volCny"`
	TS     int64  `json:"ts,string"`
}

// MarketPlatform24VolumeService 获取平台 24 小时总成交量。
type MarketPlatform24VolumeService struct {
	c *Client
}

// NewMarketPlatform24VolumeService 创建 MarketPlatform24VolumeService。
func (c *Client) NewMarketPlatform24VolumeService() *MarketPlatform24VolumeService {
	return &MarketPlatform24VolumeService{c: c}
}

var errEmptyMarketPlatform24VolumeResponse = errors.New("okx: empty market platform 24 volume response")

// Do 获取平台 24 小时总成交量（GET /api/v5/market/platform-24-volume）。
func (s *MarketPlatform24VolumeService) Do(ctx context.Context) (*MarketPlatform24Volume, error) {
	var data []MarketPlatform24Volume
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/market/platform-24-volume", nil, nil, false, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyMarketPlatform24VolumeResponse
	}
	return &data[0], nil
}
