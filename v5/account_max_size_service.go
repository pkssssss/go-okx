package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// AccountMaxSize 表示最大可下单数量（对应下单时的 sz）。
type AccountMaxSize struct {
	InstId  string `json:"instId"`
	Ccy     string `json:"ccy"`
	MaxBuy  string `json:"maxBuy"`
	MaxSell string `json:"maxSell"`
}

// AccountMaxSizeService 获取最大可下单数量。
type AccountMaxSizeService struct {
	c *Client

	instId        string
	tdMode        string
	ccy           string
	px            string
	leverage      string
	tradeQuoteCcy string
}

// NewAccountMaxSizeService 创建 AccountMaxSizeService。
func (c *Client) NewAccountMaxSizeService() *AccountMaxSizeService {
	return &AccountMaxSizeService{c: c}
}

// InstId 设置产品 ID（支持同一业务线下的多产品 ID，逗号分隔；最多 5 个）。
func (s *AccountMaxSizeService) InstId(instId string) *AccountMaxSizeService {
	s.instId = instId
	return s
}

// TdMode 设置交易模式（必填：cross/isolated/cash/spot_isolated）。
func (s *AccountMaxSizeService) TdMode(tdMode string) *AccountMaxSizeService {
	s.tdMode = tdMode
	return s
}

// Ccy 设置保证金币种（可选）。
func (s *AccountMaxSizeService) Ccy(ccy string) *AccountMaxSizeService {
	s.ccy = ccy
	return s
}

// Px 设置委托价格（可选）。
func (s *AccountMaxSizeService) Px(px string) *AccountMaxSizeService {
	s.px = px
	return s
}

// Leverage 设置开仓杠杆倍数（可选，仅适用于币币杠杆/交割/永续）。
func (s *AccountMaxSizeService) Leverage(leverage string) *AccountMaxSizeService {
	s.leverage = leverage
	return s
}

// TradeQuoteCcy 设置用于交易的计价币种（可选，仅适用于币币）。
func (s *AccountMaxSizeService) TradeQuoteCcy(tradeQuoteCcy string) *AccountMaxSizeService {
	s.tradeQuoteCcy = tradeQuoteCcy
	return s
}

var errAccountMaxSizeMissingRequired = errors.New("okx: max size requires instId and tdMode")

// Do 获取最大可下单数量（GET /api/v5/account/max-size）。
func (s *AccountMaxSizeService) Do(ctx context.Context) ([]AccountMaxSize, error) {
	if s.instId == "" || s.tdMode == "" {
		return nil, errAccountMaxSizeMissingRequired
	}

	q := url.Values{}
	q.Set("instId", s.instId)
	q.Set("tdMode", s.tdMode)
	if s.ccy != "" {
		q.Set("ccy", s.ccy)
	}
	if s.px != "" {
		q.Set("px", s.px)
	}
	if s.leverage != "" {
		q.Set("leverage", s.leverage)
	}
	if s.tradeQuoteCcy != "" {
		q.Set("tradeQuoteCcy", s.tradeQuoteCcy)
	}

	var data []AccountMaxSize
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/account/max-size", q, nil, true, &data); err != nil {
		return nil, err
	}
	return data, nil
}
