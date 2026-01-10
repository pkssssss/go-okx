package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublicDeliveryExerciseHistoryService_Do(t *testing.T) {
	t.Run("missing_instType", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewPublicDeliveryExerciseHistoryService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errPublicDeliveryExerciseHistoryMissingInstType {
			t.Fatalf("error = %v, want %v", err, errPublicDeliveryExerciseHistoryMissingInstType)
		}
	})

	t.Run("missing_instFamily", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewPublicDeliveryExerciseHistoryService().InstType("FUTURES").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errPublicDeliveryExerciseHistoryMissingInstFamily {
			t.Fatalf("error = %v, want %v", err, errPublicDeliveryExerciseHistoryMissingInstFamily)
		}
	})

	t.Run("ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodGet; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/public/delivery-exercise-history"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			wantQuery := "after=1&before=2&instFamily=BTC-USD&instType=FUTURES&limit=3"
			if got, want := r.URL.RawQuery, wantQuery; got != want {
				t.Fatalf("query = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
				t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"ts":"1597026383085","details":[{"type":"delivery","insId":"BTC-USD-190927","px":"0.016"}]}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
		)

		got, err := c.NewPublicDeliveryExerciseHistoryService().
			InstType("FUTURES").
			InstFamily("BTC-USD").
			After("1").
			Before("2").
			Limit(3).
			Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		if len(got) != 1 || got[0].TS != 1597026383085 {
			t.Fatalf("data = %#v", got)
		}
		if len(got[0].Details) != 1 || got[0].Details[0].InsId == "" || got[0].Details[0].Px == "" || got[0].Details[0].Type == "" {
			t.Fatalf("details = %#v", got[0])
		}
	})
}
