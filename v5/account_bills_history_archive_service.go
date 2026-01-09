package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// AccountBillsHistoryArchive 表示获取账单流水（自 2021 年）返回项。
type AccountBillsHistoryArchive struct {
	FileHref string    `json:"fileHref"`
	State    string    `json:"state"`
	TS       UnixMilli `json:"ts"`
}

// AccountBillsHistoryArchiveService 获取账单流水（自 2021 年）。
type AccountBillsHistoryArchiveService struct {
	c *Client

	year    string
	quarter string
}

// NewAccountBillsHistoryArchiveService 创建 AccountBillsHistoryArchiveService。
func (c *Client) NewAccountBillsHistoryArchiveService() *AccountBillsHistoryArchiveService {
	return &AccountBillsHistoryArchiveService{c: c}
}

// Year 设置年份（必填，4 位数字，如 2023）。
func (s *AccountBillsHistoryArchiveService) Year(year string) *AccountBillsHistoryArchiveService {
	s.year = year
	return s
}

// Quarter 设置季度（必填：Q1/Q2/Q3/Q4）。
func (s *AccountBillsHistoryArchiveService) Quarter(quarter string) *AccountBillsHistoryArchiveService {
	s.quarter = quarter
	return s
}

var (
	errAccountBillsHistoryArchiveMissingRequired = errors.New("okx: bills history archive requires year and quarter")
	errEmptyAccountBillsHistoryArchive           = errors.New("okx: empty bills history archive response")
)

// Do 获取账单流水（GET /api/v5/account/bills-history-archive）。
func (s *AccountBillsHistoryArchiveService) Do(ctx context.Context) (*AccountBillsHistoryArchive, error) {
	if s.year == "" || s.quarter == "" {
		return nil, errAccountBillsHistoryArchiveMissingRequired
	}

	q := url.Values{}
	q.Set("year", s.year)
	q.Set("quarter", s.quarter)

	var data []AccountBillsHistoryArchive
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/bills-history-archive", q, nil, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyAccountBillsHistoryArchive
	}
	return &data[0], nil
}
