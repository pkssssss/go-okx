package okx

import (
	"context"
	"errors"
	"net/http"
)

type financeStakingDefiCancelRequest struct {
	OrdId        string `json:"ordId"`
	ProtocolType string `json:"protocolType"`
}

// FinanceStakingDefiCancelService 撤销链上赚币项目申购/赎回。
type FinanceStakingDefiCancelService struct {
	c   *Client
	req financeStakingDefiCancelRequest
}

// NewFinanceStakingDefiCancelService 创建 FinanceStakingDefiCancelService。
func (c *Client) NewFinanceStakingDefiCancelService() *FinanceStakingDefiCancelService {
	return &FinanceStakingDefiCancelService{c: c}
}

// OrdId 设置订单 ID（必填）。
func (s *FinanceStakingDefiCancelService) OrdId(ordId string) *FinanceStakingDefiCancelService {
	s.req.OrdId = ordId
	return s
}

// ProtocolType 设置项目类型（必填；defi）。
func (s *FinanceStakingDefiCancelService) ProtocolType(protocolType string) *FinanceStakingDefiCancelService {
	s.req.ProtocolType = protocolType
	return s
}

var (
	errFinanceStakingDefiCancelMissingRequired = errors.New("okx: staking-defi cancel requires ordId and protocolType")
	errEmptyFinanceStakingDefiCancelAck        = errors.New("okx: empty staking-defi cancel response")
)

// Do 撤销链上赚币项目申购/赎回（POST /api/v5/finance/staking-defi/cancel）。
func (s *FinanceStakingDefiCancelService) Do(ctx context.Context) (*FinanceStakingDefiOrderAck, error) {
	if s.req.OrdId == "" || s.req.ProtocolType == "" {
		return nil, errFinanceStakingDefiCancelMissingRequired
	}

	var data []FinanceStakingDefiOrderAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/finance/staking-defi/cancel", nil, s.req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/finance/staking-defi/cancel", requestID, errEmptyFinanceStakingDefiCancelAck)
	}
	return &data[0], nil
}
