package okx

import (
	"context"
	"net/http"
)

// OneClickRepayCurrencyListV2Service 获取一键还债币种列表（新）。
type OneClickRepayCurrencyListV2Service struct {
	c *Client
}

// NewOneClickRepayCurrencyListV2Service 创建 OneClickRepayCurrencyListV2Service。
func (c *Client) NewOneClickRepayCurrencyListV2Service() *OneClickRepayCurrencyListV2Service {
	return &OneClickRepayCurrencyListV2Service{c: c}
}

// Do 获取一键还债币种列表（新）（GET /api/v5/trade/one-click-repay-currency-list-v2）。
func (s *OneClickRepayCurrencyListV2Service) Do(ctx context.Context) ([]OneClickRepayCurrencyListV2Item, error) {
	var data []OneClickRepayCurrencyListV2Item
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/trade/one-click-repay-currency-list-v2", nil, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
