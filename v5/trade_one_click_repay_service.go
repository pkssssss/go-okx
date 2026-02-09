package okx

import (
	"context"
	"errors"
	"net/http"
)

// OneClickRepayService 交易一键偿还债务（跨币种保证金/组合保证金）。
type OneClickRepayService struct {
	c *Client

	debtCcy  []string
	repayCcy string
}

// NewOneClickRepayService 创建 OneClickRepayService。
func (c *Client) NewOneClickRepayService() *OneClickRepayService {
	return &OneClickRepayService{c: c}
}

// DebtCcy 设置负债币种列表（必填；单次最多 5 个）。
func (s *OneClickRepayService) DebtCcy(debtCcy []string) *OneClickRepayService {
	s.debtCcy = debtCcy
	return s
}

// RepayCcy 设置偿还币种（必填；不能与 debtCcy 重复）。
func (s *OneClickRepayService) RepayCcy(repayCcy string) *OneClickRepayService {
	s.repayCcy = repayCcy
	return s
}

var (
	errOneClickRepayMissingRequired = errors.New("okx: one-click repay requires debtCcy/repayCcy")
	errOneClickRepayTooManyDebtCcy  = errors.New("okx: one-click repay supports up to 5 debtCcy")
	errOneClickRepayEmptyDebtCcy    = errors.New("okx: one-click repay requires non-empty debtCcy items")
	errOneClickRepaySameCurrency    = errors.New("okx: one-click repay requires repayCcy not in debtCcy")
	errEmptyOneClickRepayResponse   = errors.New("okx: empty one-click repay response")
	errInvalidOneClickRepayResponse = errors.New("okx: invalid one-click repay response")
)

type oneClickRepayRequest struct {
	DebtCcy  []string `json:"debtCcy"`
	RepayCcy string   `json:"repayCcy"`
}

func validateOneClickRepayAck(ack *OneClickRepayAck, req oneClickRepayRequest) error {
	if ack == nil || ack.DebtCcy == "" || ack.RepayCcy == "" || ack.Status == "" {
		return errInvalidOneClickRepayResponse
	}
	if ack.RepayCcy != req.RepayCcy {
		return errInvalidOneClickRepayResponse
	}

	allowedDebtCcy := make(map[string]struct{}, len(req.DebtCcy))
	for _, ccy := range req.DebtCcy {
		allowedDebtCcy[ccy] = struct{}{}
	}
	if _, ok := allowedDebtCcy[ack.DebtCcy]; !ok {
		return errInvalidOneClickRepayResponse
	}
	return nil
}

// Do 一键还债交易（跨币种保证金/组合保证金）（POST /api/v5/trade/one-click-repay）。
func (s *OneClickRepayService) Do(ctx context.Context) ([]OneClickRepayAck, error) {
	if len(s.debtCcy) == 0 || s.repayCcy == "" {
		return nil, errOneClickRepayMissingRequired
	}
	if len(s.debtCcy) > 5 {
		return nil, errOneClickRepayTooManyDebtCcy
	}
	for _, ccy := range s.debtCcy {
		if ccy == "" {
			return nil, errOneClickRepayEmptyDebtCcy
		}
		if ccy == s.repayCcy {
			return nil, errOneClickRepaySameCurrency
		}
	}

	req := oneClickRepayRequest{
		DebtCcy:  s.debtCcy,
		RepayCcy: s.repayCcy,
	}

	var data []OneClickRepayAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/trade/one-click-repay", nil, req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyOneClickRepayResponse
	}
	for i := range data {
		if err := validateOneClickRepayAck(&data[i], req); err != nil {
			return nil, err
		}
	}
	return data, nil
}
