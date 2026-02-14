package okx

import (
	"context"
	"errors"
	"net/http"
)

type assetMonthlyStatementApplyRequest struct {
	Month string `json:"month,omitempty"`
}

// AssetMonthlyStatementApplyAck 表示申请月结单返回项。
type AssetMonthlyStatementApplyAck struct {
	TS UnixMilli `json:"ts"`
}

// AssetMonthlyStatementApplyService 申请月结单（近一年）。
type AssetMonthlyStatementApplyService struct {
	c     *Client
	month string
}

// NewAssetMonthlyStatementApplyService 创建 AssetMonthlyStatementApplyService。
func (c *Client) NewAssetMonthlyStatementApplyService() *AssetMonthlyStatementApplyService {
	return &AssetMonthlyStatementApplyService{c: c}
}

// Month 设置月份（默认上一个月：Jan/Feb/.../Dec）。
func (s *AssetMonthlyStatementApplyService) Month(month string) *AssetMonthlyStatementApplyService {
	s.month = month
	return s
}

var errEmptyAssetMonthlyStatementApply = errors.New("okx: empty monthly statement apply response")

// Do 申请月结单（POST /api/v5/asset/monthly-statement）。
func (s *AssetMonthlyStatementApplyService) Do(ctx context.Context) (*AssetMonthlyStatementApplyAck, error) {
	req := assetMonthlyStatementApplyRequest{Month: s.month}

	var data []AssetMonthlyStatementApplyAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/asset/monthly-statement", nil, req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/asset/monthly-statement", requestID, errEmptyAssetMonthlyStatementApply)
	}
	return &data[0], nil
}
