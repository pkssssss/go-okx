package okx

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRFQCancelBatchRFQsService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	t.Run("signed_request_and_partial_failure", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.Method, http.MethodPost; got != want {
				t.Fatalf("method = %q, want %q", got, want)
			}
			if got, want := r.URL.Path, "/api/v5/rfq/cancel-batch-rfqs"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}

			bodyBytes, _ := io.ReadAll(r.Body)
			if got, want := string(bodyBytes), `{"rfqIds":["2201","2202"],"clRfqIds":["r1","r2"]}`; got != want {
				t.Fatalf("body = %q, want %q", got, want)
			}

			if got, want := r.Header.Get("OK-ACCESS-KEY"), "mykey"; got != want {
				t.Fatalf("OK-ACCESS-KEY = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-PASSPHRASE"), "mypass"; got != want {
				t.Fatalf("OK-ACCESS-PASSPHRASE = %q, want %q", got, want)
			}
			if got, want := r.Header.Get("OK-ACCESS-TIMESTAMP"), "2020-03-28T12:21:41.274Z"; got != want {
				t.Fatalf("OK-ACCESS-TIMESTAMP = %q, want %q", got, want)
			}
			if got := r.Header.Get("OK-ACCESS-SIGN"); got == "" {
				t.Fatalf("OK-ACCESS-SIGN empty")
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"rfqId":"2201","clRfqId":"r1","sCode":"0","sMsg":""},{"rfqId":"2202","clRfqId":"r2","sCode":"70001","sMsg":"RFQ does not exist."}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
			WithCredentials(Credentials{
				APIKey:     "mykey",
				SecretKey:  "mysecret",
				Passphrase: "mypass",
			}),
			WithNowFunc(func() time.Time { return fixedNow }),
		)

		acks, err := c.NewRFQCancelBatchRFQsService().
			RfqIds([]string{"2201", "2202"}).
			ClRfqIds([]string{"r1", "r2"}).
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected partial failure error, got nil")
		}
		var batchErr *RFQCancelBatchRFQsError
		if !errors.As(err, &batchErr) {
			t.Fatalf("expected RFQCancelBatchRFQsError, got %T: %v", err, err)
		}
		if len(acks) != 2 || acks[1].SCode != "70001" {
			t.Fatalf("acks = %#v", acks)
		}
	})

	t.Run("missing_ids", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewRFQCancelBatchRFQsService().Do(context.Background())
		if !errors.Is(err, errRFQCancelBatchRFQsMissingIds) {
			t.Fatalf("expected errRFQCancelBatchRFQsMissingIds, got %T: %v", err, err)
		}
	})
}
