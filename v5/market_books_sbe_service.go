package okx

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkssssss/go-okx/v5/internal/rest"
)

// MarketBooksSBEService 获取 SBE 订单簿快照（二进制）。
//
// 成功响应：
// - Content-Type: application/sbe
// - Body: SBE 二进制数据（schema: SnapshotDepthResponseEvent，templateId=1006）
//
// 失败响应：
// - Content-Type: application/json
// - Body: {"code":"...","msg":"...","data":[]}
type MarketBooksSBEService struct {
	c *Client

	instIdCode *int64
	source     int
}

// NewMarketBooksSBEService 创建 MarketBooksSBEService。
// source 当前仅支持 0（普通）。
func (c *Client) NewMarketBooksSBEService() *MarketBooksSBEService {
	return &MarketBooksSBEService{c: c, source: 0}
}

// InstIdCode 设置产品 ID 唯一标识码（必填）。
func (s *MarketBooksSBEService) InstIdCode(instIdCode int64) *MarketBooksSBEService {
	s.instIdCode = &instIdCode
	return s
}

// Source 设置订单簿来源（当前仅 0: 普通）。
func (s *MarketBooksSBEService) Source(source int) *MarketBooksSBEService {
	s.source = source
	return s
}

var errMarketBooksSBEMissingInstIdCode = errors.New("okx: market books sbe requires instIdCode")

// Do 获取 SBE 订单簿快照（GET /api/v5/market/books-sbe）。
func (s *MarketBooksSBEService) Do(ctx context.Context) ([]byte, error) {
	if s.instIdCode == nil {
		return nil, errMarketBooksSBEMissingInstIdCode
	}

	q := url.Values{}
	q.Set("instIdCode", strconv.FormatInt(*s.instIdCode, 10))
	q.Set("source", strconv.Itoa(s.source))

	endpoint := "/api/v5/market/books-sbe"
	requestPath := rest.BuildRequestPath(endpoint, q)

	header := make(http.Header)
	header.Set("Accept", "application/sbe,application/json")
	if s.c.demo {
		header.Set("x-simulated-trading", "1")
	}

	status, resp, respHeader, err := s.c.rest.Do(ctx, http.MethodGet, requestPath, nil, header)
	if err != nil {
		return nil, err
	}

	contentType := respHeader.Get("Content-Type")
	if status >= http.StatusBadRequest || strings.Contains(contentType, "application/json") {
		var env responseEnvelope
		if err := json.Unmarshal(resp, &env); err != nil {
			return nil, &APIError{
				HTTPStatus:  status,
				Method:      http.MethodGet,
				RequestPath: requestPath,
				Message:     "invalid JSON response",
				Raw:         resp,
				RequestID:   respHeader.Get("x-request-id"),
			}
		}

		if env.Code != "" && env.Code != "0" {
			return nil, &APIError{
				HTTPStatus:  status,
				Method:      http.MethodGet,
				RequestPath: requestPath,
				Code:        env.Code,
				Message:     env.Msg,
				Raw:         resp,
				RequestID:   respHeader.Get("x-request-id"),
			}
		}

		return nil, &APIError{
			HTTPStatus:  status,
			Method:      http.MethodGet,
			RequestPath: requestPath,
			Message:     "unexpected JSON response",
			Raw:         resp,
			RequestID:   respHeader.Get("x-request-id"),
		}
	}

	return resp, nil
}
