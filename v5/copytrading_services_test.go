package okx

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCopyTradingServices_Do(t *testing.T) {
	fixedNow := time.Date(2020, 12, 8, 9, 8, 57, 715_000_000, time.UTC)

	type tc struct {
		name         string
		method       string
		path         string
		rawQuery     string
		signed       bool
		bodyContains []string
		call         func(c *Client) error
	}

	cases := []tc{
		{
			name:     "private_config",
			method:   http.MethodGet,
			path:     "/api/v5/copytrading/config",
			rawQuery: "",
			signed:   true,
			call: func(c *Client) error {
				_, err := c.NewCopyTradingConfigService().Do(context.Background())
				return err
			},
		},
		{
			name:     "private_copy_settings",
			method:   http.MethodGet,
			path:     "/api/v5/copytrading/copy-settings",
			rawQuery: "instType=SWAP&uniqueCode=UC",
			signed:   true,
			call: func(c *Client) error {
				_, err := c.NewCopyTradingCopySettingsService().InstType("SWAP").UniqueCode("UC").Do(context.Background())
				return err
			},
		},
		{
			name:     "private_current_lead_traders",
			method:   http.MethodGet,
			path:     "/api/v5/copytrading/current-lead-traders",
			rawQuery: "instType=SWAP",
			signed:   true,
			call: func(c *Client) error {
				_, err := c.NewCopyTradingCurrentLeadTradersService().InstType("SWAP").Do(context.Background())
				return err
			},
		},
		{
			name:     "private_current_subpositions",
			method:   http.MethodGet,
			path:     "/api/v5/copytrading/current-subpositions",
			rawQuery: "after=1&before=2&instId=BTC-USDT-SWAP&instType=SWAP&limit=10",
			signed:   true,
			call: func(c *Client) error {
				_, err := c.NewCopyTradingCurrentSubpositionsService().
					InstType("SWAP").
					InstId("BTC-USDT-SWAP").
					After("1").
					Before("2").
					Limit(10).
					Do(context.Background())
				return err
			},
		},
		{
			name:     "private_subpositions_history",
			method:   http.MethodGet,
			path:     "/api/v5/copytrading/subpositions-history",
			rawQuery: "after=1&before=2&instId=BTC-USDT-SWAP&instType=SWAP&limit=10",
			signed:   true,
			call: func(c *Client) error {
				_, err := c.NewCopyTradingSubpositionsHistoryService().
					InstType("SWAP").
					InstId("BTC-USDT-SWAP").
					After("1").
					Before("2").
					Limit(10).
					Do(context.Background())
				return err
			},
		},
		{
			name:     "private_instruments",
			method:   http.MethodGet,
			path:     "/api/v5/copytrading/instruments",
			rawQuery: "instType=SWAP",
			signed:   true,
			call: func(c *Client) error {
				_, err := c.NewCopyTradingInstrumentsService().InstType("SWAP").Do(context.Background())
				return err
			},
		},
		{
			name:     "private_profit_sharing_details",
			method:   http.MethodGet,
			path:     "/api/v5/copytrading/profit-sharing-details",
			rawQuery: "after=1&before=2&instType=SWAP&limit=10",
			signed:   true,
			call: func(c *Client) error {
				_, err := c.NewCopyTradingProfitSharingDetailsService().
					InstType("SWAP").
					After("1").
					Before("2").
					Limit(10).
					Do(context.Background())
				return err
			},
		},
		{
			name:     "private_total_profit_sharing",
			method:   http.MethodGet,
			path:     "/api/v5/copytrading/total-profit-sharing",
			rawQuery: "instType=SWAP",
			signed:   true,
			call: func(c *Client) error {
				_, err := c.NewCopyTradingTotalProfitSharingService().InstType("SWAP").Do(context.Background())
				return err
			},
		},
		{
			name:     "private_unrealized_profit_sharing_details",
			method:   http.MethodGet,
			path:     "/api/v5/copytrading/unrealized-profit-sharing-details",
			rawQuery: "instType=SWAP",
			signed:   true,
			call: func(c *Client) error {
				_, err := c.NewCopyTradingUnrealizedProfitSharingDetailsService().InstType("SWAP").Do(context.Background())
				return err
			},
		},
		{
			name:     "private_total_unrealized_profit_sharing",
			method:   http.MethodGet,
			path:     "/api/v5/copytrading/total-unrealized-profit-sharing",
			rawQuery: "instType=SWAP",
			signed:   true,
			call: func(c *Client) error {
				_, err := c.NewCopyTradingTotalUnrealizedProfitSharingService().InstType("SWAP").Do(context.Background())
				return err
			},
		},
		{
			name:     "private_amend_profit_sharing_ratio",
			method:   http.MethodPost,
			path:     "/api/v5/copytrading/amend-profit-sharing-ratio",
			rawQuery: "",
			signed:   true,
			bodyContains: []string{
				`"instType":"SWAP"`,
				`"profitSharingRatio":"0.1"`,
			},
			call: func(c *Client) error {
				_, err := c.NewCopyTradingAmendProfitSharingRatioService().InstType("SWAP").ProfitSharingRatio("0.1").Do(context.Background())
				return err
			},
		},
		{
			name:     "private_set_instruments",
			method:   http.MethodPost,
			path:     "/api/v5/copytrading/set-instruments",
			rawQuery: "",
			signed:   true,
			bodyContains: []string{
				`"instId":"BTC-USDT-SWAP,ETH-USDT-SWAP"`,
			},
			call: func(c *Client) error {
				_, err := c.NewCopyTradingSetInstrumentsService().InstType("SWAP").InstId("BTC-USDT-SWAP,ETH-USDT-SWAP").Do(context.Background())
				return err
			},
		},
		{
			name:     "private_stop_copy_trading",
			method:   http.MethodPost,
			path:     "/api/v5/copytrading/stop-copy-trading",
			rawQuery: "",
			signed:   true,
			bodyContains: []string{
				`"uniqueCode":"UC"`,
				`"subPosCloseType":"manual_close"`,
			},
			call: func(c *Client) error {
				_, err := c.NewCopyTradingStopCopyTradingService().InstType("SWAP").UniqueCode("UC").SubPosCloseType("manual_close").Do(context.Background())
				return err
			},
		},
		{
			name:     "private_close_subposition",
			method:   http.MethodPost,
			path:     "/api/v5/copytrading/close-subposition",
			rawQuery: "",
			signed:   true,
			bodyContains: []string{
				`"subPosId":"SPID"`,
			},
			call: func(c *Client) error {
				_, err := c.NewCopyTradingCloseSubpositionService().InstType("SWAP").SubPosId("SPID").Do(context.Background())
				return err
			},
		},
		{
			name:     "private_algo_order",
			method:   http.MethodPost,
			path:     "/api/v5/copytrading/algo-order",
			rawQuery: "",
			signed:   true,
			bodyContains: []string{
				`"subPosId":"SPID"`,
				`"tpTriggerPx":"10000"`,
			},
			call: func(c *Client) error {
				_, err := c.NewCopyTradingAlgoOrderService().InstType("SWAP").SubPosId("SPID").TpTriggerPx("10000").Do(context.Background())
				return err
			},
		},
		{
			name:     "private_first_copy_settings",
			method:   http.MethodPost,
			path:     "/api/v5/copytrading/first-copy-settings",
			rawQuery: "",
			signed:   true,
			bodyContains: []string{
				`"uniqueCode":"UC"`,
				`"copyMgnMode":"cross"`,
				`"copyInstIdType":"copy"`,
				`"copyTotalAmt":"500"`,
				`"copyAmt":"20"`,
				`"subPosCloseType":"copy_close"`,
			},
			call: func(c *Client) error {
				_, err := c.NewCopyTradingFirstCopySettingsService().
					InstType("SWAP").
					UniqueCode("UC").
					CopyMgnMode("cross").
					CopyInstIdType("copy").
					CopyTotalAmt("500").
					CopyAmt("20").
					SubPosCloseType("copy_close").
					Do(context.Background())
				return err
			},
		},
		{
			name:     "private_amend_copy_settings",
			method:   http.MethodPost,
			path:     "/api/v5/copytrading/amend-copy-settings",
			rawQuery: "",
			signed:   true,
			bodyContains: []string{
				`"uniqueCode":"UC"`,
				`"copyMgnMode":"cross"`,
				`"copyInstIdType":"copy"`,
				`"copyTotalAmt":"500"`,
				`"copyAmt":"20"`,
				`"subPosCloseType":"copy_close"`,
			},
			call: func(c *Client) error {
				_, err := c.NewCopyTradingAmendCopySettingsService().
					InstType("SWAP").
					UniqueCode("UC").
					CopyMgnMode("cross").
					CopyInstIdType("copy").
					CopyTotalAmt("500").
					CopyAmt("20").
					SubPosCloseType("copy_close").
					Do(context.Background())
				return err
			},
		},

		{
			name:     "public_config",
			method:   http.MethodGet,
			path:     "/api/v5/copytrading/public-config",
			rawQuery: "instType=SWAP",
			signed:   false,
			call: func(c *Client) error {
				_, err := c.NewCopyTradingPublicConfigService().InstType("SWAP").Do(context.Background())
				return err
			},
		},
		{
			name:     "public_lead_traders",
			method:   http.MethodGet,
			path:     "/api/v5/copytrading/public-lead-traders",
			rawQuery: "dataVer=1&instType=SWAP&limit=2&page=3&sortType=overview&state=1",
			signed:   false,
			call: func(c *Client) error {
				_, err := c.NewCopyTradingPublicLeadTradersService().
					InstType("SWAP").
					SortType("overview").
					State("1").
					DataVer("1").
					Page(3).
					Limit(2).
					Do(context.Background())
				return err
			},
		},
		{
			name:     "public_weekly_pnl",
			method:   http.MethodGet,
			path:     "/api/v5/copytrading/public-weekly-pnl",
			rawQuery: "instType=SWAP&uniqueCode=UC",
			signed:   false,
			call: func(c *Client) error {
				_, err := c.NewCopyTradingPublicWeeklyPnlService().InstType("SWAP").UniqueCode("UC").Do(context.Background())
				return err
			},
		},
		{
			name:     "public_pnl",
			method:   http.MethodGet,
			path:     "/api/v5/copytrading/public-pnl",
			rawQuery: "instType=SWAP&lastDays=1&uniqueCode=UC",
			signed:   false,
			call: func(c *Client) error {
				_, err := c.NewCopyTradingPublicPnlService().InstType("SWAP").UniqueCode("UC").LastDays("1").Do(context.Background())
				return err
			},
		},
		{
			name:     "public_stats",
			method:   http.MethodGet,
			path:     "/api/v5/copytrading/public-stats",
			rawQuery: "instType=SWAP&lastDays=1&uniqueCode=UC",
			signed:   false,
			call: func(c *Client) error {
				_, err := c.NewCopyTradingPublicStatsService().InstType("SWAP").UniqueCode("UC").LastDays("1").Do(context.Background())
				return err
			},
		},
		{
			name:     "public_preference_currency",
			method:   http.MethodGet,
			path:     "/api/v5/copytrading/public-preference-currency",
			rawQuery: "instType=SWAP&uniqueCode=UC",
			signed:   false,
			call: func(c *Client) error {
				_, err := c.NewCopyTradingPublicPreferenceCurrencyService().InstType("SWAP").UniqueCode("UC").Do(context.Background())
				return err
			},
		},
		{
			name:     "public_current_subpositions",
			method:   http.MethodGet,
			path:     "/api/v5/copytrading/public-current-subpositions",
			rawQuery: "after=1&before=2&instType=SWAP&limit=10&uniqueCode=UC",
			signed:   false,
			call: func(c *Client) error {
				_, err := c.NewCopyTradingPublicCurrentSubpositionsService().
					InstType("SWAP").
					UniqueCode("UC").
					After("1").
					Before("2").
					Limit(10).
					Do(context.Background())
				return err
			},
		},
		{
			name:     "public_subpositions_history",
			method:   http.MethodGet,
			path:     "/api/v5/copytrading/public-subpositions-history",
			rawQuery: "after=1&before=2&instType=SWAP&limit=10&uniqueCode=UC",
			signed:   false,
			call: func(c *Client) error {
				_, err := c.NewCopyTradingPublicSubpositionsHistoryService().
					InstType("SWAP").
					UniqueCode("UC").
					After("1").
					Before("2").
					Limit(10).
					Do(context.Background())
				return err
			},
		},
		{
			name:     "public_copy_traders",
			method:   http.MethodGet,
			path:     "/api/v5/copytrading/public-copy-traders",
			rawQuery: "instType=SWAP&limit=10&uniqueCode=UC",
			signed:   false,
			call: func(c *Client) error {
				_, err := c.NewCopyTradingPublicCopyTradersService().
					InstType("SWAP").
					UniqueCode("UC").
					Limit(10).
					Do(context.Background())
				return err
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if got, want := r.Method, tt.method; got != want {
					t.Fatalf("method = %q, want %q", got, want)
				}
				if got, want := r.URL.Path, tt.path; got != want {
					t.Fatalf("path = %q, want %q", got, want)
				}
				if got, want := r.URL.RawQuery, tt.rawQuery; got != want {
					t.Fatalf("query = %q, want %q", got, want)
				}

				if tt.signed {
					if got, want := r.Header.Get("OK-ACCESS-KEY"), "mykey"; got != want {
						t.Fatalf("OK-ACCESS-KEY = %q, want %q", got, want)
					}
					if got := r.Header.Get("OK-ACCESS-SIGN"); got == "" {
						t.Fatalf("expected OK-ACCESS-SIGN")
					}
				} else {
					if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
						t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
					}
				}

				if len(tt.bodyContains) > 0 {
					bodyBytes, _ := io.ReadAll(r.Body)
					body := string(bodyBytes)
					for _, s := range tt.bodyContains {
						if !strings.Contains(body, s) {
							t.Fatalf("body %q missing %q", body, s)
						}
					}
				}

				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{}]}`))
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

			if err := tt.call(c); err != nil {
				t.Fatalf("Do() error = %v", err)
			}
		})
	}
}

func TestCopyTradingServices_Validation(t *testing.T) {
	t.Run("copy_settings_missing_uniqueCode", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewCopyTradingCopySettingsService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errCopyTradingCopySettingsMissingUniqueCode {
			t.Fatalf("error = %v, want %v", err, errCopyTradingCopySettingsMissingUniqueCode)
		}
	})

	t.Run("algo_order_missing_trigger_px", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewCopyTradingAlgoOrderService().SubPosId("SPID").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errCopyTradingAlgoOrderMissingTriggerPx {
			t.Fatalf("error = %v, want %v", err, errCopyTradingAlgoOrderMissingTriggerPx)
		}
	})

	t.Run("upsert_copy_settings_missing_required", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewCopyTradingFirstCopySettingsService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errCopyTradingUpsertCopySettingsMissingRequired {
			t.Fatalf("error = %v, want %v", err, errCopyTradingUpsertCopySettingsMissingRequired)
		}
	})

	t.Run("stop_copy_trading_missing_uniqueCode", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewCopyTradingStopCopyTradingService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errCopyTradingStopCopyTradingMissingUniqueCode {
			t.Fatalf("error = %v, want %v", err, errCopyTradingStopCopyTradingMissingUniqueCode)
		}
	})

	t.Run("public_weekly_pnl_missing_uniqueCode", func(t *testing.T) {
		c := NewClient()
		_, err := c.NewCopyTradingPublicWeeklyPnlService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errCopyTradingPublicWeeklyPnlMissingUniqueCode {
			t.Fatalf("error = %v, want %v", err, errCopyTradingPublicWeeklyPnlMissingUniqueCode)
		}
	})
}
