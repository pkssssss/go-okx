package okx

import (
	"context"
	"errors"
	"net/http"
)

// OneClickRepayV2Service 交易一键偿还债务（新）。
type OneClickRepayV2Service struct {
	c *Client

	debtCcy      string
	repayCcyList []string
}

// NewOneClickRepayV2Service 创建 OneClickRepayV2Service。
func (c *Client) NewOneClickRepayV2Service() *OneClickRepayV2Service {
	return &OneClickRepayV2Service{c: c}
}

// DebtCcy 设置负债币种（必填）。
func (s *OneClickRepayV2Service) DebtCcy(debtCcy string) *OneClickRepayV2Service {
	s.debtCcy = debtCcy
	return s
}

// RepayCcyList 设置偿还币种列表（必填，排序代表偿还优先级，排第一的优先级最高）。
func (s *OneClickRepayV2Service) RepayCcyList(repayCcyList []string) *OneClickRepayV2Service {
	s.repayCcyList = repayCcyList
	return s
}

var (
	errOneClickRepayV2MissingRequired = errors.New("okx: one-click repay v2 requires debtCcy/repayCcyList")
	errOneClickRepayV2InvalidCcy      = errors.New("okx: one-click repay v2 requires non-empty repayCcyList items")
	errEmptyOneClickRepayV2Response   = errors.New("okx: empty one-click repay v2 response")
	errInvalidOneClickRepayV2Response = errors.New("okx: invalid one-click repay v2 response")
)

type oneClickRepayV2Request struct {
	DebtCcy      string   `json:"debtCcy"`
	RepayCcyList []string `json:"repayCcyList"`
}

func validateOneClickRepayV2Ack(ack *OneClickRepayV2Ack, req oneClickRepayV2Request) error {
	if ack == nil || ack.DebtCcy == "" || len(ack.RepayCcyList) == 0 || ack.TS <= 0 {
		return errInvalidOneClickRepayV2Response
	}
	if ack.DebtCcy != req.DebtCcy || len(ack.RepayCcyList) != len(req.RepayCcyList) {
		return errInvalidOneClickRepayV2Response
	}

	repayCcyCount := make(map[string]int, len(req.RepayCcyList))
	for _, ccy := range req.RepayCcyList {
		repayCcyCount[ccy]++
	}
	for _, ccy := range ack.RepayCcyList {
		if ccy == "" {
			return errInvalidOneClickRepayV2Response
		}
		count := repayCcyCount[ccy]
		if count == 0 {
			return errInvalidOneClickRepayV2Response
		}
		repayCcyCount[ccy] = count - 1
	}

	for _, count := range repayCcyCount {
		if count != 0 {
			return errInvalidOneClickRepayV2Response
		}
	}
	return nil
}

// Do 一键还债交易（新）（POST /api/v5/trade/one-click-repay-v2）。
func (s *OneClickRepayV2Service) Do(ctx context.Context) (*OneClickRepayV2Ack, error) {
	if s.debtCcy == "" || len(s.repayCcyList) == 0 {
		return nil, errOneClickRepayV2MissingRequired
	}
	for _, ccy := range s.repayCcyList {
		if ccy == "" {
			return nil, errOneClickRepayV2InvalidCcy
		}
	}

	req := oneClickRepayV2Request{
		DebtCcy:      s.debtCcy,
		RepayCcyList: s.repayCcyList,
	}

	var data []OneClickRepayV2Ack
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/trade/one-click-repay-v2", nil, req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/trade/one-click-repay-v2", requestID, errEmptyOneClickRepayV2Response)
	}
	if err := validateOneClickRepayV2Ack(&data[0], req); err != nil {
		return nil, newInvalidDataAPIError(http.MethodPost, "/api/v5/trade/one-click-repay-v2", requestID, err)
	}
	return &data[0], nil
}
