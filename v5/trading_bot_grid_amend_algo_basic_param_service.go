package okx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type tradingBotGridAmendAlgoBasicParamRequest struct {
	AlgoId  string `json:"algoId"`
	MinPx   string `json:"minPx"`
	MaxPx   string `json:"maxPx"`
	GridNum string `json:"gridNum"`
}

// TradingBotGridAmendAlgoBasicParamService 修改网格策略基本参数。
type TradingBotGridAmendAlgoBasicParamService struct {
	c *Client
	r tradingBotGridAmendAlgoBasicParamRequest
}

// NewTradingBotGridAmendAlgoBasicParamService 创建 TradingBotGridAmendAlgoBasicParamService。
func (c *Client) NewTradingBotGridAmendAlgoBasicParamService() *TradingBotGridAmendAlgoBasicParamService {
	return &TradingBotGridAmendAlgoBasicParamService{c: c}
}

func (s *TradingBotGridAmendAlgoBasicParamService) AlgoId(algoId string) *TradingBotGridAmendAlgoBasicParamService {
	s.r.AlgoId = algoId
	return s
}

func (s *TradingBotGridAmendAlgoBasicParamService) MinPx(minPx string) *TradingBotGridAmendAlgoBasicParamService {
	s.r.MinPx = minPx
	return s
}

func (s *TradingBotGridAmendAlgoBasicParamService) MaxPx(maxPx string) *TradingBotGridAmendAlgoBasicParamService {
	s.r.MaxPx = maxPx
	return s
}

func (s *TradingBotGridAmendAlgoBasicParamService) GridNum(gridNum string) *TradingBotGridAmendAlgoBasicParamService {
	s.r.GridNum = gridNum
	return s
}

var (
	errTradingBotGridAmendAlgoBasicParamMissingRequired = errors.New("okx: tradingBot grid amend-algo-basic-param requires algoId, minPx, maxPx and gridNum")
	errEmptyTradingBotGridAmendAlgoBasicParamResponse   = errors.New("okx: empty tradingBot grid amend-algo-basic-param response")
	errInvalidTradingBotGridAmendAlgoBasicParamResponse = errors.New("okx: invalid tradingBot grid amend-algo-basic-param response")
)

func validateTradingBotGridAmendAlgoBasicParamResult(result *TradingBotGridAmendAlgoBasicParamResult) error {
	if result == nil {
		return errInvalidTradingBotGridAmendAlgoBasicParamResponse
	}
	if result.AlgoId == "" {
		return errInvalidTradingBotGridAmendAlgoBasicParamResponse
	}
	return nil
}

// Do 修改网格策略基本参数（POST /api/v5/tradingBot/grid/amend-algo-basic-param）。
//
// 注意：OKX 文档示例中该接口的 data 可能为对象或数组；为提高兼容性，这里同时兼容两种形态。
func (s *TradingBotGridAmendAlgoBasicParamService) Do(ctx context.Context) (*TradingBotGridAmendAlgoBasicParamResult, error) {
	if s.r.AlgoId == "" || s.r.MinPx == "" || s.r.MaxPx == "" || s.r.GridNum == "" {
		return nil, errTradingBotGridAmendAlgoBasicParamMissingRequired
	}

	var raw json.RawMessage
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/tradingBot/grid/amend-algo-basic-param", nil, s.r, true, nil, &raw)
	if err != nil {
		return nil, err
	}
	b := bytes.TrimSpace(raw)
	if len(b) == 0 || string(b) == "null" {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/tradingBot/grid/amend-algo-basic-param", requestID, errEmptyTradingBotGridAmendAlgoBasicParamResponse)
	}

	switch b[0] {
	case '{':
		var v TradingBotGridAmendAlgoBasicParamResult
		if err := json.Unmarshal(b, &v); err != nil {
			return nil, err
		}
		if err := validateTradingBotGridAmendAlgoBasicParamResult(&v); err != nil {
			return nil, newInvalidDataAPIError(http.MethodPost, "/api/v5/tradingBot/grid/amend-algo-basic-param", requestID, err)
		}
		return &v, nil
	case '[':
		var vs []TradingBotGridAmendAlgoBasicParamResult
		if err := json.Unmarshal(b, &vs); err != nil {
			return nil, err
		}
		if len(vs) == 0 {
			return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/tradingBot/grid/amend-algo-basic-param", requestID, errEmptyTradingBotGridAmendAlgoBasicParamResponse)
		}
		if err := validateTradingBotGridAmendAlgoBasicParamResult(&vs[0]); err != nil {
			return nil, newInvalidDataAPIError(http.MethodPost, "/api/v5/tradingBot/grid/amend-algo-basic-param", requestID, err)
		}
		return &vs[0], nil
	default:
		return nil, errInvalidTradingBotGridAmendAlgoBasicParamResponse
	}
}
