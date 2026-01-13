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

func TestFinanceServices_Do(t *testing.T) {
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
			name:     "flexible_loan_borrow_currencies",
			method:   http.MethodGet,
			path:     "/api/v5/finance/flexible-loan/borrow-currencies",
			rawQuery: "",
			signed:   true,
			call: func(c *Client) error {
				_, err := c.NewFinanceFlexibleLoanBorrowCurrenciesService().Do(context.Background())
				return err
			},
		},
		{
			name:     "flexible_loan_collateral_assets",
			method:   http.MethodGet,
			path:     "/api/v5/finance/flexible-loan/collateral-assets",
			rawQuery: "ccy=BTC",
			signed:   true,
			call: func(c *Client) error {
				_, err := c.NewFinanceFlexibleLoanCollateralAssetsService().Ccy("BTC").Do(context.Background())
				return err
			},
		},
		{
			name:     "flexible_loan_max_loan",
			method:   http.MethodPost,
			path:     "/api/v5/finance/flexible-loan/max-loan",
			rawQuery: "",
			signed:   true,
			bodyContains: []string{
				`"borrowCcy":"USDT"`,
				`"supCollateral":[`,
				`"ccy":"BTC"`,
				`"amt":"0.1"`,
			},
			call: func(c *Client) error {
				_, err := c.NewFinanceFlexibleLoanMaxLoanService().
					BorrowCcy("USDT").
					SupCollateral([]FinanceFlexibleLoanSupCollateral{{Ccy: "BTC", Amt: "0.1"}}).
					Do(context.Background())
				return err
			},
		},
		{
			name:     "flexible_loan_max_collateral_redeem_amount",
			method:   http.MethodGet,
			path:     "/api/v5/finance/flexible-loan/max-collateral-redeem-amount",
			rawQuery: "ccy=USDT",
			signed:   true,
			call: func(c *Client) error {
				_, err := c.NewFinanceFlexibleLoanMaxCollateralRedeemAmountService().Ccy("USDT").Do(context.Background())
				return err
			},
		},
		{
			name:     "flexible_loan_adjust_collateral",
			method:   http.MethodPost,
			path:     "/api/v5/finance/flexible-loan/adjust-collateral",
			rawQuery: "",
			signed:   true,
			bodyContains: []string{
				`"type":"add"`,
				`"collateralCcy":"BTC"`,
				`"collateralAmt":"0.1"`,
			},
			call: func(c *Client) error {
				return c.NewFinanceFlexibleLoanAdjustCollateralService().
					Type("add").
					CollateralCcy("BTC").
					CollateralAmt("0.1").
					Do(context.Background())
			},
		},
		{
			name:     "flexible_loan_loan_info",
			method:   http.MethodGet,
			path:     "/api/v5/finance/flexible-loan/loan-info",
			rawQuery: "",
			signed:   true,
			call: func(c *Client) error {
				_, err := c.NewFinanceFlexibleLoanLoanInfoService().Do(context.Background())
				return err
			},
		},
		{
			name:     "flexible_loan_loan_history",
			method:   http.MethodGet,
			path:     "/api/v5/finance/flexible-loan/loan-history",
			rawQuery: "after=1&before=2&limit=100&type=borrowed",
			signed:   true,
			call: func(c *Client) error {
				_, err := c.NewFinanceFlexibleLoanLoanHistoryService().
					Type("borrowed").
					After("1").
					Before("2").
					Limit(100).
					Do(context.Background())
				return err
			},
		},
		{
			name:     "flexible_loan_interest_accrued",
			method:   http.MethodGet,
			path:     "/api/v5/finance/flexible-loan/interest-accrued",
			rawQuery: "after=1&before=2&ccy=USDT&limit=100",
			signed:   true,
			call: func(c *Client) error {
				_, err := c.NewFinanceFlexibleLoanInterestAccruedService().
					Ccy("USDT").
					After("1").
					Before("2").
					Limit(100).
					Do(context.Background())
				return err
			},
		},

		{
			name:     "savings_balance",
			method:   http.MethodGet,
			path:     "/api/v5/finance/savings/balance",
			rawQuery: "ccy=USDT",
			signed:   true,
			call: func(c *Client) error {
				_, err := c.NewFinanceSavingsBalanceService().Ccy("USDT").Do(context.Background())
				return err
			},
		},
		{
			name:     "savings_purchase_redempt",
			method:   http.MethodPost,
			path:     "/api/v5/finance/savings/purchase-redempt",
			rawQuery: "",
			signed:   true,
			bodyContains: []string{
				`"ccy":"USDT"`,
				`"amt":"0.1"`,
				`"side":"purchase"`,
				`"rate":"0.02"`,
			},
			call: func(c *Client) error {
				_, err := c.NewFinanceSavingsPurchaseRedemptService().
					Ccy("USDT").
					Amt("0.1").
					Side("purchase").
					Rate("0.02").
					Do(context.Background())
				return err
			},
		},
		{
			name:     "savings_set_lending_rate",
			method:   http.MethodPost,
			path:     "/api/v5/finance/savings/set-lending-rate",
			rawQuery: "",
			signed:   true,
			bodyContains: []string{
				`"ccy":"BTC"`,
				`"rate":"0.02"`,
			},
			call: func(c *Client) error {
				_, err := c.NewFinanceSavingsSetLendingRateService().Ccy("BTC").Rate("0.02").Do(context.Background())
				return err
			},
		},
		{
			name:     "savings_lending_history",
			method:   http.MethodGet,
			path:     "/api/v5/finance/savings/lending-history",
			rawQuery: "after=1&before=2&ccy=BTC&limit=10",
			signed:   true,
			call: func(c *Client) error {
				_, err := c.NewFinanceSavingsLendingHistoryService().
					Ccy("BTC").
					After("1").
					Before("2").
					Limit(10).
					Do(context.Background())
				return err
			},
		},
		{
			name:     "savings_lending_rate_summary_public",
			method:   http.MethodGet,
			path:     "/api/v5/finance/savings/lending-rate-summary",
			rawQuery: "ccy=BTC",
			signed:   false,
			call: func(c *Client) error {
				_, err := c.NewFinanceSavingsLendingRateSummaryService().Ccy("BTC").Do(context.Background())
				return err
			},
		},
		{
			name:     "savings_lending_rate_history_public",
			method:   http.MethodGet,
			path:     "/api/v5/finance/savings/lending-rate-history",
			rawQuery: "after=1&before=2&ccy=BTC&limit=10",
			signed:   false,
			call: func(c *Client) error {
				_, err := c.NewFinanceSavingsLendingRateHistoryService().
					Ccy("BTC").
					After("1").
					Before("2").
					Limit(10).
					Do(context.Background())
				return err
			},
		},

		{
			name:     "staking_defi_offers",
			method:   http.MethodGet,
			path:     "/api/v5/finance/staking-defi/offers",
			rawQuery: "ccy=USDT&productId=123&protocolType=defi",
			signed:   true,
			call: func(c *Client) error {
				_, err := c.NewFinanceStakingDefiOffersService().
					ProductId("123").
					ProtocolType("defi").
					Ccy("USDT").
					Do(context.Background())
				return err
			},
		},
		{
			name:     "staking_defi_purchase",
			method:   http.MethodPost,
			path:     "/api/v5/finance/staking-defi/purchase",
			rawQuery: "",
			signed:   true,
			bodyContains: []string{
				`"productId":"123"`,
				`"investData":[`,
				`"ccy":"USDT"`,
				`"amt":"100"`,
				`"term":"30"`,
				`"tag":"TAG"`,
			},
			call: func(c *Client) error {
				_, err := c.NewFinanceStakingDefiPurchaseService().
					ProductId("123").
					InvestData([]FinanceStakingDefiInvestData{{Ccy: "USDT", Amt: "100"}}).
					Term("30").
					Tag("TAG").
					Do(context.Background())
				return err
			},
		},
		{
			name:     "staking_defi_redeem",
			method:   http.MethodPost,
			path:     "/api/v5/finance/staking-defi/redeem",
			rawQuery: "",
			signed:   true,
			bodyContains: []string{
				`"ordId":"OID"`,
				`"protocolType":"defi"`,
				`"allowEarlyRedeem":true`,
			},
			call: func(c *Client) error {
				_, err := c.NewFinanceStakingDefiRedeemService().
					OrdId("OID").
					ProtocolType("defi").
					AllowEarlyRedeem(true).
					Do(context.Background())
				return err
			},
		},
		{
			name:     "staking_defi_cancel",
			method:   http.MethodPost,
			path:     "/api/v5/finance/staking-defi/cancel",
			rawQuery: "",
			signed:   true,
			bodyContains: []string{
				`"ordId":"OID"`,
				`"protocolType":"defi"`,
			},
			call: func(c *Client) error {
				_, err := c.NewFinanceStakingDefiCancelService().
					OrdId("OID").
					ProtocolType("defi").
					Do(context.Background())
				return err
			},
		},
		{
			name:     "staking_defi_orders_active",
			method:   http.MethodGet,
			path:     "/api/v5/finance/staking-defi/orders-active",
			rawQuery: "ccy=USDT&productId=123&protocolType=defi&state=1",
			signed:   true,
			call: func(c *Client) error {
				_, err := c.NewFinanceStakingDefiOrdersActiveService().
					ProductId("123").
					ProtocolType("defi").
					Ccy("USDT").
					State("1").
					Do(context.Background())
				return err
			},
		},
		{
			name:     "staking_defi_orders_history",
			method:   http.MethodGet,
			path:     "/api/v5/finance/staking-defi/orders-history",
			rawQuery: "after=1&before=2&ccy=USDT&limit=10&productId=123&protocolType=defi",
			signed:   true,
			call: func(c *Client) error {
				_, err := c.NewFinanceStakingDefiOrdersHistoryService().
					ProductId("123").
					ProtocolType("defi").
					Ccy("USDT").
					After("1").
					Before("2").
					Limit(10).
					Do(context.Background())
				return err
			},
		},

		{
			name:     "staking_defi_eth_product_info",
			method:   http.MethodGet,
			path:     "/api/v5/finance/staking-defi/eth/product-info",
			rawQuery: "",
			signed:   true,
			call: func(c *Client) error {
				_, err := c.NewFinanceStakingDefiETHProductInfoService().Do(context.Background())
				return err
			},
		},
		{
			name:     "staking_defi_eth_purchase",
			method:   http.MethodPost,
			path:     "/api/v5/finance/staking-defi/eth/purchase",
			rawQuery: "",
			signed:   true,
			bodyContains: []string{
				`"amt":"1"`,
			},
			call: func(c *Client) error {
				return c.NewFinanceStakingDefiETHPurchaseService().Amt("1").Do(context.Background())
			},
		},
		{
			name:     "staking_defi_eth_redeem",
			method:   http.MethodPost,
			path:     "/api/v5/finance/staking-defi/eth/redeem",
			rawQuery: "",
			signed:   true,
			bodyContains: []string{
				`"amt":"1"`,
			},
			call: func(c *Client) error {
				return c.NewFinanceStakingDefiETHRedeemService().Amt("1").Do(context.Background())
			},
		},
		{
			name:     "staking_defi_eth_cancel_redeem",
			method:   http.MethodPost,
			path:     "/api/v5/finance/staking-defi/eth/cancel-redeem",
			rawQuery: "",
			signed:   true,
			bodyContains: []string{
				`"ordId":"OID"`,
			},
			call: func(c *Client) error {
				_, err := c.NewFinanceStakingDefiETHCancelRedeemService().OrdId("OID").Do(context.Background())
				return err
			},
		},
		{
			name:     "staking_defi_eth_balance",
			method:   http.MethodGet,
			path:     "/api/v5/finance/staking-defi/eth/balance",
			rawQuery: "",
			signed:   true,
			call: func(c *Client) error {
				_, err := c.NewFinanceStakingDefiETHBalanceService().Do(context.Background())
				return err
			},
		},
		{
			name:     "staking_defi_eth_purchase_redeem_history",
			method:   http.MethodGet,
			path:     "/api/v5/finance/staking-defi/eth/purchase-redeem-history",
			rawQuery: "after=1&before=2&limit=10&status=success&type=purchase",
			signed:   true,
			call: func(c *Client) error {
				_, err := c.NewFinanceStakingDefiETHPurchaseRedeemHistoryService().
					Type("purchase").
					Status("success").
					After("1").
					Before("2").
					Limit(10).
					Do(context.Background())
				return err
			},
		},
		{
			name:     "staking_defi_eth_apy_history_public",
			method:   http.MethodGet,
			path:     "/api/v5/finance/staking-defi/eth/apy-history",
			rawQuery: "days=7",
			signed:   false,
			call: func(c *Client) error {
				_, err := c.NewFinanceStakingDefiETHAPYHistoryService().Days("7").Do(context.Background())
				return err
			},
		},

		{
			name:     "staking_defi_sol_product_info",
			method:   http.MethodGet,
			path:     "/api/v5/finance/staking-defi/sol/product-info",
			rawQuery: "",
			signed:   true,
			call: func(c *Client) error {
				_, err := c.NewFinanceStakingDefiSOLProductInfoService().Do(context.Background())
				return err
			},
		},
		{
			name:     "staking_defi_sol_purchase",
			method:   http.MethodPost,
			path:     "/api/v5/finance/staking-defi/sol/purchase",
			rawQuery: "",
			signed:   true,
			bodyContains: []string{
				`"amt":"1"`,
			},
			call: func(c *Client) error {
				return c.NewFinanceStakingDefiSOLPurchaseService().Amt("1").Do(context.Background())
			},
		},
		{
			name:     "staking_defi_sol_redeem",
			method:   http.MethodPost,
			path:     "/api/v5/finance/staking-defi/sol/redeem",
			rawQuery: "",
			signed:   true,
			bodyContains: []string{
				`"amt":"1"`,
			},
			call: func(c *Client) error {
				return c.NewFinanceStakingDefiSOLRedeemService().Amt("1").Do(context.Background())
			},
		},
		{
			name:     "staking_defi_sol_balance",
			method:   http.MethodGet,
			path:     "/api/v5/finance/staking-defi/sol/balance",
			rawQuery: "",
			signed:   true,
			call: func(c *Client) error {
				_, err := c.NewFinanceStakingDefiSOLBalanceService().Do(context.Background())
				return err
			},
		},
		{
			name:     "staking_defi_sol_purchase_redeem_history",
			method:   http.MethodGet,
			path:     "/api/v5/finance/staking-defi/sol/purchase-redeem-history",
			rawQuery: "after=1&before=2&limit=10&status=success&type=purchase",
			signed:   true,
			call: func(c *Client) error {
				_, err := c.NewFinanceStakingDefiSOLPurchaseRedeemHistoryService().
					Type("purchase").
					Status("success").
					After("1").
					Before("2").
					Limit(10).
					Do(context.Background())
				return err
			},
		},
		{
			name:     "staking_defi_sol_apy_history_public",
			method:   http.MethodGet,
			path:     "/api/v5/finance/staking-defi/sol/apy-history",
			rawQuery: "days=7",
			signed:   false,
			call: func(c *Client) error {
				_, err := c.NewFinanceStakingDefiSOLAPYHistoryService().Days("7").Do(context.Background())
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
				if tt.path == "/api/v5/finance/staking-defi/sol/product-info" {
					_, _ = w.Write([]byte(`{"code":"0","msg":"","data":{}}`))
					return
				}
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

func TestFinanceServices_Validation(t *testing.T) {
	c := NewClient()

	t.Run("flexible_loan_max_loan_missing_borrowCcy", func(t *testing.T) {
		_, err := c.NewFinanceFlexibleLoanMaxLoanService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errFinanceFlexibleLoanMaxLoanMissingBorrowCcy {
			t.Fatalf("error = %v, want %v", err, errFinanceFlexibleLoanMaxLoanMissingBorrowCcy)
		}
	})

	t.Run("flexible_loan_max_loan_invalid_supCollateral", func(t *testing.T) {
		_, err := c.NewFinanceFlexibleLoanMaxLoanService().
			BorrowCcy("USDT").
			SupCollateral([]FinanceFlexibleLoanSupCollateral{{Ccy: "BTC"}}).
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errFinanceFlexibleLoanMaxLoanInvalidSupCollateral {
			t.Fatalf("error = %v, want %v", err, errFinanceFlexibleLoanMaxLoanInvalidSupCollateral)
		}
	})

	t.Run("flexible_loan_adjust_collateral_missing_required", func(t *testing.T) {
		err := c.NewFinanceFlexibleLoanAdjustCollateralService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errFinanceFlexibleLoanAdjustCollateralMissingRequired {
			t.Fatalf("error = %v, want %v", err, errFinanceFlexibleLoanAdjustCollateralMissingRequired)
		}
	})

	t.Run("flexible_loan_max_collateral_redeem_amount_missing_ccy", func(t *testing.T) {
		_, err := c.NewFinanceFlexibleLoanMaxCollateralRedeemAmountService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errFinanceFlexibleLoanMaxCollateralRedeemAmountMissingCcy {
			t.Fatalf("error = %v, want %v", err, errFinanceFlexibleLoanMaxCollateralRedeemAmountMissingCcy)
		}
	})

	t.Run("savings_purchase_redempt_missing_required", func(t *testing.T) {
		_, err := c.NewFinanceSavingsPurchaseRedemptService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errFinanceSavingsPurchaseRedemptMissingRequired {
			t.Fatalf("error = %v, want %v", err, errFinanceSavingsPurchaseRedemptMissingRequired)
		}
	})

	t.Run("savings_set_lending_rate_missing_required", func(t *testing.T) {
		_, err := c.NewFinanceSavingsSetLendingRateService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errFinanceSavingsSetLendingRateMissingRequired {
			t.Fatalf("error = %v, want %v", err, errFinanceSavingsSetLendingRateMissingRequired)
		}
	})

	t.Run("staking_defi_purchase_missing_productId", func(t *testing.T) {
		_, err := c.NewFinanceStakingDefiPurchaseService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errFinanceStakingDefiPurchaseMissingProductId {
			t.Fatalf("error = %v, want %v", err, errFinanceStakingDefiPurchaseMissingProductId)
		}
	})

	t.Run("staking_defi_purchase_missing_investData", func(t *testing.T) {
		_, err := c.NewFinanceStakingDefiPurchaseService().ProductId("123").Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errFinanceStakingDefiPurchaseMissingInvestData {
			t.Fatalf("error = %v, want %v", err, errFinanceStakingDefiPurchaseMissingInvestData)
		}
	})

	t.Run("staking_defi_purchase_invalid_investData", func(t *testing.T) {
		_, err := c.NewFinanceStakingDefiPurchaseService().
			ProductId("123").
			InvestData([]FinanceStakingDefiInvestData{{Ccy: "USDT"}}).
			Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errFinanceStakingDefiPurchaseInvalidInvestData {
			t.Fatalf("error = %v, want %v", err, errFinanceStakingDefiPurchaseInvalidInvestData)
		}
	})

	t.Run("staking_defi_redeem_missing_required", func(t *testing.T) {
		_, err := c.NewFinanceStakingDefiRedeemService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errFinanceStakingDefiRedeemMissingRequired {
			t.Fatalf("error = %v, want %v", err, errFinanceStakingDefiRedeemMissingRequired)
		}
	})

	t.Run("staking_defi_cancel_missing_required", func(t *testing.T) {
		_, err := c.NewFinanceStakingDefiCancelService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errFinanceStakingDefiCancelMissingRequired {
			t.Fatalf("error = %v, want %v", err, errFinanceStakingDefiCancelMissingRequired)
		}
	})

	t.Run("staking_defi_eth_purchase_missing_amt", func(t *testing.T) {
		err := c.NewFinanceStakingDefiETHPurchaseService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errFinanceStakingDefiETHPurchaseMissingAmt {
			t.Fatalf("error = %v, want %v", err, errFinanceStakingDefiETHPurchaseMissingAmt)
		}
	})

	t.Run("staking_defi_eth_redeem_missing_amt", func(t *testing.T) {
		err := c.NewFinanceStakingDefiETHRedeemService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errFinanceStakingDefiETHRedeemMissingAmt {
			t.Fatalf("error = %v, want %v", err, errFinanceStakingDefiETHRedeemMissingAmt)
		}
	})

	t.Run("staking_defi_eth_cancel_redeem_missing_ordId", func(t *testing.T) {
		_, err := c.NewFinanceStakingDefiETHCancelRedeemService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errFinanceStakingDefiETHCancelRedeemMissingOrdId {
			t.Fatalf("error = %v, want %v", err, errFinanceStakingDefiETHCancelRedeemMissingOrdId)
		}
	})

	t.Run("staking_defi_eth_apy_history_missing_days", func(t *testing.T) {
		_, err := c.NewFinanceStakingDefiETHAPYHistoryService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errFinanceStakingDefiETHAPYHistoryMissingDays {
			t.Fatalf("error = %v, want %v", err, errFinanceStakingDefiETHAPYHistoryMissingDays)
		}
	})

	t.Run("staking_defi_sol_purchase_missing_amt", func(t *testing.T) {
		err := c.NewFinanceStakingDefiSOLPurchaseService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errFinanceStakingDefiSOLPurchaseMissingAmt {
			t.Fatalf("error = %v, want %v", err, errFinanceStakingDefiSOLPurchaseMissingAmt)
		}
	})

	t.Run("staking_defi_sol_redeem_missing_amt", func(t *testing.T) {
		err := c.NewFinanceStakingDefiSOLRedeemService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errFinanceStakingDefiSOLRedeemMissingAmt {
			t.Fatalf("error = %v, want %v", err, errFinanceStakingDefiSOLRedeemMissingAmt)
		}
	})

	t.Run("staking_defi_sol_apy_history_missing_days", func(t *testing.T) {
		_, err := c.NewFinanceStakingDefiSOLAPYHistoryService().Do(context.Background())
		if err == nil {
			t.Fatalf("expected error")
		}
		if err != errFinanceStakingDefiSOLAPYHistoryMissingDays {
			t.Fatalf("error = %v, want %v", err, errFinanceStakingDefiSOLAPYHistoryMissingDays)
		}
	})
}
