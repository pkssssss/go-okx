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
