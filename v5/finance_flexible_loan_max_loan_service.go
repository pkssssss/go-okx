package okx

import (
	"context"
	"errors"
	"net/http"
)

type financeFlexibleLoanMaxLoanRequest struct {
	BorrowCcy     string                             `json:"borrowCcy"`
	SupCollateral []FinanceFlexibleLoanSupCollateral `json:"supCollateral,omitempty"`
}

// FinanceFlexibleLoanMaxLoanService 获取最大可借。
type FinanceFlexibleLoanMaxLoanService struct {
	c   *Client
	req financeFlexibleLoanMaxLoanRequest
}

// NewFinanceFlexibleLoanMaxLoanService 创建 FinanceFlexibleLoanMaxLoanService。
func (c *Client) NewFinanceFlexibleLoanMaxLoanService() *FinanceFlexibleLoanMaxLoanService {
	return &FinanceFlexibleLoanMaxLoanService{c: c}
}

// BorrowCcy 设置借币币种（必填）。
func (s *FinanceFlexibleLoanMaxLoanService) BorrowCcy(borrowCcy string) *FinanceFlexibleLoanMaxLoanService {
	s.req.BorrowCcy = borrowCcy
	return s
}

// SupCollateral 设置补充抵押资产信息（可选）。
func (s *FinanceFlexibleLoanMaxLoanService) SupCollateral(sup []FinanceFlexibleLoanSupCollateral) *FinanceFlexibleLoanMaxLoanService {
	s.req.SupCollateral = sup
	return s
}

var (
	errFinanceFlexibleLoanMaxLoanMissingBorrowCcy     = errors.New("okx: flexible loan max loan requires borrowCcy")
	errFinanceFlexibleLoanMaxLoanInvalidSupCollateral = errors.New("okx: flexible loan max loan requires supCollateral[].ccy and supCollateral[].amt")
	errEmptyFinanceFlexibleLoanMaxLoan                = errors.New("okx: empty flexible loan max loan response")
)

// Do 获取最大可借（POST /api/v5/finance/flexible-loan/max-loan）。
func (s *FinanceFlexibleLoanMaxLoanService) Do(ctx context.Context) (*FinanceFlexibleLoanMaxLoan, error) {
	if s.req.BorrowCcy == "" {
		return nil, errFinanceFlexibleLoanMaxLoanMissingBorrowCcy
	}
	for _, it := range s.req.SupCollateral {
		if it.Ccy == "" || it.Amt == "" {
			return nil, errFinanceFlexibleLoanMaxLoanInvalidSupCollateral
		}
	}

	var data []FinanceFlexibleLoanMaxLoan
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/finance/flexible-loan/max-loan", nil, s.req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/finance/flexible-loan/max-loan", requestID, errEmptyFinanceFlexibleLoanMaxLoan)
	}
	return &data[0], nil
}
