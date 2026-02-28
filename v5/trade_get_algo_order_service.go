package okx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// GetAlgoOrderService 获取策略委托单信息（单笔）。
type GetAlgoOrderService struct {
	c *Client

	algoId      string
	algoClOrdId string
}

// NewGetAlgoOrderService 创建 GetAlgoOrderService。
func (c *Client) NewGetAlgoOrderService() *GetAlgoOrderService {
	return &GetAlgoOrderService{c: c}
}

func (s *GetAlgoOrderService) AlgoId(algoId string) *GetAlgoOrderService {
	s.algoId = algoId
	return s
}

func (s *GetAlgoOrderService) AlgoClOrdId(algoClOrdId string) *GetAlgoOrderService {
	s.algoClOrdId = algoClOrdId
	return s
}

var (
	errGetAlgoOrderMissingId       = errors.New("okx: get algo order requires algoId or algoClOrdId")
	errEmptyGetAlgoOrderResponse   = errors.New("okx: empty get algo order response")
	errInvalidGetAlgoOrderResponse = errors.New("okx: invalid get algo order response")
)

// Do 获取策略委托单信息（GET /api/v5/trade/order-algo）。
func (s *GetAlgoOrderService) Do(ctx context.Context) (*TradeAlgoOrder, error) {
	if s.algoId == "" && s.algoClOrdId == "" {
		return nil, errGetAlgoOrderMissingId
	}

	q := url.Values{}
	if s.algoId != "" {
		q.Set("algoId", s.algoId)
	} else {
		q.Set("algoClOrdId", s.algoClOrdId)
	}

	var data []TradeAlgoOrder
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodGet, "/api/v5/trade/order-algo", q, nil, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodGet, "/api/v5/trade/order-algo", requestID, errEmptyGetAlgoOrderResponse)
	}
	if len(data) != 1 {
		return nil, newInvalidDataAPIError(
			http.MethodGet,
			"/api/v5/trade/order-algo",
			requestID,
			fmt.Errorf("%w: expected 1 item, got %d", errInvalidGetAlgoOrderResponse, len(data)),
		)
	}
	return &data[0], nil
}
