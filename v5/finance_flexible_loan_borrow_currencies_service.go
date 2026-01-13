package okx

import (
	"context"
	"net/http"
)

// FinanceFlexibleLoanBorrowCurrenciesService 获取可借币种列表。
type FinanceFlexibleLoanBorrowCurrenciesService struct {
	c *Client
}

// NewFinanceFlexibleLoanBorrowCurrenciesService 创建 FinanceFlexibleLoanBorrowCurrenciesService。
func (c *Client) NewFinanceFlexibleLoanBorrowCurrenciesService() *FinanceFlexibleLoanBorrowCurrenciesService {
	return &FinanceFlexibleLoanBorrowCurrenciesService{c: c}
}

// Do 获取可借币种列表（GET /api/v5/finance/flexible-loan/borrow-currencies）。
func (s *FinanceFlexibleLoanBorrowCurrenciesService) Do(ctx context.Context) ([]FinanceFlexibleLoanBorrowCurrency, error) {
	var data []FinanceFlexibleLoanBorrowCurrency
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/finance/flexible-loan/borrow-currencies", nil, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
