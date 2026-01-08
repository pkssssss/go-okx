package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// AssetMonthlyStatement 表示月结单状态与下载链接。
type AssetMonthlyStatement struct {
	FileHref string    `json:"fileHref"`
	State    string    `json:"state"`
	TS       UnixMilli `json:"ts"`
}

// AssetMonthlyStatementService 获取月结单（近一年）。
type AssetMonthlyStatementService struct {
	c     *Client
	month string
}

// NewAssetMonthlyStatementService 创建 AssetMonthlyStatementService。
func (c *Client) NewAssetMonthlyStatementService() *AssetMonthlyStatementService {
	return &AssetMonthlyStatementService{c: c}
}

// Month 设置月份（必填：Jan/Feb/.../Dec）。
func (s *AssetMonthlyStatementService) Month(month string) *AssetMonthlyStatementService {
	s.month = month
	return s
}

var errAssetMonthlyStatementMissingMonth = errors.New("okx: monthly statement requires month")

// Do 获取月结单（GET /api/v5/asset/monthly-statement）。
func (s *AssetMonthlyStatementService) Do(ctx context.Context) ([]AssetMonthlyStatement, error) {
	if s.month == "" {
		return nil, errAssetMonthlyStatementMissingMonth
	}

	q := url.Values{}
	q.Set("month", s.month)

	var data []AssetMonthlyStatement
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/asset/monthly-statement", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
