package okx

import (
	"context"
	"net/http"
	"net/url"
)

// DiscountRateInterestFreeQuotaDetail 表示币种折算率详情。
//
// 说明：数值字段保持为 string（无损）。
type DiscountRateInterestFreeQuotaDetail struct {
	DiscountRate   string `json:"discountRate"`
	LiqPenaltyRate string `json:"liqPenaltyRate"`
	MaxAmt         string `json:"maxAmt"`
	MinAmt         string `json:"minAmt"`
	Tier           string `json:"tier"`
	DisCcyEq       string `json:"disCcyEq"`
}

// DiscountRateInterestFreeQuota 表示免息额度和币种折算率等级。
//
// 说明：数值字段保持为 string（无损）。
type DiscountRateInterestFreeQuota struct {
	Ccy string `json:"ccy"`

	ColRes             string `json:"colRes"`
	CollateralRestrict bool   `json:"collateralRestrict"`

	Amt             string `json:"amt"`
	DiscountLv      string `json:"discountLv"`
	MinDiscountRate string `json:"minDiscountRate"`

	Details []DiscountRateInterestFreeQuotaDetail `json:"details"`
}

// PublicDiscountRateInterestFreeQuotaService 获取免息额度和币种折算率等级。
type PublicDiscountRateInterestFreeQuotaService struct {
	c *Client

	ccy        string
	discountLv string
}

// NewPublicDiscountRateInterestFreeQuotaService 创建 PublicDiscountRateInterestFreeQuotaService。
func (c *Client) NewPublicDiscountRateInterestFreeQuotaService() *PublicDiscountRateInterestFreeQuotaService {
	return &PublicDiscountRateInterestFreeQuotaService{c: c}
}

// Ccy 设置币种（可选）。
func (s *PublicDiscountRateInterestFreeQuotaService) Ccy(ccy string) *PublicDiscountRateInterestFreeQuotaService {
	s.ccy = ccy
	return s
}

// DiscountLv 设置折算率等级（已弃用；兼容文档参数名）。
func (s *PublicDiscountRateInterestFreeQuotaService) DiscountLv(discountLv string) *PublicDiscountRateInterestFreeQuotaService {
	s.discountLv = discountLv
	return s
}

// Do 获取免息额度和币种折算率等级（GET /api/v5/public/discount-rate-interest-free-quota）。
func (s *PublicDiscountRateInterestFreeQuotaService) Do(ctx context.Context) ([]DiscountRateInterestFreeQuota, error) {
	q := url.Values{}
	if s.ccy != "" {
		q.Set("ccy", s.ccy)
	}
	if s.discountLv != "" {
		q.Set("discountLv", s.discountLv)
	}

	var data []DiscountRateInterestFreeQuota
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/public/discount-rate-interest-free-quota", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
