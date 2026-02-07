package okx

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type accountBillsHistoryArchiveApplyRequest struct {
	Year    string `json:"year"`
	Quarter string `json:"quarter"`
}

// AccountBillsHistoryArchiveApplyAck 表示申请账单流水（自 2021 年）返回项。
type AccountBillsHistoryArchiveApplyAck struct {
	Result string    `json:"result"`
	TS     UnixMilli `json:"ts"`
}

// AccountBillsHistoryArchiveApplyService 申请账单流水（自 2021 年，不包括当前季度）。
type AccountBillsHistoryArchiveApplyService struct {
	c   *Client
	req accountBillsHistoryArchiveApplyRequest
}

// NewAccountBillsHistoryArchiveApplyService 创建 AccountBillsHistoryArchiveApplyService。
func (c *Client) NewAccountBillsHistoryArchiveApplyService() *AccountBillsHistoryArchiveApplyService {
	return &AccountBillsHistoryArchiveApplyService{c: c}
}

// Year 设置年份（必填，4 位数字，如 2023）。
func (s *AccountBillsHistoryArchiveApplyService) Year(year string) *AccountBillsHistoryArchiveApplyService {
	s.req.Year = year
	return s
}

// Quarter 设置季度（必填：Q1/Q2/Q3/Q4）。
func (s *AccountBillsHistoryArchiveApplyService) Quarter(quarter string) *AccountBillsHistoryArchiveApplyService {
	s.req.Quarter = quarter
	return s
}

var (
	errAccountBillsHistoryArchiveApplyMissingRequired = errors.New("okx: bills history archive apply requires year and quarter")
	errEmptyAccountBillsHistoryArchiveApply           = errors.New("okx: empty bills history archive apply response")
)

// Do 申请账单流水（POST /api/v5/account/bills-history-archive）。
func (s *AccountBillsHistoryArchiveApplyService) Do(ctx context.Context) (*AccountBillsHistoryArchiveApplyAck, error) {
	if s.req.Year == "" || s.req.Quarter == "" {
		return nil, errAccountBillsHistoryArchiveApplyMissingRequired
	}

	var data []AccountBillsHistoryArchiveApplyAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/account/bills-history-archive", nil, s.req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyAccountBillsHistoryArchiveApply
	}
	if !strings.EqualFold(data[0].Result, "true") {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodPost,
			RequestPath: "/api/v5/account/bills-history-archive",
			RequestID:   requestID,
			Code:        "0",
			Message:     "bills history archive apply result is false",
		}
	}
	return &data[0], nil
}
