package okx

import (
	"context"
	"errors"
	"net/http"
)

type assetTransferRequest struct {
	Type string `json:"type,omitempty"`
	Ccy  string `json:"ccy"`
	Amt  string `json:"amt"`
	From string `json:"from"`
	To   string `json:"to"`

	SubAcct     string `json:"subAcct,omitempty"`
	OmitPosRisk *bool  `json:"omitPosRisk,omitempty"`
	ClientId    string `json:"clientId,omitempty"`
}

// AssetTransferAck 表示资金划转返回项。
type AssetTransferAck struct {
	TransId  string `json:"transId"`
	Ccy      string `json:"ccy"`
	ClientId string `json:"clientId"`
	From     string `json:"from"`
	Amt      string `json:"amt"`
	To       string `json:"to"`
}

// AssetTransferService 资金划转（资金账户/交易账户/子账户）。
type AssetTransferService struct {
	c   *Client
	req assetTransferRequest
}

// NewAssetTransferService 创建 AssetTransferService。
func (c *Client) NewAssetTransferService() *AssetTransferService {
	return &AssetTransferService{c: c}
}

// Type 设置划转类型（0=账户内划转，1/2/3/4=母子账户相关）。
func (s *AssetTransferService) Type(typ string) *AssetTransferService {
	s.req.Type = typ
	return s
}

// Ccy 设置划转币种（必填）。
func (s *AssetTransferService) Ccy(ccy string) *AssetTransferService {
	s.req.Ccy = ccy
	return s
}

// Amt 设置划转数量（必填）。
func (s *AssetTransferService) Amt(amt string) *AssetTransferService {
	s.req.Amt = amt
	return s
}

// From 设置转出账户（必填），如 6=资金账户，18=交易账户。
func (s *AssetTransferService) From(from string) *AssetTransferService {
	s.req.From = from
	return s
}

// To 设置转入账户（必填），如 6=资金账户，18=交易账户。
func (s *AssetTransferService) To(to string) *AssetTransferService {
	s.req.To = to
	return s
}

// SubAcct 设置子账户名称（母子账户划转时使用）。
func (s *AssetTransferService) SubAcct(subAcct string) *AssetTransferService {
	s.req.SubAcct = subAcct
	return s
}

// OmitPosRisk 设置是否忽略仓位风险（仅组合保证金模式适用）。
func (s *AssetTransferService) OmitPosRisk(enable bool) *AssetTransferService {
	s.req.OmitPosRisk = &enable
	return s
}

// ClientId 设置客户自定义 ID（1-32）。
func (s *AssetTransferService) ClientId(clientId string) *AssetTransferService {
	s.req.ClientId = clientId
	return s
}

var errAssetTransferMissingRequired = errors.New("okx: asset transfer requires ccy/amt/from/to")

// Do 发起资金划转（POST /api/v5/asset/transfer）。
//
// 注意：请求成功/失败不一定反映最终划转结果，建议结合 AssetTransferStateService 查询最终状态。
func (s *AssetTransferService) Do(ctx context.Context) (*AssetTransferAck, error) {
	if s.req.Ccy == "" || s.req.Amt == "" || s.req.From == "" || s.req.To == "" {
		return nil, errAssetTransferMissingRequired
	}

	var data []AssetTransferAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/asset/transfer", nil, s.req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/asset/transfer", requestID, errors.New("okx: empty asset transfer response"))
	}
	return &data[0], nil
}
