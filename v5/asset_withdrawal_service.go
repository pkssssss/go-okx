package okx

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type assetWithdrawalRequest struct {
	Ccy        string          `json:"ccy"`
	Amt        string          `json:"amt"`
	Dest       string          `json:"dest"`
	ToAddr     string          `json:"toAddr"`
	ToAddrType string          `json:"toAddrType,omitempty"`
	Chain      string          `json:"chain,omitempty"`
	AreaCode   string          `json:"areaCode,omitempty"`
	RcvrInfo   json.RawMessage `json:"rcvrInfo,omitempty"`
	ClientId   string          `json:"clientId,omitempty"`
}

// AssetWithdrawalReceiverInfo 表示接收方信息（特定主体用户链上提币/闪电网络提币）。
type AssetWithdrawalReceiverInfo struct {
	WalletType string `json:"walletType"`

	ExchId string `json:"exchId,omitempty"`

	RcvrFirstName string `json:"rcvrFirstName,omitempty"`
	RcvrLastName  string `json:"rcvrLastName,omitempty"`

	RcvrCountry            string `json:"rcvrCountry,omitempty"`
	RcvrCountrySubDivision string `json:"rcvrCountrySubDivision,omitempty"`
	RcvrTownName           string `json:"rcvrTownName,omitempty"`
	RcvrStreetName         string `json:"rcvrStreetName,omitempty"`
}

// AssetWithdrawalAck 表示提币返回项。
type AssetWithdrawalAck struct {
	Ccy      string `json:"ccy"`
	Chain    string `json:"chain"`
	Amt      string `json:"amt"`
	WdId     string `json:"wdId"`
	ClientId string `json:"clientId"`
}

// AssetWithdrawalService 提币（内部转账 / 链上提币）。
type AssetWithdrawalService struct {
	c   *Client
	req assetWithdrawalRequest
}

// NewAssetWithdrawalService 创建 AssetWithdrawalService。
func (c *Client) NewAssetWithdrawalService() *AssetWithdrawalService {
	return &AssetWithdrawalService{c: c}
}

// Ccy 设置币种（必填）。
func (s *AssetWithdrawalService) Ccy(ccy string) *AssetWithdrawalService {
	s.req.Ccy = ccy
	return s
}

// Amt 设置提币数量（必填，不包含手续费）。
func (s *AssetWithdrawalService) Amt(amt string) *AssetWithdrawalService {
	s.req.Amt = amt
	return s
}

// Dest 设置提币方式（必填：3=内部转账，4=链上提币）。
func (s *AssetWithdrawalService) Dest(dest string) *AssetWithdrawalService {
	s.req.Dest = dest
	return s
}

// ToAddr 设置提币地址/账户（必填）。
func (s *AssetWithdrawalService) ToAddr(toAddr string) *AssetWithdrawalService {
	s.req.ToAddr = toAddr
	return s
}

// ToAddrType 设置地址类型（1=钱包地址/邮箱/手机号/登录账户名，2=UID）。
func (s *AssetWithdrawalService) ToAddrType(toAddrType string) *AssetWithdrawalService {
	s.req.ToAddrType = toAddrType
	return s
}

// Chain 设置币种链信息（链上提币可选）。
func (s *AssetWithdrawalService) Chain(chain string) *AssetWithdrawalService {
	s.req.Chain = chain
	return s
}

// AreaCode 设置手机区号（内部转账且 toAddr 为手机号时必填）。
func (s *AssetWithdrawalService) AreaCode(areaCode string) *AssetWithdrawalService {
	s.req.AreaCode = areaCode
	return s
}

// RcvrInfoJSON 设置接收方信息（特定主体用户链上提币需要）。
// 传入 JSON object 的 bytes；为空/nil 表示不设置。
func (s *AssetWithdrawalService) RcvrInfoJSON(rcvrInfo json.RawMessage) *AssetWithdrawalService {
	if len(rcvrInfo) == 0 {
		s.req.RcvrInfo = nil
		return s
	}
	s.req.RcvrInfo = rcvrInfo
	return s
}

// RcvrInfo 设置接收方信息（特定主体用户链上提币需要）。
// 传入 nil 表示清空。
func (s *AssetWithdrawalService) RcvrInfo(info *AssetWithdrawalReceiverInfo) *AssetWithdrawalService {
	if info == nil {
		s.req.RcvrInfo = nil
		return s
	}
	b, err := json.Marshal(info)
	if err != nil {
		s.req.RcvrInfo = nil
		return s
	}
	s.req.RcvrInfo = b
	return s
}

// ClientId 设置客户自定义 ID（1-32）。
func (s *AssetWithdrawalService) ClientId(clientId string) *AssetWithdrawalService {
	s.req.ClientId = clientId
	return s
}

var (
	errAssetWithdrawalMissingRequired = errors.New("okx: withdrawal requires ccy/amt/dest/toAddr")
	errAssetWithdrawalInvalidRcvrInfo = errors.New("okx: invalid rcvrInfo json")
)

// Do 发起提币（POST /api/v5/asset/withdrawal）。
func (s *AssetWithdrawalService) Do(ctx context.Context) (*AssetWithdrawalAck, error) {
	if s.req.Ccy == "" || s.req.Amt == "" || s.req.Dest == "" || s.req.ToAddr == "" {
		return nil, errAssetWithdrawalMissingRequired
	}
	if len(s.req.RcvrInfo) > 0 && !json.Valid(s.req.RcvrInfo) {
		return nil, errAssetWithdrawalInvalidRcvrInfo
	}

	var data []AssetWithdrawalAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/asset/withdrawal", nil, s.req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/asset/withdrawal", requestID, errors.New("okx: empty withdrawal response"))
	}
	return &data[0], nil
}
