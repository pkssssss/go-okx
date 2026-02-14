package okx

import (
	"context"
	"errors"
	"net/http"
)

type assetSubaccountTransferRequest struct {
	Ccy            string `json:"ccy"`
	Amt            string `json:"amt"`
	From           string `json:"from"`
	To             string `json:"to"`
	FromSubAccount string `json:"fromSubAccount"`
	ToSubAccount   string `json:"toSubAccount"`

	LoanTrans   *bool `json:"loanTrans,omitempty"`
	OmitPosRisk *bool `json:"omitPosRisk,omitempty"`
}

// AssetSubaccountTransferAck 表示子账户间资金划转返回项。
type AssetSubaccountTransferAck struct {
	TransId string `json:"transId"`
}

// AssetSubaccountTransferService 子账户间资金划转（母账户）。
type AssetSubaccountTransferService struct {
	c   *Client
	req assetSubaccountTransferRequest
}

// NewAssetSubaccountTransferService 创建 AssetSubaccountTransferService。
func (c *Client) NewAssetSubaccountTransferService() *AssetSubaccountTransferService {
	return &AssetSubaccountTransferService{c: c}
}

// Ccy 设置币种（必填）。
func (s *AssetSubaccountTransferService) Ccy(ccy string) *AssetSubaccountTransferService {
	s.req.Ccy = ccy
	return s
}

// Amt 设置划转数量（必填）。
func (s *AssetSubaccountTransferService) Amt(amt string) *AssetSubaccountTransferService {
	s.req.Amt = amt
	return s
}

// From 设置转出子账户类型（必填：6=资金账户，18=交易账户）。
func (s *AssetSubaccountTransferService) From(from string) *AssetSubaccountTransferService {
	s.req.From = from
	return s
}

// To 设置转入子账户类型（必填：6=资金账户，18=交易账户）。
func (s *AssetSubaccountTransferService) To(to string) *AssetSubaccountTransferService {
	s.req.To = to
	return s
}

// FromSubAccount 设置转出子账户名称（必填）。
func (s *AssetSubaccountTransferService) FromSubAccount(fromSubAccount string) *AssetSubaccountTransferService {
	s.req.FromSubAccount = fromSubAccount
	return s
}

// ToSubAccount 设置转入子账户名称（必填）。
func (s *AssetSubaccountTransferService) ToSubAccount(toSubAccount string) *AssetSubaccountTransferService {
	s.req.ToSubAccount = toSubAccount
	return s
}

// LoanTrans 设置是否支持跨币种保证金/组合保证金下的借币转入/转出（默认 false）。
func (s *AssetSubaccountTransferService) LoanTrans(enable bool) *AssetSubaccountTransferService {
	s.req.LoanTrans = &enable
	return s
}

// OmitPosRisk 设置是否忽略仓位风险（仅组合保证金模式适用，默认 false）。
func (s *AssetSubaccountTransferService) OmitPosRisk(enable bool) *AssetSubaccountTransferService {
	s.req.OmitPosRisk = &enable
	return s
}

var errAssetSubaccountTransferMissingRequired = errors.New("okx: subaccount transfer requires ccy/amt/from/to/fromSubAccount/toSubAccount")
var errEmptyAssetSubaccountTransfer = errors.New("okx: empty subaccount transfer response")

// Do 子账户间资金划转（POST /api/v5/asset/subaccount/transfer）。
func (s *AssetSubaccountTransferService) Do(ctx context.Context) (*AssetSubaccountTransferAck, error) {
	if s.req.Ccy == "" || s.req.Amt == "" || s.req.From == "" || s.req.To == "" || s.req.FromSubAccount == "" || s.req.ToSubAccount == "" {
		return nil, errAssetSubaccountTransferMissingRequired
	}

	var data []AssetSubaccountTransferAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/asset/subaccount/transfer", nil, s.req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/asset/subaccount/transfer", requestID, errEmptyAssetSubaccountTransfer)
	}
	return &data[0], nil
}
