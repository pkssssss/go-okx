package okx

import "net/http"

func newInvalidDataAPIError(method, requestPath, requestID string, invalidErr error) *APIError {
	message := "okx: invalid response data"
	if invalidErr != nil {
		message = invalidErr.Error()
	}
	return &APIError{
		HTTPStatus:  http.StatusOK,
		Method:      method,
		RequestPath: requestPath,
		RequestID:   requestID,
		Code:        "0",
		Message:     message,
		Err:         invalidErr,
	}
}
