package okx

import (
	"context"
	"errors"
	"net/http"
)

type tradingBotGridOrderAlgoRequest struct {
	InstId      string `json:"instId"`
	AlgoOrdType string `json:"algoOrdType"`
	MaxPx       string `json:"maxPx"`
	MinPx       string `json:"minPx"`
	GridNum     string `json:"gridNum"`

	RunType            string                       `json:"runType,omitempty"`
	TpTriggerPx        string                       `json:"tpTriggerPx,omitempty"`
	SlTriggerPx        string                       `json:"slTriggerPx,omitempty"`
	AlgoClOrdId        string                       `json:"algoClOrdId,omitempty"`
	Tag                string                       `json:"tag,omitempty"`
	ProfitSharingRatio string                       `json:"profitSharingRatio,omitempty"`
	TriggerParams      []TradingBotGridTriggerParam `json:"triggerParams,omitempty"`

	QuoteSz       string `json:"quoteSz,omitempty"`
	BaseSz        string `json:"baseSz,omitempty"`
	TradeQuoteCcy string `json:"tradeQuoteCcy,omitempty"`

	Sz        string `json:"sz,omitempty"`
	Direction string `json:"direction,omitempty"`
	Lever     string `json:"lever,omitempty"`
	BasePos   *bool  `json:"basePos,omitempty"`
	TpRatio   string `json:"tpRatio,omitempty"`
	SlRatio   string `json:"slRatio,omitempty"`
}

// TradingBotGridOrderAlgoService 网格策略委托下单。
type TradingBotGridOrderAlgoService struct {
	c *Client
	r tradingBotGridOrderAlgoRequest
}

// NewTradingBotGridOrderAlgoService 创建 TradingBotGridOrderAlgoService。
func (c *Client) NewTradingBotGridOrderAlgoService() *TradingBotGridOrderAlgoService {
	return &TradingBotGridOrderAlgoService{c: c}
}

func (s *TradingBotGridOrderAlgoService) InstId(instId string) *TradingBotGridOrderAlgoService {
	s.r.InstId = instId
	return s
}

func (s *TradingBotGridOrderAlgoService) AlgoOrdType(algoOrdType string) *TradingBotGridOrderAlgoService {
	s.r.AlgoOrdType = algoOrdType
	return s
}

func (s *TradingBotGridOrderAlgoService) MaxPx(maxPx string) *TradingBotGridOrderAlgoService {
	s.r.MaxPx = maxPx
	return s
}

func (s *TradingBotGridOrderAlgoService) MinPx(minPx string) *TradingBotGridOrderAlgoService {
	s.r.MinPx = minPx
	return s
}

func (s *TradingBotGridOrderAlgoService) GridNum(gridNum string) *TradingBotGridOrderAlgoService {
	s.r.GridNum = gridNum
	return s
}

func (s *TradingBotGridOrderAlgoService) RunType(runType string) *TradingBotGridOrderAlgoService {
	s.r.RunType = runType
	return s
}

func (s *TradingBotGridOrderAlgoService) TpTriggerPx(tpTriggerPx string) *TradingBotGridOrderAlgoService {
	s.r.TpTriggerPx = tpTriggerPx
	return s
}

func (s *TradingBotGridOrderAlgoService) SlTriggerPx(slTriggerPx string) *TradingBotGridOrderAlgoService {
	s.r.SlTriggerPx = slTriggerPx
	return s
}

func (s *TradingBotGridOrderAlgoService) AlgoClOrdId(algoClOrdId string) *TradingBotGridOrderAlgoService {
	s.r.AlgoClOrdId = algoClOrdId
	return s
}

func (s *TradingBotGridOrderAlgoService) Tag(tag string) *TradingBotGridOrderAlgoService {
	s.r.Tag = tag
	return s
}

func (s *TradingBotGridOrderAlgoService) ProfitSharingRatio(ratio string) *TradingBotGridOrderAlgoService {
	s.r.ProfitSharingRatio = ratio
	return s
}

func (s *TradingBotGridOrderAlgoService) TriggerParams(params []TradingBotGridTriggerParam) *TradingBotGridOrderAlgoService {
	s.r.TriggerParams = params
	return s
}

// QuoteSz 设置计价币投入数量（现货网格）。
func (s *TradingBotGridOrderAlgoService) QuoteSz(quoteSz string) *TradingBotGridOrderAlgoService {
	s.r.QuoteSz = quoteSz
	return s
}

// BaseSz 设置交易币投入数量（现货网格）。
func (s *TradingBotGridOrderAlgoService) BaseSz(baseSz string) *TradingBotGridOrderAlgoService {
	s.r.BaseSz = baseSz
	return s
}

// TradeQuoteCcy 设置用于交易的计价币种（现货网格）。
func (s *TradingBotGridOrderAlgoService) TradeQuoteCcy(tradeQuoteCcy string) *TradingBotGridOrderAlgoService {
	s.r.TradeQuoteCcy = tradeQuoteCcy
	return s
}

// Sz 设置投入保证金（合约网格，单位 USDT）。
func (s *TradingBotGridOrderAlgoService) Sz(sz string) *TradingBotGridOrderAlgoService {
	s.r.Sz = sz
	return s
}

// Direction 设置合约网格类型（long/short/neutral）。
func (s *TradingBotGridOrderAlgoService) Direction(direction string) *TradingBotGridOrderAlgoService {
	s.r.Direction = direction
	return s
}

// Lever 设置杠杆倍数（合约网格）。
func (s *TradingBotGridOrderAlgoService) Lever(lever string) *TradingBotGridOrderAlgoService {
	s.r.Lever = lever
	return s
}

// BasePos 设置是否开底仓（合约网格）。
func (s *TradingBotGridOrderAlgoService) BasePos(basePos bool) *TradingBotGridOrderAlgoService {
	s.r.BasePos = &basePos
	return s
}

func (s *TradingBotGridOrderAlgoService) TpRatio(tpRatio string) *TradingBotGridOrderAlgoService {
	s.r.TpRatio = tpRatio
	return s
}

func (s *TradingBotGridOrderAlgoService) SlRatio(slRatio string) *TradingBotGridOrderAlgoService {
	s.r.SlRatio = slRatio
	return s
}

var (
	errTradingBotGridOrderAlgoMissingRequired     = errors.New("okx: tradingBot grid order-algo requires instId, algoOrdType, maxPx, minPx and gridNum")
	errTradingBotGridOrderAlgoMissingSpotSize     = errors.New("okx: tradingBot grid order-algo spot requires quoteSz or baseSz")
	errTradingBotGridOrderAlgoMissingContractSize = errors.New("okx: tradingBot grid order-algo contract requires sz, direction and lever")
	errTradingBotGridOrderAlgoInvalidTriggerParam = errors.New("okx: tradingBot grid order-algo invalid triggerParams")
	errEmptyTradingBotGridOrderAlgoResponse       = errors.New("okx: empty tradingBot grid order-algo response")
)

// Do 网格策略委托下单（POST /api/v5/tradingBot/grid/order-algo）。
func (s *TradingBotGridOrderAlgoService) Do(ctx context.Context) (*TradingBotOrderAck, error) {
	if s.r.InstId == "" || s.r.AlgoOrdType == "" || s.r.MaxPx == "" || s.r.MinPx == "" || s.r.GridNum == "" {
		return nil, errTradingBotGridOrderAlgoMissingRequired
	}

	switch s.r.AlgoOrdType {
	case "grid":
		if s.r.QuoteSz == "" && s.r.BaseSz == "" {
			return nil, errTradingBotGridOrderAlgoMissingSpotSize
		}
	case "contract_grid":
		if s.r.Sz == "" || s.r.Direction == "" || s.r.Lever == "" {
			return nil, errTradingBotGridOrderAlgoMissingContractSize
		}
	}

	if len(s.r.TriggerParams) > 0 {
		for _, p := range s.r.TriggerParams {
			if p.TriggerAction == "" || p.TriggerStrategy == "" {
				return nil, errTradingBotGridOrderAlgoInvalidTriggerParam
			}
		}
	}

	var data []TradingBotOrderAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/tradingBot/grid/order-algo", nil, s.r, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyTradingBotGridOrderAlgoResponse
	}
	if data[0].SCode != "" && data[0].SCode != "0" {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodPost,
			RequestPath: "/api/v5/tradingBot/grid/order-algo",
			RequestID:   requestID,
			Code:        data[0].SCode,
			Message:     data[0].SMsg,
		}
	}
	return &data[0], nil
}
