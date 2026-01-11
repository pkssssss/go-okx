package okx

import (
	"context"
	"net/http"
)

// RFQCounterparty 表示可选择的报价方信息。
type RFQCounterparty struct {
	TraderName string `json:"traderName"`
	TraderCode string `json:"traderCode"`
	Type       string `json:"type"`
}

// RFQCounterpartiesService 获取报价方信息。
type RFQCounterpartiesService struct {
	c *Client
}

// NewRFQCounterpartiesService 创建 RFQCounterpartiesService。
func (c *Client) NewRFQCounterpartiesService() *RFQCounterpartiesService {
	return &RFQCounterpartiesService{c: c}
}

// Do 获取报价方信息（GET /api/v5/rfq/counterparties）。
func (s *RFQCounterpartiesService) Do(ctx context.Context) ([]RFQCounterparty, error) {
	var data []RFQCounterparty
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/rfq/counterparties", nil, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
