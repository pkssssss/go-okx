package okx

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTradingBotServices_RequestShape(t *testing.T) {
	fixedNow := time.Date(2020, 6, 30, 12, 34, 56, 789_000_000, time.UTC)

	mkClient := func(srv *httptest.Server) *Client {
		return NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
			WithCredentials(Credentials{
				APIKey:     "mykey",
				SecretKey:  "mysecret",
				Passphrase: "mypass",
			}),
			WithNowFunc(func() time.Time { return fixedNow }),
		)
	}

	type tc struct {
		name     string
		method   string
		path     string
		query    string
		body     string
		signed   bool
		response string
		invokeDo func(c *Client) error
	}

	okResp := `{"code":"0","msg":"","data":[{}]}`

	cases := []tc{
		// TradingBot Grid - public
		{
			name:   "grid_ai_param_public",
			method: http.MethodGet,
			path:   "/api/v5/tradingBot/grid/ai-param",
			query:  "algoOrdType=grid&instId=BTC-USDT",
			signed: false,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotGridAIParamService().AlgoOrdType("grid").InstId("BTC-USDT").Do(context.Background())
				return err
			},
		},
		{
			name:   "grid_grid_quantity_public",
			method: http.MethodGet,
			path:   "/api/v5/tradingBot/grid/grid-quantity",
			query:  "algoOrdType=grid&instId=BTC-USDT&maxPx=1&minPx=0&runType=1",
			signed: false,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotGridGridQuantityService().
					InstId("BTC-USDT").
					RunType("1").
					AlgoOrdType("grid").
					MaxPx("1").
					MinPx("0").
					Do(context.Background())
				return err
			},
		},
		{
			name:   "grid_min_investment_public",
			method: http.MethodPost,
			path:   "/api/v5/tradingBot/grid/min-investment",
			body:   `{"instId":"BTC-USDT","algoOrdType":"grid","gridNum":"10","maxPx":"1","minPx":"0","runType":"1"}`,
			signed: false,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotGridMinInvestmentService().
					InstId("BTC-USDT").
					AlgoOrdType("grid").
					GridNum("10").
					MaxPx("1").
					MinPx("0").
					RunType("1").
					Do(context.Background())
				return err
			},
		},
		{
			name:   "public_rsi_back_testing",
			method: http.MethodGet,
			path:   "/api/v5/tradingBot/public/rsi-back-testing",
			query:  "instId=BTC-USDT&thold=70&timePeriod=14&timeframe=1H",
			signed: false,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotPublicRSIBackTestingService().
					InstId("BTC-USDT").
					Timeframe("1H").
					Thold("70").
					TimePeriod("14").
					Do(context.Background())
				return err
			},
		},

		// TradingBot Grid - private
		{
			name:   "grid_order_algo",
			method: http.MethodPost,
			path:   "/api/v5/tradingBot/grid/order-algo",
			body:   `{"instId":"BTC-USDT","algoOrdType":"grid","maxPx":"1","minPx":"0","gridNum":"10","quoteSz":"100"}`,
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotGridOrderAlgoService().
					InstId("BTC-USDT").
					AlgoOrdType("grid").
					MaxPx("1").
					MinPx("0").
					GridNum("10").
					QuoteSz("100").
					Do(context.Background())
				return err
			},
		},
		{
			name:   "grid_amend_algo_basic_param",
			method: http.MethodPost,
			path:   "/api/v5/tradingBot/grid/amend-algo-basic-param",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotGridAmendAlgoBasicParamService().
					AlgoId("1").
					MinPx("0").
					MaxPx("1").
					GridNum("10").
					Do(context.Background())
				return err
			},
		},
		{
			name:   "grid_amend_order_algo",
			method: http.MethodPost,
			path:   "/api/v5/tradingBot/grid/amend-order-algo",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotGridAmendOrderAlgoService().
					AlgoId("1").
					InstId("BTC-USDT").
					SlTriggerPx("0.9").
					Do(context.Background())
				return err
			},
		},
		{
			name:   "grid_stop_order_algo",
			method: http.MethodPost,
			path:   "/api/v5/tradingBot/grid/stop-order-algo",
			body:   `[{"algoId":"1","instId":"BTC-USDT","algoOrdType":"grid","stopType":"1"}]`,
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotGridStopOrderAlgoService().
					Orders([]TradingBotGridStopOrder{
						{AlgoId: "1", InstId: "BTC-USDT", AlgoOrdType: "grid", StopType: "1"},
					}).
					Do(context.Background())
				return err
			},
		},
		{
			name:   "grid_orders_algo_details",
			method: http.MethodGet,
			path:   "/api/v5/tradingBot/grid/orders-algo-details",
			query:  "algoId=1&algoOrdType=grid",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotGridOrdersAlgoDetailsService().
					AlgoOrdType("grid").
					AlgoId("1").
					Do(context.Background())
				return err
			},
		},
		{
			name:   "grid_orders_algo_pending",
			method: http.MethodGet,
			path:   "/api/v5/tradingBot/grid/orders-algo-pending",
			query:  "algoOrdType=grid",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotGridOrdersAlgoPendingService().AlgoOrdType("grid").Do(context.Background())
				return err
			},
		},
		{
			name:   "grid_orders_algo_history",
			method: http.MethodGet,
			path:   "/api/v5/tradingBot/grid/orders-algo-history",
			query:  "algoOrdType=grid",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotGridOrdersAlgoHistoryService().AlgoOrdType("grid").Do(context.Background())
				return err
			},
		},
		{
			name:   "grid_positions",
			method: http.MethodGet,
			path:   "/api/v5/tradingBot/grid/positions",
			query:  "algoId=1&algoOrdType=grid",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotGridPositionsService().AlgoOrdType("grid").AlgoId("1").Do(context.Background())
				return err
			},
		},
		{
			name:   "grid_sub_orders",
			method: http.MethodGet,
			path:   "/api/v5/tradingBot/grid/sub-orders",
			query:  "algoId=1&algoOrdType=grid&type=live",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotGridSubOrdersService().AlgoId("1").AlgoOrdType("grid").Type("live").Do(context.Background())
				return err
			},
		},
		{
			name:   "grid_margin_balance",
			method: http.MethodPost,
			path:   "/api/v5/tradingBot/grid/margin-balance",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotGridMarginBalanceService().AlgoId("1").Type("add").Amt("1").Do(context.Background())
				return err
			},
		},
		{
			name:   "grid_compute_margin_balance",
			method: http.MethodPost,
			path:   "/api/v5/tradingBot/grid/compute-margin-balance",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotGridComputeMarginBalanceService().AlgoId("1").Type("add").Amt("1").Do(context.Background())
				return err
			},
		},
		{
			name:   "grid_order_instant_trigger",
			method: http.MethodPost,
			path:   "/api/v5/tradingBot/grid/order-instant-trigger",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotGridOrderInstantTriggerService().AlgoId("1").Do(context.Background())
				return err
			},
		},
		{
			name:   "grid_close_position",
			method: http.MethodPost,
			path:   "/api/v5/tradingBot/grid/close-position",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotGridClosePositionService().AlgoId("1").MktClose(true).Do(context.Background())
				return err
			},
		},
		{
			name:   "grid_cancel_close_order",
			method: http.MethodPost,
			path:   "/api/v5/tradingBot/grid/cancel-close-order",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotGridCancelCloseOrderService().AlgoId("1").OrdId("2").Do(context.Background())
				return err
			},
		},
		{
			name:   "grid_withdraw_income",
			method: http.MethodPost,
			path:   "/api/v5/tradingBot/grid/withdraw-income",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotGridWithdrawIncomeService().AlgoId("1").Do(context.Background())
				return err
			},
		},
		{
			name:   "grid_adjust_investment",
			method: http.MethodPost,
			path:   "/api/v5/tradingBot/grid/adjust-investment",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotGridAdjustInvestmentService().AlgoId("1").Amt("1").Do(context.Background())
				return err
			},
		},

		// TradingBot Recurring - private
		{
			name:   "recurring_order_algo",
			method: http.MethodPost,
			path:   "/api/v5/tradingBot/recurring/order-algo",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotRecurringOrderAlgoService().
					StgyName("stgy").
					RecurringList([]TradingBotRecurringListItem{{Ccy: "BTC", Ratio: "1"}}).
					Period("daily").
					RecurringTime("00:00").
					TimeZone("UTC").
					Amt("1").
					InvestmentCcy("USDT").
					TdMode("cash").
					Do(context.Background())
				return err
			},
		},
		{
			name:   "recurring_amend_order_algo",
			method: http.MethodPost,
			path:   "/api/v5/tradingBot/recurring/amend-order-algo",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotRecurringAmendOrderAlgoService().AlgoId("1").StgyName("stgy").Do(context.Background())
				return err
			},
		},
		{
			name:   "recurring_stop_order_algo",
			method: http.MethodPost,
			path:   "/api/v5/tradingBot/recurring/stop-order-algo",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotRecurringStopOrderAlgoService().
					Orders([]TradingBotRecurringStopOrder{{AlgoId: "1"}}).
					Do(context.Background())
				return err
			},
		},
		{
			name:   "recurring_orders_algo_pending",
			method: http.MethodGet,
			path:   "/api/v5/tradingBot/recurring/orders-algo-pending",
			query:  "algoId=1",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotRecurringOrdersAlgoPendingService().AlgoId("1").Do(context.Background())
				return err
			},
		},
		{
			name:   "recurring_orders_algo_history",
			method: http.MethodGet,
			path:   "/api/v5/tradingBot/recurring/orders-algo-history",
			query:  "algoId=1",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotRecurringOrdersAlgoHistoryService().AlgoId("1").Do(context.Background())
				return err
			},
		},
		{
			name:   "recurring_orders_algo_details",
			method: http.MethodGet,
			path:   "/api/v5/tradingBot/recurring/orders-algo-details",
			query:  "algoId=1",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotRecurringOrdersAlgoDetailsService().AlgoId("1").Do(context.Background())
				return err
			},
		},
		{
			name:   "recurring_sub_orders",
			method: http.MethodGet,
			path:   "/api/v5/tradingBot/recurring/sub-orders",
			query:  "algoId=1",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotRecurringSubOrdersService().AlgoId("1").Do(context.Background())
				return err
			},
		},

		// TradingBot Signal - private
		{
			name:   "signal_create_signal",
			method: http.MethodPost,
			path:   "/api/v5/tradingBot/signal/create-signal",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotSignalCreateSignalService().SignalChanName("sig").Do(context.Background())
				return err
			},
		},
		{
			name:   "signal_signals",
			method: http.MethodGet,
			path:   "/api/v5/tradingBot/signal/signals",
			query:  "signalSourceType=1",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotSignalSignalsService().SignalSourceType("1").Do(context.Background())
				return err
			},
		},
		{
			name:   "signal_order_algo",
			method: http.MethodPost,
			path:   "/api/v5/tradingBot/signal/order-algo",
			body:   `{"signalChanId":"1","includeAll":true,"lever":"1","investAmt":"1","subOrdType":"1"}`,
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotSignalOrderAlgoService().
					SignalChanId("1").
					IncludeAll(true).
					Lever("1").
					InvestAmt("1").
					SubOrdType("1").
					Do(context.Background())
				return err
			},
		},
		{
			name:   "signal_stop_order_algo",
			method: http.MethodPost,
			path:   "/api/v5/tradingBot/signal/stop-order-algo",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotSignalStopOrderAlgoService().
					Orders([]TradingBotSignalStopOrder{{AlgoId: "1"}}).
					Do(context.Background())
				return err
			},
		},
		{
			name:   "signal_margin_balance",
			method: http.MethodPost,
			path:   "/api/v5/tradingBot/signal/margin-balance",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotSignalMarginBalanceService().AlgoId("1").Type("add").Amt("1").Do(context.Background())
				return err
			},
		},
		{
			name:   "signal_set_instruments",
			method: http.MethodPost,
			path:   "/api/v5/tradingBot/signal/set-instruments",
			body:   `{"algoId":"1","instIds":["BTC-USDT-SWAP"],"includeAll":false}`,
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotSignalSetInstrumentsService().
					AlgoId("1").
					InstIds([]string{"BTC-USDT-SWAP"}).
					IncludeAll(false).
					Do(context.Background())
				return err
			},
		},
		{
			name:   "signal_amend_tpsl",
			method: http.MethodPost,
			path:   "/api/v5/tradingBot/signal/amendTPSL",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotSignalAmendTPSLService().
					AlgoId("1").
					ExitSettingParam(TradingBotSignalExitSettingParam{TpSlType: "price"}).
					Do(context.Background())
				return err
			},
		},
		{
			name:   "signal_orders_algo_details",
			method: http.MethodGet,
			path:   "/api/v5/tradingBot/signal/orders-algo-details",
			query:  "algoId=1&algoOrdType=contract",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotSignalOrdersAlgoDetailsService().AlgoOrdType("contract").AlgoId("1").Do(context.Background())
				return err
			},
		},
		{
			name:   "signal_orders_algo_pending",
			method: http.MethodGet,
			path:   "/api/v5/tradingBot/signal/orders-algo-pending",
			query:  "algoOrdType=contract",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotSignalOrdersAlgoPendingService().AlgoOrdType("contract").Do(context.Background())
				return err
			},
		},
		{
			name:   "signal_orders_algo_history",
			method: http.MethodGet,
			path:   "/api/v5/tradingBot/signal/orders-algo-history",
			query:  "algoId=1&algoOrdType=contract",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotSignalOrdersAlgoHistoryService().AlgoOrdType("contract").AlgoId("1").Do(context.Background())
				return err
			},
		},
		{
			name:   "signal_positions",
			method: http.MethodGet,
			path:   "/api/v5/tradingBot/signal/positions",
			query:  "algoId=1&algoOrdType=contract",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotSignalPositionsService().AlgoOrdType("contract").AlgoId("1").Do(context.Background())
				return err
			},
		},
		{
			name:   "signal_positions_history",
			method: http.MethodGet,
			path:   "/api/v5/tradingBot/signal/positions-history",
			query:  "algoId=1",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotSignalPositionsHistoryService().AlgoId("1").Do(context.Background())
				return err
			},
		},
		{
			name:   "signal_sub_orders",
			method: http.MethodGet,
			path:   "/api/v5/tradingBot/signal/sub-orders",
			query:  "algoId=1&algoOrdType=contract&state=filled",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotSignalSubOrdersService().AlgoId("1").AlgoOrdType("contract").State("filled").Do(context.Background())
				return err
			},
		},
		{
			name:   "signal_event_history",
			method: http.MethodGet,
			path:   "/api/v5/tradingBot/signal/event-history",
			query:  "algoId=1",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotSignalEventHistoryService().AlgoId("1").Do(context.Background())
				return err
			},
		},
		{
			name:   "signal_close_position",
			method: http.MethodPost,
			path:   "/api/v5/tradingBot/signal/close-position",
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotSignalClosePositionService().AlgoId("1").InstId("BTC-USDT-SWAP").Do(context.Background())
				return err
			},
		},
		{
			name:     "signal_sub_order",
			method:   http.MethodPost,
			path:     "/api/v5/tradingBot/signal/sub-order",
			body:     `{"instId":"BTC-USDT-SWAP","algoId":"1","side":"buy","ordType":"market","sz":"1"}`,
			signed:   true,
			response: `{"code":"0","msg":"","data":[]}`,
			invokeDo: func(c *Client) error {
				return c.NewTradingBotSignalSubOrderService().
					AlgoId("1").
					InstId("BTC-USDT-SWAP").
					Side("buy").
					OrdType("market").
					Sz("1").
					Do(context.Background())
			},
		},
		{
			name:   "signal_cancel_sub_order",
			method: http.MethodPost,
			path:   "/api/v5/tradingBot/signal/cancel-sub-order",
			body:   `{"algoId":"1","instId":"BTC-USDT-SWAP","signalOrdId":"O1"}`,
			signed: true,
			invokeDo: func(c *Client) error {
				_, err := c.NewTradingBotSignalCancelSubOrderService().AlgoId("1").InstId("BTC-USDT-SWAP").SignalOrdId("O1").Do(context.Background())
				return err
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if got, want := r.Method, tc.method; got != want {
					t.Fatalf("method = %q, want %q", got, want)
				}
				if got, want := r.URL.Path, tc.path; got != want {
					t.Fatalf("path = %q, want %q", got, want)
				}
				if got, want := r.URL.RawQuery, tc.query; got != want {
					t.Fatalf("query = %q, want %q", got, want)
				}

				if tc.body != "" {
					bodyBytes, _ := io.ReadAll(r.Body)
					if got, want := string(bodyBytes), tc.body; got != want {
						t.Fatalf("body = %q, want %q", got, want)
					}
				}

				if tc.signed {
					if got, want := r.Header.Get("OK-ACCESS-KEY"), "mykey"; got != want {
						t.Fatalf("OK-ACCESS-KEY = %q, want %q", got, want)
					}
					if got, want := r.Header.Get("OK-ACCESS-PASSPHRASE"), "mypass"; got != want {
						t.Fatalf("OK-ACCESS-PASSPHRASE = %q, want %q", got, want)
					}
					if got := r.Header.Get("OK-ACCESS-TIMESTAMP"); got == "" {
						t.Fatalf("OK-ACCESS-TIMESTAMP empty")
					}
					if got := r.Header.Get("OK-ACCESS-SIGN"); got == "" {
						t.Fatalf("OK-ACCESS-SIGN empty")
					}
				} else if got := r.Header.Get("OK-ACCESS-KEY"); got != "" {
					t.Fatalf("unexpected signed header OK-ACCESS-KEY = %q", got)
				}

				w.Header().Set("Content-Type", "application/json")
				resp := tc.response
				if resp == "" {
					resp = okResp
				}
				_, _ = w.Write([]byte(resp))
			}))
			t.Cleanup(srv.Close)

			c := mkClient(srv)
			if err := tc.invokeDo(c); err != nil {
				t.Fatalf("Do() error = %v", err)
			}
		})
	}
}

func TestTradingBotGridAmendAlgoBasicParamService_Do_DataCompat(t *testing.T) {
	fixedNow := time.Date(2020, 6, 30, 12, 34, 56, 789_000_000, time.UTC)

	t.Run("data_object", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":{"algoId":"1","requiredTopupAmount":"0"}}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
			WithCredentials(Credentials{APIKey: "mykey", SecretKey: "mysecret", Passphrase: "mypass"}),
			WithNowFunc(func() time.Time { return fixedNow }),
		)

		_, err := c.NewTradingBotGridAmendAlgoBasicParamService().AlgoId("1").MinPx("0").MaxPx("1").GridNum("10").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
	})

	t.Run("data_array", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":"0","msg":"","data":[{"algoId":"1","requiredTopupAmount":"0"}]}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient(
			WithBaseURL(srv.URL),
			WithHTTPClient(srv.Client()),
			WithCredentials(Credentials{APIKey: "mykey", SecretKey: "mysecret", Passphrase: "mypass"}),
			WithNowFunc(func() time.Time { return fixedNow }),
		)

		_, err := c.NewTradingBotGridAmendAlgoBasicParamService().AlgoId("1").MinPx("0").MaxPx("1").GridNum("10").Do(context.Background())
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
	})
}

func TestTradingBotSignalSubOrdersService_Do_Validation(t *testing.T) {
	c := NewClient()
	_, err := c.NewTradingBotSignalSubOrdersService().AlgoId("1").AlgoOrdType("contract").Do(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
	if err != errTradingBotSignalSubOrdersMissingStateOrSignalOrdId {
		t.Fatalf("error = %v, want %v", err, errTradingBotSignalSubOrdersMissingStateOrSignalOrdId)
	}
}

func TestTradingBotSignalSubOrderService_Do_Validation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := io.ReadAll(r.Body)
		if got, want := string(bodyBytes), `{"instId":"BTC-USDT-SWAP","algoId":"1","side":"buy","ordType":"limit","sz":"1"}`; got != want {
			t.Fatalf("body = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":"51000","msg":"invalid","data":[]}`))
	}))
	t.Cleanup(srv.Close)

	c := NewClient(
		WithBaseURL(srv.URL),
		WithHTTPClient(srv.Client()),
		WithCredentials(Credentials{APIKey: "mykey", SecretKey: "mysecret", Passphrase: "mypass"}),
	)

	err := c.NewTradingBotSignalSubOrderService().
		AlgoId("1").
		InstId("BTC-USDT-SWAP").
		Side("buy").
		OrdType("limit").
		Sz("1").
		Do(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("err = %T, want *APIError: %v", err, err)
	}
	if apiErr.Code != "51000" {
		t.Fatalf("apiErr.Code = %q, want %q", apiErr.Code, "51000")
	}
}
