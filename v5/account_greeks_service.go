package okx

import (
	"context"
	"net/http"
	"net/url"
)

// AccountGreeks 表示账户资产 Greeks 信息。
type AccountGreeks struct {
	DeltaBS string `json:"deltaBS"`
	DeltaPA string `json:"deltaPA"`
	GammaBS string `json:"gammaBS"`
	GammaPA string `json:"gammaPA"`
	ThetaBS string `json:"thetaBS"`
	ThetaPA string `json:"thetaPA"`
	VegaBS  string `json:"vegaBS"`
	VegaPA  string `json:"vegaPA"`

	Ccy string    `json:"ccy"`
	TS  UnixMilli `json:"ts"`
}

// AccountGreeksService 查看账户 Greeks。
type AccountGreeksService struct {
	c   *Client
	ccy string
}

// NewAccountGreeksService 创建 AccountGreeksService。
func (c *Client) NewAccountGreeksService() *AccountGreeksService {
	return &AccountGreeksService{c: c}
}

// Ccy 设置币种（可选，如 BTC）。
func (s *AccountGreeksService) Ccy(ccy string) *AccountGreeksService {
	s.ccy = ccy
	return s
}

// Do 查看账户 Greeks（GET /api/v5/account/greeks）。
func (s *AccountGreeksService) Do(ctx context.Context) ([]AccountGreeks, error) {
	var q url.Values
	if s.ccy != "" {
		q = url.Values{}
		q.Set("ccy", s.ccy)
	}

	var data []AccountGreeks
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/greeks", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
