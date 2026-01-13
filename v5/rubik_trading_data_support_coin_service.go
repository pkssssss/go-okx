package okx

import (
	"context"
	"net/http"
)

// RubikSupportCoinService 获取交易大数据支持币种。
type RubikSupportCoinService struct {
	c *Client
}

// NewRubikSupportCoinService 创建 RubikSupportCoinService。
func (c *Client) NewRubikSupportCoinService() *RubikSupportCoinService {
	return &RubikSupportCoinService{c: c}
}

// Do 获取交易大数据支持币种（GET /api/v5/rubik/stat/trading-data/support-coin）。
func (s *RubikSupportCoinService) Do(ctx context.Context) (*RubikSupportCoin, error) {
	var data RubikSupportCoin
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/rubik/stat/trading-data/support-coin", nil, nil, false, &data); err != nil {
		return nil, err
	}
	return &data, nil
}
