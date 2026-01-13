package okx

import (
	"context"
	"errors"
	"net/http"
)

type financeStakingDefiPurchaseRequest struct {
	ProductId  string                         `json:"productId"`
	InvestData []FinanceStakingDefiInvestData `json:"investData"`
	Term       string                         `json:"term,omitempty"`
	Tag        string                         `json:"tag,omitempty"`
}

// FinanceStakingDefiPurchaseService 申购链上赚币项目。
type FinanceStakingDefiPurchaseService struct {
	c   *Client
	req financeStakingDefiPurchaseRequest
}

// NewFinanceStakingDefiPurchaseService 创建 FinanceStakingDefiPurchaseService。
func (c *Client) NewFinanceStakingDefiPurchaseService() *FinanceStakingDefiPurchaseService {
	return &FinanceStakingDefiPurchaseService{c: c}
}

// ProductId 设置项目 ID（必填）。
func (s *FinanceStakingDefiPurchaseService) ProductId(productId string) *FinanceStakingDefiPurchaseService {
	s.req.ProductId = productId
	return s
}

// InvestData 设置投资信息（必填）。
func (s *FinanceStakingDefiPurchaseService) InvestData(investData []FinanceStakingDefiInvestData) *FinanceStakingDefiPurchaseService {
	s.req.InvestData = investData
	return s
}

// Term 设置投资期限（可选；定期项目必须指定）。
func (s *FinanceStakingDefiPurchaseService) Term(term string) *FinanceStakingDefiPurchaseService {
	s.req.Term = term
	return s
}

// Tag 设置订单标签（可选，1-16 位字母/数字组合）。
func (s *FinanceStakingDefiPurchaseService) Tag(tag string) *FinanceStakingDefiPurchaseService {
	s.req.Tag = tag
	return s
}

var (
	errFinanceStakingDefiPurchaseMissingProductId  = errors.New("okx: staking-defi purchase requires productId")
	errFinanceStakingDefiPurchaseMissingInvestData = errors.New("okx: staking-defi purchase requires investData")
	errFinanceStakingDefiPurchaseInvalidInvestData = errors.New("okx: staking-defi purchase requires investData[].ccy and investData[].amt")
	errEmptyFinanceStakingDefiPurchaseAck          = errors.New("okx: empty staking-defi purchase response")
)

// Do 申购链上赚币项目（POST /api/v5/finance/staking-defi/purchase）。
func (s *FinanceStakingDefiPurchaseService) Do(ctx context.Context) (*FinanceStakingDefiOrderAck, error) {
	if s.req.ProductId == "" {
		return nil, errFinanceStakingDefiPurchaseMissingProductId
	}
	if len(s.req.InvestData) == 0 {
		return nil, errFinanceStakingDefiPurchaseMissingInvestData
	}
	for _, it := range s.req.InvestData {
		if it.Ccy == "" || it.Amt == "" {
			return nil, errFinanceStakingDefiPurchaseInvalidInvestData
		}
	}

	var data []FinanceStakingDefiOrderAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/finance/staking-defi/purchase", nil, s.req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyFinanceStakingDefiPurchaseAck
	}
	return &data[0], nil
}
