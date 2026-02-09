package okx

import (
	"context"
	"errors"
	"net/http"
)

type tradingBotGridAmendOrderAlgoRequest struct {
	AlgoId string `json:"algoId"`
	InstId string `json:"instId"`

	SlTriggerPx *string `json:"slTriggerPx,omitempty"`
	TpTriggerPx *string `json:"tpTriggerPx,omitempty"`
	TpRatio     *string `json:"tpRatio,omitempty"`
	SlRatio     *string `json:"slRatio,omitempty"`

	TopUpAmt      string                       `json:"topUpAmt,omitempty"`
	TriggerParams []TradingBotGridTriggerParam `json:"triggerParams,omitempty"`
}

// TradingBotGridAmendOrderAlgoService 修改网格策略订单。
type TradingBotGridAmendOrderAlgoService struct {
	c *Client
	r tradingBotGridAmendOrderAlgoRequest
}

// NewTradingBotGridAmendOrderAlgoService 创建 TradingBotGridAmendOrderAlgoService。
func (c *Client) NewTradingBotGridAmendOrderAlgoService() *TradingBotGridAmendOrderAlgoService {
	return &TradingBotGridAmendOrderAlgoService{c: c}
}

func (s *TradingBotGridAmendOrderAlgoService) AlgoId(algoId string) *TradingBotGridAmendOrderAlgoService {
	s.r.AlgoId = algoId
	return s
}

func (s *TradingBotGridAmendOrderAlgoService) InstId(instId string) *TradingBotGridAmendOrderAlgoService {
	s.r.InstId = instId
	return s
}

// SlTriggerPx 设置新的止损触发价（可选；传空字符串表示取消止损触发价）。
func (s *TradingBotGridAmendOrderAlgoService) SlTriggerPx(slTriggerPx string) *TradingBotGridAmendOrderAlgoService {
	s.r.SlTriggerPx = &slTriggerPx
	return s
}

// TpTriggerPx 设置新的止盈触发价（可选；传空字符串表示取消止盈触发价）。
func (s *TradingBotGridAmendOrderAlgoService) TpTriggerPx(tpTriggerPx string) *TradingBotGridAmendOrderAlgoService {
	s.r.TpTriggerPx = &tpTriggerPx
	return s
}

// TpRatio 设置止盈比率（可选；合约网格；传空字符串表示取消）。
func (s *TradingBotGridAmendOrderAlgoService) TpRatio(tpRatio string) *TradingBotGridAmendOrderAlgoService {
	s.r.TpRatio = &tpRatio
	return s
}

// SlRatio 设置止损比率（可选；合约网格；传空字符串表示取消）。
func (s *TradingBotGridAmendOrderAlgoService) SlRatio(slRatio string) *TradingBotGridAmendOrderAlgoService {
	s.r.SlRatio = &slRatio
	return s
}

// TopUpAmt 设置增加的投资额（可选，仅适用于现货网格）。
func (s *TradingBotGridAmendOrderAlgoService) TopUpAmt(topUpAmt string) *TradingBotGridAmendOrderAlgoService {
	s.r.TopUpAmt = topUpAmt
	return s
}

// TriggerParams 设置信号触发参数（可选）。
func (s *TradingBotGridAmendOrderAlgoService) TriggerParams(params []TradingBotGridTriggerParam) *TradingBotGridAmendOrderAlgoService {
	s.r.TriggerParams = params
	return s
}

var (
	errTradingBotGridAmendOrderAlgoMissingRequired = errors.New("okx: tradingBot grid amend-order-algo requires algoId and instId")
	errTradingBotGridAmendOrderAlgoMissingUpdate   = errors.New("okx: tradingBot grid amend-order-algo requires at least one update field")
	errTradingBotGridAmendOrderAlgoInvalidTrigger  = errors.New("okx: tradingBot grid amend-order-algo invalid triggerParams")
	errEmptyTradingBotGridAmendOrderAlgoResponse   = errors.New("okx: empty tradingBot grid amend-order-algo response")
)

// Do 修改网格策略订单（POST /api/v5/tradingBot/grid/amend-order-algo）。
func (s *TradingBotGridAmendOrderAlgoService) Do(ctx context.Context) (*TradingBotOrderAck, error) {
	if s.r.AlgoId == "" || s.r.InstId == "" {
		return nil, errTradingBotGridAmendOrderAlgoMissingRequired
	}

	hasUpdate := s.r.SlTriggerPx != nil || s.r.TpTriggerPx != nil || s.r.TpRatio != nil || s.r.SlRatio != nil || s.r.TopUpAmt != "" || len(s.r.TriggerParams) > 0
	if !hasUpdate {
		return nil, errTradingBotGridAmendOrderAlgoMissingUpdate
	}

	if len(s.r.TriggerParams) > 0 {
		for _, p := range s.r.TriggerParams {
			if p.TriggerAction == "" || p.TriggerStrategy == "" {
				return nil, errTradingBotGridAmendOrderAlgoInvalidTrigger
			}
		}
	}

	var data []TradingBotOrderAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/tradingBot/grid/amend-order-algo", nil, s.r, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyTradingBotGridAmendOrderAlgoResponse
	}
	if data[0].SCode != "" && data[0].SCode != "0" {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodPost,
			RequestPath: "/api/v5/tradingBot/grid/amend-order-algo",
			RequestID:   requestID,
			Code:        data[0].SCode,
			Message:     data[0].SMsg,
		}
	}
	return &data[0], nil
}
