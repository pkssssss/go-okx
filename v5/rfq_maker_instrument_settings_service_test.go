package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRFQMakerInstrumentSettingsService_Do(t *testing.T) {
	fixedNow := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, http.MethodGet; got != want {
			t.Fatalf("method = %q, want %q", got, want)
		}
		if got, want := r.URL.Path, "/api/v5/rfq/maker-instrument-settings"; got != want {
			t.Fatalf("path = %q, want %q", got, want)
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
		_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"instType":"OPTION","includeAll":true,"data":[{"instFamily":"BTC-USD","maxBlockSz":"10000","makerPxBand":"5"}]},{"instType":"SPOT","includeAll":false,"data":[{"instId":"BTC-USDT"}]}]}`))
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

	settings, err := c.NewRFQMakerInstrumentSettingsService().Do(context.Background())
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if len(settings) != 2 || settings[0].InstType != "OPTION" || settings[1].InstType != "SPOT" {
		t.Fatalf("settings = %#v", settings)
	}
	if len(settings[0].Data) != 1 || settings[0].Data[0].InstFamily != "BTC-USD" {
		t.Fatalf("settings[0].Data = %#v", settings[0].Data)
	}
}
