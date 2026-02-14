package okx

import (
	"errors"
	"testing"
)

func assertEmptyDataAPIError(t *testing.T, err error, wantErr error) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error = %T, want *APIError", err)
	}
	if got, want := apiErr.Code, "0"; got != want {
		t.Fatalf("apiErr.Code = %q, want %q", got, want)
	}
	if got, want := apiErr.Message, wantErr.Error(); got != want {
		t.Fatalf("apiErr.Message = %q, want %q", got, want)
	}
}
