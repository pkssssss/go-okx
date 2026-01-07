package okx

import (
	"context"
	"net/http"
	"net/url"
)

// AccountPosition 表示持仓信息（字段按 OKX 返回保持为 string，无损）。
// v0.1 仅保留量化常用字段，其他字段后续按需补齐。
type AccountPosition struct {
	InstType string `json:"instType"`
	InstId   string `json:"instId"`
	PosId    string `json:"posId"`
	PosSide  string `json:"posSide"`

	Pos      string `json:"pos"`
	AvailPos string `json:"availPos"`
	AvgPx    string `json:"avgPx"`
	MarkPx   string `json:"markPx"`
	LiqPx    string `json:"liqPx"`

	Upl      string `json:"upl"`
	UplRatio string `json:"uplRatio"`

	Lever   string `json:"lever"`
	MgnMode string `json:"mgnMode"`
	Ccy     string `json:"ccy"`
}

// AccountPositionsService 查看持仓信息。
type AccountPositionsService struct {
	c        *Client
	instType string
	instId   string
	posId    string
}

// NewAccountPositionsService 创建 AccountPositionsService。
func (c *Client) NewAccountPositionsService() *AccountPositionsService {
	return &AccountPositionsService{c: c}
}

// InstType 设置产品类型（MARGIN/SWAP/FUTURES/OPTION）。
func (s *AccountPositionsService) InstType(instType string) *AccountPositionsService {
	s.instType = instType
	return s
}

// InstId 设置交易产品 ID（支持多个 instId，逗号分隔）。
func (s *AccountPositionsService) InstId(instId string) *AccountPositionsService {
	s.instId = instId
	return s
}

// PosId 设置持仓 ID（支持多个 posId，逗号分隔）。
func (s *AccountPositionsService) PosId(posId string) *AccountPositionsService {
	s.posId = posId
	return s
}

// Do 查看持仓信息（GET /api/v5/account/positions）。
func (s *AccountPositionsService) Do(ctx context.Context) ([]AccountPosition, error) {
	var q url.Values
	if s.instType != "" || s.instId != "" || s.posId != "" {
		q = url.Values{}
		if s.instType != "" {
			q.Set("instType", s.instType)
		}
		if s.instId != "" {
			q.Set("instId", s.instId)
		}
		if s.posId != "" {
			q.Set("posId", s.posId)
		}
	}

	var data []AccountPosition
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/positions", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
