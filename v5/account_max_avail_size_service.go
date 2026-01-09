package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// AccountMaxAvailSize 表示最大可用余额/保证金。
type AccountMaxAvailSize struct {
	InstId    string `json:"instId"`
	AvailBuy  string `json:"availBuy"`
	AvailSell string `json:"availSell"`
}

// AccountMaxAvailSizeService 获取最大可用余额/保证金。
type AccountMaxAvailSizeService struct {
	c *Client

	instId        string
	tdMode        string
	ccy           string
	reduceOnly    *bool
	px            string
	tradeQuoteCcy string
}

// NewAccountMaxAvailSizeService 创建 AccountMaxAvailSizeService。
func (c *Client) NewAccountMaxAvailSizeService() *AccountMaxAvailSizeService {
	return &AccountMaxAvailSizeService{c: c}
}

// InstId 设置产品 ID（支持多产品 ID 查询，逗号分隔；最多 5 个）。
func (s *AccountMaxAvailSizeService) InstId(instId string) *AccountMaxAvailSizeService {
	s.instId = instId
	return s
}

// TdMode 设置交易模式（必填：cross/isolated/cash/spot_isolated）。
func (s *AccountMaxAvailSizeService) TdMode(tdMode string) *AccountMaxAvailSizeService {
	s.tdMode = tdMode
	return s
}

// Ccy 设置保证金币种（可选）。
func (s *AccountMaxAvailSizeService) Ccy(ccy string) *AccountMaxAvailSizeService {
	s.ccy = ccy
	return s
}

// ReduceOnly 设置只减仓模式（可选，仅适用于币币杠杆）。
func (s *AccountMaxAvailSizeService) ReduceOnly(reduceOnly bool) *AccountMaxAvailSizeService {
	s.reduceOnly = &reduceOnly
	return s
}

// Px 设置平仓价格（可选，仅适用于杠杆只减仓；默认市价）。
func (s *AccountMaxAvailSizeService) Px(px string) *AccountMaxAvailSizeService {
	s.px = px
	return s
}

// TradeQuoteCcy 设置用于交易的计价币种（可选，仅适用于币币）。
func (s *AccountMaxAvailSizeService) TradeQuoteCcy(tradeQuoteCcy string) *AccountMaxAvailSizeService {
	s.tradeQuoteCcy = tradeQuoteCcy
	return s
}

var errAccountMaxAvailSizeMissingRequired = errors.New("okx: max avail size requires instId and tdMode")

// Do 获取最大可用余额/保证金（GET /api/v5/account/max-avail-size）。
func (s *AccountMaxAvailSizeService) Do(ctx context.Context) ([]AccountMaxAvailSize, error) {
	if s.instId == "" || s.tdMode == "" {
		return nil, errAccountMaxAvailSizeMissingRequired
	}

	q := url.Values{}
	q.Set("instId", s.instId)
	q.Set("tdMode", s.tdMode)
	if s.ccy != "" {
		q.Set("ccy", s.ccy)
	}
	if s.reduceOnly != nil {
		q.Set("reduceOnly", strconv.FormatBool(*s.reduceOnly))
	}
	if s.px != "" {
		q.Set("px", s.px)
	}
	if s.tradeQuoteCcy != "" {
		q.Set("tradeQuoteCcy", s.tradeQuoteCcy)
	}

	var data []AccountMaxAvailSize
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/max-avail-size", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
