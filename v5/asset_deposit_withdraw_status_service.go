package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// AssetDepositWithdrawStatus 表示充值/提现的详细状态与预估完成时间。
type AssetDepositWithdrawStatus struct {
	EstCompleteTime string `json:"estCompleteTime"`
	State           string `json:"state"`
	TxId            string `json:"txId"`
	WdId            string `json:"wdId"`
}

// AssetDepositWithdrawStatusService 获取充值/提现的详细状态（含预估完成时间）。
type AssetDepositWithdrawStatusService struct {
	c *Client

	wdId string
	txId string

	ccy   string
	to    string
	chain string
}

// NewAssetDepositWithdrawStatusService 创建 AssetDepositWithdrawStatusService。
func (c *Client) NewAssetDepositWithdrawStatusService() *AssetDepositWithdrawStatusService {
	return &AssetDepositWithdrawStatusService{c: c}
}

// WdId 设置提币申请 ID（用于查询提现状态；wdId 与 txId 必传其一且仅可传其一）。
func (s *AssetDepositWithdrawStatusService) WdId(wdId string) *AssetDepositWithdrawStatusService {
	s.wdId = wdId
	return s
}

// TxId 设置区块转账哈希记录 ID（用于查询充值状态；wdId 与 txId 必传其一且仅可传其一）。
func (s *AssetDepositWithdrawStatusService) TxId(txId string) *AssetDepositWithdrawStatusService {
	s.txId = txId
	return s
}

// Ccy 设置币种（查询充值时必填，需要与 txId 一并提供）。
func (s *AssetDepositWithdrawStatusService) Ccy(ccy string) *AssetDepositWithdrawStatusService {
	s.ccy = ccy
	return s
}

// To 设置资金充值到账账户地址（查询充值时必填，需要与 txId 一并提供）。
func (s *AssetDepositWithdrawStatusService) To(to string) *AssetDepositWithdrawStatusService {
	s.to = to
	return s
}

// Chain 设置币种链信息（查询充值时必填，需要与 txId 一并提供）。
func (s *AssetDepositWithdrawStatusService) Chain(chain string) *AssetDepositWithdrawStatusService {
	s.chain = chain
	return s
}

var (
	errAssetDepositWithdrawStatusMissingID      = errors.New("okx: deposit-withdraw-status requires wdId or txId")
	errAssetDepositWithdrawStatusMultipleID     = errors.New("okx: deposit-withdraw-status requires exactly one of wdId or txId")
	errAssetDepositWithdrawStatusDepositMissing = errors.New("okx: deposit-withdraw-status deposit query requires ccy/to/chain with txId")
)

// Do 获取充值/提现的详细状态（GET /api/v5/asset/deposit-withdraw-status）。
func (s *AssetDepositWithdrawStatusService) Do(ctx context.Context) ([]AssetDepositWithdrawStatus, error) {
	if s.wdId == "" && s.txId == "" {
		return nil, errAssetDepositWithdrawStatusMissingID
	}
	if s.wdId != "" && s.txId != "" {
		return nil, errAssetDepositWithdrawStatusMultipleID
	}

	q := url.Values{}
	if s.wdId != "" {
		q.Set("wdId", s.wdId)
	} else {
		if s.ccy == "" || s.to == "" || s.chain == "" {
			return nil, errAssetDepositWithdrawStatusDepositMissing
		}
		q.Set("txId", s.txId)
		q.Set("ccy", s.ccy)
		q.Set("to", s.to)
		q.Set("chain", s.chain)
	}

	var data []AssetDepositWithdrawStatus
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/asset/deposit-withdraw-status", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
