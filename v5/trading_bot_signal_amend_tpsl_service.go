package okx

import (
	"context"
	"errors"
	"net/http"
)

type tradingBotSignalAmendTPSLRequest struct {
	AlgoId           string                           `json:"algoId"`
	ExitSettingParam TradingBotSignalExitSettingParam `json:"exitSettingParam"`
}

// TradingBotSignalAmendTPSLService 修改止盈止损。
type TradingBotSignalAmendTPSLService struct {
	c *Client
	r tradingBotSignalAmendTPSLRequest
}

// NewTradingBotSignalAmendTPSLService 创建 TradingBotSignalAmendTPSLService。
func (c *Client) NewTradingBotSignalAmendTPSLService() *TradingBotSignalAmendTPSLService {
	return &TradingBotSignalAmendTPSLService{c: c}
}

func (s *TradingBotSignalAmendTPSLService) AlgoId(algoId string) *TradingBotSignalAmendTPSLService {
	s.r.AlgoId = algoId
	return s
}

func (s *TradingBotSignalAmendTPSLService) ExitSettingParam(p TradingBotSignalExitSettingParam) *TradingBotSignalAmendTPSLService {
	s.r.ExitSettingParam = p
	return s
}

var (
	errTradingBotSignalAmendTPSLMissingRequired = errors.New("okx: tradingBot signal amendTPSL requires algoId and exitSettingParam.tpSlType")
	errEmptyTradingBotSignalAmendTPSLResponse   = errors.New("okx: empty tradingBot signal amendTPSL response")
)

// Do 修改止盈止损（POST /api/v5/tradingBot/signal/amendTPSL）。
func (s *TradingBotSignalAmendTPSLService) Do(ctx context.Context) (*TradingBotAlgoIdAck, error) {
	if s.r.AlgoId == "" || s.r.ExitSettingParam.TpSlType == "" {
		return nil, errTradingBotSignalAmendTPSLMissingRequired
	}

	var data []TradingBotAlgoIdAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/tradingBot/signal/amendTPSL", nil, s.r, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyTradingBotSignalAmendTPSLResponse
	}
	return &data[0], nil
}
