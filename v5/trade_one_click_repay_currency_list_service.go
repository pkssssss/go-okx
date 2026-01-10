package okx

import (
	"context"
	"net/http"
	"net/url"
)

// OneClickRepayCurrencyListService 获取一键还债币种列表（跨币种保证金/组合保证金）。
type OneClickRepayCurrencyListService struct {
	c *Client

	debtType string
}

// NewOneClickRepayCurrencyListService 创建 OneClickRepayCurrencyListService。
func (c *Client) NewOneClickRepayCurrencyListService() *OneClickRepayCurrencyListService {
	return &OneClickRepayCurrencyListService{c: c}
}

// DebtType 设置负债类型（可选：cross=全仓，isolated=逐仓）。
func (s *OneClickRepayCurrencyListService) DebtType(debtType string) *OneClickRepayCurrencyListService {
	s.debtType = debtType
	return s
}

// Do 获取一键还债币种列表（跨币种保证金/组合保证金）（GET /api/v5/trade/one-click-repay-currency-list）。
func (s *OneClickRepayCurrencyListService) Do(ctx context.Context) ([]OneClickRepayCurrencyList, error) {
	var q url.Values
	if s.debtType != "" {
		q = url.Values{}
		q.Set("debtType", s.debtType)
	}

	var data []OneClickRepayCurrencyList
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/trade/one-click-repay-currency-list", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
