package okx

import (
	"errors"
	"strings"
	"testing"
)

func TestAPIError_Error_IncludesRequestID(t *testing.T) {
	err := &APIError{
		HTTPStatus:  400,
		Method:      "GET",
		RequestPath: "/api/v5/demo",
		Code:        "51000",
		Message:     "bad",
		RequestID:   "rid-123",
	}
	if !strings.Contains(err.Error(), "requestId=rid-123") {
		t.Fatalf("Error() = %q", err.Error())
	}
}

func TestAPIError_Unwrap(t *testing.T) {
	root := errors.New("root")
	err := &APIError{
		HTTPStatus:  200,
		Method:      "POST",
		RequestPath: "/api/v5/demo",
		Code:        "0",
		Message:     "invalid ack",
		Err:         root,
	}
	if !errors.Is(err, root) {
		t.Fatalf("expected errors.Is(err, root)")
	}
}
