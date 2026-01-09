package okx

import (
	"context"
	"errors"
	"net/http"
)

type accountSpotManualBorrowRepayRequest struct {
	Ccy  string `json:"ccy"`
	Side string `json:"side"`
	Amt  string `json:"amt"`
}

// AccountSpotManualBorrowRepayAck 表示手动借/还币返回项。
type AccountSpotManualBorrowRepayAck struct {
	Ccy  string `json:"ccy"`
	Side string `json:"side"`
	Amt  string `json:"amt"`
}

// AccountSpotManualBorrowRepayService 现货手动借/还币。
type AccountSpotManualBorrowRepayService struct {
	c   *Client
	req accountSpotManualBorrowRepayRequest
}

// NewAccountSpotManualBorrowRepayService 创建 AccountSpotManualBorrowRepayService。
func (c *Client) NewAccountSpotManualBorrowRepayService() *AccountSpotManualBorrowRepayService {
	return &AccountSpotManualBorrowRepayService{c: c}
}

// Ccy 设置币种（必填）。
func (s *AccountSpotManualBorrowRepayService) Ccy(ccy string) *AccountSpotManualBorrowRepayService {
	s.req.Ccy = ccy
	return s
}

// Side 设置方向（必填：borrow/repay）。
func (s *AccountSpotManualBorrowRepayService) Side(side string) *AccountSpotManualBorrowRepayService {
	s.req.Side = side
	return s
}

// Amt 设置数量（必填）。
func (s *AccountSpotManualBorrowRepayService) Amt(amt string) *AccountSpotManualBorrowRepayService {
	s.req.Amt = amt
	return s
}

var (
	errAccountSpotManualBorrowRepayMissingRequired = errors.New("okx: spot manual borrow repay requires ccy/side/amt")
	errEmptyAccountSpotManualBorrowRepay           = errors.New("okx: empty spot manual borrow repay response")
)

// Do 手动借/还币（POST /api/v5/account/spot-manual-borrow-repay）。
func (s *AccountSpotManualBorrowRepayService) Do(ctx context.Context) (*AccountSpotManualBorrowRepayAck, error) {
	if s.req.Ccy == "" || s.req.Side == "" || s.req.Amt == "" {
		return nil, errAccountSpotManualBorrowRepayMissingRequired
	}

	var data []AccountSpotManualBorrowRepayAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/account/spot-manual-borrow-repay", nil, s.req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyAccountSpotManualBorrowRepay
	}
	return &data[0], nil
}
