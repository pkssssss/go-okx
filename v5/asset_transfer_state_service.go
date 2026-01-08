package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// AssetTransferState 表示资金划转状态。
type AssetTransferState struct {
	TransId  string `json:"transId"`
	ClientId string `json:"clientId"`
	Ccy      string `json:"ccy"`
	Amt      string `json:"amt"`
	Type     string `json:"type"`
	From     string `json:"from"`
	To       string `json:"to"`
	SubAcct  string `json:"subAcct"`

	InstId   string `json:"instId"`   // 已废弃
	ToInstId string `json:"toInstId"` // 已废弃

	State string `json:"state"`
}

// AssetTransferStateService 查询资金划转状态（近两周）。
type AssetTransferStateService struct {
	c        *Client
	transId  string
	clientId string
	typ      string
}

// NewAssetTransferStateService 创建 AssetTransferStateService。
func (c *Client) NewAssetTransferStateService() *AssetTransferStateService {
	return &AssetTransferStateService{c: c}
}

// TransId 设置划转 ID（transId 和 clientId 至少填一个；若都填以 transId 为主）。
func (s *AssetTransferStateService) TransId(transId string) *AssetTransferStateService {
	s.transId = transId
	return s
}

// ClientId 设置客户自定义 ID（transId 和 clientId 至少填一个）。
func (s *AssetTransferStateService) ClientId(clientId string) *AssetTransferStateService {
	s.clientId = clientId
	return s
}

// Type 设置划转类型（默认 0）。
func (s *AssetTransferStateService) Type(typ string) *AssetTransferStateService {
	s.typ = typ
	return s
}

var errAssetTransferStateMissingID = errors.New("okx: asset transfer state requires transId or clientId")

// Do 查询资金划转状态（GET /api/v5/asset/transfer-state）。
func (s *AssetTransferStateService) Do(ctx context.Context) ([]AssetTransferState, error) {
	if s.transId == "" && s.clientId == "" {
		return nil, errAssetTransferStateMissingID
	}

	q := url.Values{}
	if s.transId != "" {
		q.Set("transId", s.transId)
	} else {
		q.Set("clientId", s.clientId)
	}
	if s.typ != "" {
		q.Set("type", s.typ)
	}

	var data []AssetTransferState
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/asset/transfer-state", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
