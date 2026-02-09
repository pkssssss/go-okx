package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientDo_OutNil_MissingEnvelopeFailClose(t *testing.T) {
	tests := []struct {
		name      string
		response  string
		requestID string
	}{
		{
			name:      "empty_object",
			response:  `{}`,
			requestID: "rid-invalid-envelope-empty",
		},
		{
			name:      "status_error_object",
			response:  `{"status":"error"}`,
			requestID: "rid-invalid-envelope-status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-Request-Id", tt.requestID)
				_, _ = w.Write([]byte(tt.response))
			}))
			t.Cleanup(srv.Close)

			c := NewClient(
				WithBaseURL(srv.URL),
				WithHTTPClient(srv.Client()),
			)

			err := c.do(context.Background(), http.MethodPost, "/api/v5/test", nil, map[string]string{"k": "v"}, false, nil)
			if err == nil {
				t.Fatalf("expected error")
			}

			apiErr, ok := err.(*APIError)
			if !ok {
				t.Fatalf("error type = %T, want *APIError", err)
			}
			if got, want := apiErr.HTTPStatus, http.StatusOK; got != want {
				t.Fatalf("apiErr.HTTPStatus = %d, want %d", got, want)
			}
			if got, want := apiErr.Method, http.MethodPost; got != want {
				t.Fatalf("apiErr.Method = %q, want %q", got, want)
			}
			if got, want := apiErr.RequestPath, "/api/v5/test"; got != want {
				t.Fatalf("apiErr.RequestPath = %q, want %q", got, want)
			}
			if got, want := apiErr.RequestID, tt.requestID; got != want {
				t.Fatalf("apiErr.RequestID = %q, want %q", got, want)
			}
			if got, want := apiErr.Message, "invalid response envelope"; got != want {
				t.Fatalf("apiErr.Message = %q, want %q", got, want)
			}
		})
	}
}

func TestClientDo_NonEnvelope_OutStruct_Compatible(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	t.Cleanup(srv.Close)

	c := NewClient(
		WithBaseURL(srv.URL),
		WithHTTPClient(srv.Client()),
	)

	var out struct {
		Status string `json:"status"`
	}
	err := c.do(context.Background(), http.MethodGet, "/api/v5/test", nil, nil, false, &out)
	if err != nil {
		t.Fatalf("do() error = %v", err)
	}
	if got, want := out.Status, "ok"; got != want {
		t.Fatalf("out.Status = %q, want %q", got, want)
	}
}
