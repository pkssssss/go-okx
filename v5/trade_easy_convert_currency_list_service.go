package okx

import (
	"context"
	"net/http"
	"net/url"
)

// EasyConvertCurrencyListService 获取一键兑换主流币币种列表。
type EasyConvertCurrencyListService struct {
	c *Client

	source string
}

// NewEasyConvertCurrencyListService 创建 EasyConvertCurrencyListService。
func (c *Client) NewEasyConvertCurrencyListService() *EasyConvertCurrencyListService {
	return &EasyConvertCurrencyListService{c: c}
}

// Source 设置资金来源（可选：1=交易账户，2=资金账户）。
func (s *EasyConvertCurrencyListService) Source(source string) *EasyConvertCurrencyListService {
	s.source = source
	return s
}

// Do 获取一键兑换主流币币种列表（GET /api/v5/trade/easy-convert-currency-list）。
func (s *EasyConvertCurrencyListService) Do(ctx context.Context) ([]EasyConvertCurrencyList, error) {
	var q url.Values
	if s.source != "" {
		q = url.Values{}
		q.Set("source", s.source)
	}

	var data []EasyConvertCurrencyList
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/trade/easy-convert-currency-list", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
