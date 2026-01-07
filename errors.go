package okx

import (
	"errors"
	"fmt"
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
	if e.Code != "" || e.Message != "" {
		return fmt.Sprintf("<OKX APIError> http=%d code=%s msg=%s method=%s path=%s", e.HTTPStatus, e.Code, e.Message, e.Method, e.RequestPath)
	}
	return fmt.Sprintf("<OKX APIError> http=%d method=%s path=%s", e.HTTPStatus, e.Method, e.RequestPath)
}

// IsAPIError 判断 err 是否为 APIError。
func IsAPIError(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr)
}
