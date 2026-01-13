package okx

import (
	"context"
	"errors"
	"net/http"
)

// FinanceFlexibleLoanLoanInfoService 获取借贷信息。
type FinanceFlexibleLoanLoanInfoService struct {
	c *Client
}

// NewFinanceFlexibleLoanLoanInfoService 创建 FinanceFlexibleLoanLoanInfoService。
func (c *Client) NewFinanceFlexibleLoanLoanInfoService() *FinanceFlexibleLoanLoanInfoService {
	return &FinanceFlexibleLoanLoanInfoService{c: c}
}

var errEmptyFinanceFlexibleLoanLoanInfo = errors.New("okx: empty flexible loan loan-info response")

// Do 获取借贷信息（GET /api/v5/finance/flexible-loan/loan-info）。
func (s *FinanceFlexibleLoanLoanInfoService) Do(ctx context.Context) (*FinanceFlexibleLoanInfo, error) {
	var data []FinanceFlexibleLoanInfo
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/finance/flexible-loan/loan-info", nil, nil, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyFinanceFlexibleLoanLoanInfo
	}
	return &data[0], nil
}
