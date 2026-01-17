package okx

import (
	"context"
	"errors"
	"net/http"
)

// TradeAccountRateLimit 表示账户限速信息。
// 数值字段按 OKX 返回保持为 string（无损）。
type TradeAccountRateLimit struct {
	AccRateLimit     string `json:"accRateLimit"`
	FillRatio        string `json:"fillRatio"`
	MainFillRatio    string `json:"mainFillRatio"`
	NextAccRateLimit string `json:"nextAccRateLimit"`
	TS               int64  `json:"ts,string"`
}

// TradeAccountRateLimitService 获取账户限速信息。
type TradeAccountRateLimitService struct {
	c *Client
}

// NewTradeAccountRateLimitService 创建 TradeAccountRateLimitService。
func (c *Client) NewTradeAccountRateLimitService() *TradeAccountRateLimitService {
	return &TradeAccountRateLimitService{c: c}
}

var errEmptyTradeAccountRateLimitResponse = errors.New("okx: empty trade account rate limit response")

// Do 获取账户限速信息（GET /api/v5/trade/account-rate-limit）。
func (s *TradeAccountRateLimitService) Do(ctx context.Context) (*TradeAccountRateLimit, error) {
	var data []TradeAccountRateLimit
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/trade/account-rate-limit", nil, nil, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyTradeAccountRateLimitResponse
	}

	info := &data[0]
	if err := s.c.applyTradeAccountRateLimit(info); err != nil {
		return info, err
	}
	return info, nil
}
