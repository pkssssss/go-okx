package okx

import (
	"net/http"
	"testing"
)

func TestErrorClassification(t *testing.T) {
	t.Run("auth_error", func(t *testing.T) {
		err := &APIError{HTTPStatus: http.StatusUnauthorized, Code: "50111"}
		if !IsAuthError(err) {
			t.Fatalf("expected IsAuthError=true")
		}
		if IsRateLimitError(err) {
			t.Fatalf("expected IsRateLimitError=false")
		}
		if IsTimeSkewError(err) {
			t.Fatalf("expected IsTimeSkewError=false")
		}
	})

	t.Run("rate_limit_error", func(t *testing.T) {
		err := &APIError{HTTPStatus: http.StatusOK, Code: "50011"}
		if IsAuthError(err) {
			t.Fatalf("expected IsAuthError=false")
		}
		if !IsRateLimitError(err) {
			t.Fatalf("expected IsRateLimitError=true")
		}
		if IsTimeSkewError(err) {
			t.Fatalf("expected IsTimeSkewError=false")
		}
	})

	t.Run("time_skew_error", func(t *testing.T) {
		err := &APIError{HTTPStatus: http.StatusUnauthorized, Code: "50102"}
		if !IsTimeSkewError(err) {
			t.Fatalf("expected IsTimeSkewError=true")
		}
		if !IsAuthError(err) {
			t.Fatalf("expected IsAuthError=true")
		}
		if IsRateLimitError(err) {
			t.Fatalf("expected IsRateLimitError=false")
		}
	})

	t.Run("other_error", func(t *testing.T) {
		err := &APIError{HTTPStatus: http.StatusBadRequest, Code: "51000"}
		if IsAuthError(err) || IsRateLimitError(err) || IsTimeSkewError(err) {
			t.Fatalf("unexpected classification: auth=%v rate=%v time=%v", IsAuthError(err), IsRateLimitError(err), IsTimeSkewError(err))
		}
	})
}
