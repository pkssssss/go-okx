package okx

import (
	"context"
	"net/http"
	"net/url"
)

// InterestRateLoanQuotaBasic 表示基础利率和借币限额。
//
// 说明：数值字段保持为 string（无损）。
type InterestRateLoanQuotaBasic struct {
	Ccy   string `json:"ccy"`
	Rate  string `json:"rate"`
	Quota string `json:"quota"`
}

// InterestRateLoanQuotaLevelCoef 表示借币限额系数（专业/普通用户）。
//
// 说明：数值字段保持为 string（无损）。
type InterestRateLoanQuotaLevelCoef struct {
	Level         string `json:"level"`
	LoanQuotaCoef string `json:"loanQuotaCoef"`
	IRDiscount    string `json:"irDiscount"`
}

// InterestRateLoanQuotaConfigCcy 表示自定义绝对值方式配置借币限额的币种及基础利率。
//
// 说明：数值字段保持为 string（无损）。
type InterestRateLoanQuotaConfigCcy struct {
	Ccy  string `json:"ccy"`
	Rate string `json:"rate"`
}

// InterestRateLoanQuotaConfig 表示自定义绝对值方式配置借币限额的币种详情。
//
// 说明：数值字段保持为 string（无损）。
type InterestRateLoanQuotaConfig struct {
	Ccy      string `json:"ccy"`
	StgyType string `json:"stgyType"`
	Quota    string `json:"quota"`
	Level    string `json:"level"`
}

// InterestRateLoanQuota 表示市场借币杠杆利率和借币限额。
type InterestRateLoanQuota struct {
	Basic   []InterestRateLoanQuotaBasic     `json:"basic"`
	VIP     []InterestRateLoanQuotaLevelCoef `json:"vip"`
	Regular []InterestRateLoanQuotaLevelCoef `json:"regular"`

	ConfigCcyList []InterestRateLoanQuotaConfigCcy `json:"configCcyList"`
	Config        []InterestRateLoanQuotaConfig    `json:"config"`
}

// PublicInterestRateLoanQuotaService 获取市场借币杠杆利率和借币限额。
type PublicInterestRateLoanQuotaService struct {
	c *Client
}

// NewPublicInterestRateLoanQuotaService 创建 PublicInterestRateLoanQuotaService。
func (c *Client) NewPublicInterestRateLoanQuotaService() *PublicInterestRateLoanQuotaService {
	return &PublicInterestRateLoanQuotaService{c: c}
}

// Do 获取市场借币杠杆利率和借币限额（GET /api/v5/public/interest-rate-loan-quota）。
func (s *PublicInterestRateLoanQuotaService) Do(ctx context.Context) ([]InterestRateLoanQuota, error) {
	q := url.Values{}

	var data []InterestRateLoanQuota
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/public/interest-rate-loan-quota", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
