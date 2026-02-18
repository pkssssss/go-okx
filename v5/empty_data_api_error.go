package okx

import "net/http"

func newEmptyDataAPIError(method, requestPath, requestID string, emptyErr error) *APIError {
	message := "okx: empty response data"
	if emptyErr != nil {
		message = emptyErr.Error()
	}
	return &APIError{
		HTTPStatus:  http.StatusOK,
		Method:      method,
		RequestPath: requestPath,
		RequestID:   requestID,
		Code:        "0",
		Message:     message,
		Err:         emptyErr,
	}
}
