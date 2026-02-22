package okx

import (
	"context"
	"errors"
	"fmt"
)

// ClientStats 是 Client 的 REST 运行统计快照（并发安全）。
type ClientStats struct {
	RequestTotal uint64
	SuccessTotal uint64
	FailureTotal uint64
	RetryTotal   uint64

	// ErrorCodeCounts 聚合失败请求的错误码分布：
	// - 业务错误：OKX code（如 50011）
	// - HTTP 错误：HTTP_XXX（如 HTTP_500）
	// - 其他错误：REQUEST_HTTP / CONTEXT_CANCELED / UNKNOWN 等
	ErrorCodeCounts map[string]uint64
}

// ClientStats 返回 Client 的 REST 运行统计快照（并发安全）。
func (c *Client) ClientStats() ClientStats {
	var s ClientStats
	if c == nil {
		return s
	}

	s.RequestTotal = c.statsRequestTotal.Load()
	s.SuccessTotal = c.statsSuccessTotal.Load()
	s.FailureTotal = c.statsFailureTotal.Load()
	s.RetryTotal = c.statsRetryTotal.Load()

	c.statsErrorCodeMu.Lock()
	if len(c.statsErrorCodes) > 0 {
		s.ErrorCodeCounts = make(map[string]uint64, len(c.statsErrorCodes))
		for code, n := range c.statsErrorCodes {
			s.ErrorCodeCounts[code] = n
		}
	}
	c.statsErrorCodeMu.Unlock()

	return s
}

func (c *Client) recordClientRequest() {
	if c == nil {
		return
	}
	c.statsRequestTotal.Add(1)
}

func (c *Client) recordClientSuccess() {
	if c == nil {
		return
	}
	c.statsSuccessTotal.Add(1)
}

func (c *Client) recordClientRetry() {
	if c == nil {
		return
	}
	c.statsRetryTotal.Add(1)
}

func (c *Client) recordClientFailure(err error) {
	if c == nil {
		return
	}
	c.statsFailureTotal.Add(1)

	code := classifyClientErrorCode(err)
	if code == "" {
		return
	}

	c.statsErrorCodeMu.Lock()
	if c.statsErrorCodes == nil {
		c.statsErrorCodes = make(map[string]uint64)
	}
	c.statsErrorCodes[code]++
	c.statsErrorCodeMu.Unlock()
}

func classifyClientErrorCode(err error) string {
	if err == nil {
		return ""
	}
	if errors.Is(err, errMissingCredentials) {
		return "MISSING_CREDENTIALS"
	}

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		if apiErr.Code != "" && apiErr.Code != "0" {
			return apiErr.Code
		}
		if apiErr.HTTPStatus > 0 {
			return fmt.Sprintf("HTTP_%d", apiErr.HTTPStatus)
		}
		if apiErr.Code != "" {
			return apiErr.Code
		}
		return "API_ERROR"
	}

	var reqErr *RequestStateError
	if errors.As(err, &reqErr) {
		switch reqErr.Stage {
		case RequestStagePreflight:
			return "REQUEST_PREFLIGHT"
		case RequestStageGate:
			return "REQUEST_GATE"
		case RequestStageHTTP:
			return "REQUEST_HTTP"
		default:
			return "REQUEST_UNKNOWN"
		}
	}

	if errors.Is(err, context.Canceled) {
		return "CONTEXT_CANCELED"
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return "CONTEXT_DEADLINE_EXCEEDED"
	}
	return "UNKNOWN"
}
