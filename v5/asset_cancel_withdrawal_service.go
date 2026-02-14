package okx

import (
	"context"
	"errors"
	"net/http"
)

type assetCancelWithdrawalRequest struct {
	WdId string `json:"wdId"`
}

// AssetCancelWithdrawalAck 表示撤销提币返回项。
type AssetCancelWithdrawalAck struct {
	WdId string `json:"wdId"`
}

// AssetCancelWithdrawalService 撤销提币（不支持撤销闪电网络提币）。
type AssetCancelWithdrawalService struct {
	c    *Client
	wdId string
}

// NewAssetCancelWithdrawalService 创建 AssetCancelWithdrawalService。
func (c *Client) NewAssetCancelWithdrawalService() *AssetCancelWithdrawalService {
	return &AssetCancelWithdrawalService{c: c}
}

// WdId 设置提币申请 ID（必填）。
func (s *AssetCancelWithdrawalService) WdId(wdId string) *AssetCancelWithdrawalService {
	s.wdId = wdId
	return s
}

var errAssetCancelWithdrawalMissingWdId = errors.New("okx: cancel withdrawal requires wdId")

// Do 撤销提币（POST /api/v5/asset/cancel-withdrawal）。
func (s *AssetCancelWithdrawalService) Do(ctx context.Context) (*AssetCancelWithdrawalAck, error) {
	if s.wdId == "" {
		return nil, errAssetCancelWithdrawalMissingWdId
	}

	req := assetCancelWithdrawalRequest{WdId: s.wdId}

	var data []AssetCancelWithdrawalAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/asset/cancel-withdrawal", nil, req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/asset/cancel-withdrawal", requestID, errors.New("okx: empty cancel withdrawal response"))
	}
	return &data[0], nil
}
