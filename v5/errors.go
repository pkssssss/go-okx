package okx

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// APIError 表示 OKX API 返回的错误（HTTP 错误或业务 code != "0"）。
//
// 约定：业务成功时 code == "0"。
type APIError struct {
	HTTPStatus  int
	Method      string
	RequestPath string
	RequestID   string

	Code    string
	Message string
	Raw     []byte
	Err     error
}

func (e *APIError) Error() string {
	if e == nil {
		return "<OKX APIError>"
	}
	requestIDPart := ""
	if e.RequestID != "" {
		requestIDPart = " requestId=" + e.RequestID
	}
	if e.Code != "" || e.Message != "" {
		return fmt.Sprintf("<OKX APIError> http=%d code=%s msg=%s method=%s path=%s%s", e.HTTPStatus, e.Code, e.Message, e.Method, e.RequestPath, requestIDPart)
	}
	return fmt.Sprintf("<OKX APIError> http=%d method=%s path=%s%s", e.HTTPStatus, e.Method, e.RequestPath, requestIDPart)
}

func (e *APIError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// RequestStage 表示 REST 请求失败发生的阶段。
// 用于区分“未发送”（例如 gate 排队超时）与“已发送但未收到响应”（例如网络/服务端超时）。
type RequestStage string

const (
	RequestStagePreflight RequestStage = "preflight"
	RequestStageGate      RequestStage = "gate"
	RequestStageHTTP      RequestStage = "http"
)

// RequestStateError 表示 REST 请求在“未形成 HTTP 响应”之前失败的错误。
//
// 约定：
// - 若是 *APIError，则表示已收到 HTTP 响应（可能是业务错误/HTTP 错误）。
// - 若是 *RequestStateError：
//   - Dispatched=false 表示请求未发出（例如 gate 排队/获取并发名额阶段失败）；
//   - Dispatched=true 表示已调用底层 HTTP Do（请求已发出或已尝试发出），但未拿到可解析的响应。
type RequestStateError struct {
	Stage       RequestStage
	Dispatched  bool
	Method      string
	RequestPath string
	Err         error
}

func (e *RequestStateError) Error() string {
	if e == nil {
		return "<OKX RequestStateError>"
	}
	return fmt.Sprintf("okx: request failed stage=%s dispatched=%t method=%s path=%s: %v", e.Stage, e.Dispatched, e.Method, e.RequestPath, e.Err)
}

func (e *RequestStateError) Unwrap() error { return e.Err }

// IsAPIError 判断 err 是否为 APIError。
func IsAPIError(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr)
}

// IsAuthError 判断 err 是否为鉴权/授权相关错误。
func IsAuthError(err error) bool {
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		return false
	}
	if apiErr.HTTPStatus == http.StatusUnauthorized {
		return true
	}
	// OKX：API 类错误码（501xx）覆盖 APIKey/签名/时间戳等问题。
	return strings.HasPrefix(apiErr.Code, "501")
}

// IsRateLimitError 判断 err 是否为限速错误。
func IsRateLimitError(err error) bool {
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		return false
	}
	if apiErr.HTTPStatus == http.StatusTooManyRequests {
		return true
	}
	switch apiErr.Code {
	case "50011", "50061":
		return true
	default:
		return false
	}
}

// IsTimeSkewError 判断 err 是否为时间戳相关错误（常见于本地时间偏差或时间戳格式错误）。
func IsTimeSkewError(err error) bool {
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		return false
	}
	switch apiErr.Code {
	case "50102", "50112":
		return true
	default:
		return false
	}
}
