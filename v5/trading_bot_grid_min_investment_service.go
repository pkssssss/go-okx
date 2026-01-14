package okx

import (
	"context"
	"errors"
	"net/http"
)

type tradingBotGridMinInvestmentRequest struct {
	InstId      string `json:"instId"`
	AlgoOrdType string `json:"algoOrdType"`
	GridNum     string `json:"gridNum"`
	MaxPx       string `json:"maxPx"`
	MinPx       string `json:"minPx"`
	RunType     string `json:"runType"`

	Direction string `json:"direction,omitempty"`
	Lever     string `json:"lever,omitempty"`
	BasePos   *bool  `json:"basePos,omitempty"`

	InvestmentType  string `json:"investmentType,omitempty"`
	TriggerStrategy string `json:"triggerStrategy,omitempty"`
	TopUpAmt        string `json:"topUpAmt,omitempty"`

	InvestmentData []TradingBotGridInvestmentData `json:"investmentData,omitempty"`
}

// TradingBotGridMinInvestmentService 计算最小投资数量（公共）。
type TradingBotGridMinInvestmentService struct {
	c *Client
	r tradingBotGridMinInvestmentRequest
}

// NewTradingBotGridMinInvestmentService 创建 TradingBotGridMinInvestmentService。
func (c *Client) NewTradingBotGridMinInvestmentService() *TradingBotGridMinInvestmentService {
	return &TradingBotGridMinInvestmentService{c: c}
}

func (s *TradingBotGridMinInvestmentService) InstId(instId string) *TradingBotGridMinInvestmentService {
	s.r.InstId = instId
	return s
}

func (s *TradingBotGridMinInvestmentService) AlgoOrdType(algoOrdType string) *TradingBotGridMinInvestmentService {
	s.r.AlgoOrdType = algoOrdType
	return s
}

func (s *TradingBotGridMinInvestmentService) GridNum(gridNum string) *TradingBotGridMinInvestmentService {
	s.r.GridNum = gridNum
	return s
}

func (s *TradingBotGridMinInvestmentService) MaxPx(maxPx string) *TradingBotGridMinInvestmentService {
	s.r.MaxPx = maxPx
	return s
}

func (s *TradingBotGridMinInvestmentService) MinPx(minPx string) *TradingBotGridMinInvestmentService {
	s.r.MinPx = minPx
	return s
}

func (s *TradingBotGridMinInvestmentService) RunType(runType string) *TradingBotGridMinInvestmentService {
	s.r.RunType = runType
	return s
}

func (s *TradingBotGridMinInvestmentService) Direction(direction string) *TradingBotGridMinInvestmentService {
	s.r.Direction = direction
	return s
}

func (s *TradingBotGridMinInvestmentService) Lever(lever string) *TradingBotGridMinInvestmentService {
	s.r.Lever = lever
	return s
}

func (s *TradingBotGridMinInvestmentService) BasePos(basePos bool) *TradingBotGridMinInvestmentService {
	s.r.BasePos = &basePos
	return s
}

func (s *TradingBotGridMinInvestmentService) InvestmentType(investmentType string) *TradingBotGridMinInvestmentService {
	s.r.InvestmentType = investmentType
	return s
}

func (s *TradingBotGridMinInvestmentService) TriggerStrategy(triggerStrategy string) *TradingBotGridMinInvestmentService {
	s.r.TriggerStrategy = triggerStrategy
	return s
}

func (s *TradingBotGridMinInvestmentService) TopUpAmt(topUpAmt string) *TradingBotGridMinInvestmentService {
	s.r.TopUpAmt = topUpAmt
	return s
}

func (s *TradingBotGridMinInvestmentService) InvestmentData(investmentData []TradingBotGridInvestmentData) *TradingBotGridMinInvestmentService {
	s.r.InvestmentData = investmentData
	return s
}

var (
	errTradingBotGridMinInvestmentMissingRequired   = errors.New("okx: tradingBot grid min-investment requires instId, algoOrdType, gridNum, maxPx, minPx and runType")
	errTradingBotGridMinInvestmentInvalidInvestment = errors.New("okx: tradingBot grid min-investment invalid investmentData")
	errEmptyTradingBotGridMinInvestmentResponse     = errors.New("okx: empty tradingBot grid min-investment response")
)

// Do 计算最小投资数量（POST /api/v5/tradingBot/grid/min-investment）。
func (s *TradingBotGridMinInvestmentService) Do(ctx context.Context) (*TradingBotGridMinInvestmentResult, error) {
	if s.r.InstId == "" || s.r.AlgoOrdType == "" || s.r.GridNum == "" || s.r.MaxPx == "" || s.r.MinPx == "" || s.r.RunType == "" {
		return nil, errTradingBotGridMinInvestmentMissingRequired
	}
	if len(s.r.InvestmentData) > 0 {
		for _, it := range s.r.InvestmentData {
			if it.Amt == "" || it.Ccy == "" {
				return nil, errTradingBotGridMinInvestmentInvalidInvestment
			}
		}
	}

	var data []TradingBotGridMinInvestmentResult
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/tradingBot/grid/min-investment", nil, s.r, false, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyTradingBotGridMinInvestmentResponse
	}
	return &data[0], nil
}
