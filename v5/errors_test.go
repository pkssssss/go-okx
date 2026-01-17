package okx

import (
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
