package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}

	borrowCcy := os.Getenv("OKX_BORROW_CCY")
	if borrowCcy == "" {
		log.Fatal("missing env: OKX_BORROW_CCY")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	supCcy := os.Getenv("OKX_SUP_COLLATERAL_CCY")
	supAmt := os.Getenv("OKX_SUP_COLLATERAL_AMT")

	c := okx.NewClient(
		okx.WithCredentials(okx.Credentials{APIKey: apiKey, SecretKey: secretKey, Passphrase: passphrase}),
		okx.WithDemoTrading(demo),
	)

	if _, err := c.SyncTime(context.Background()); err != nil {
		log.Fatal(err)
	}

	svc := c.NewFinanceFlexibleLoanMaxLoanService().BorrowCcy(borrowCcy)
	if supCcy != "" || supAmt != "" {
		if supCcy == "" || supAmt == "" {
			log.Fatal("env must be both set: OKX_SUP_COLLATERAL_CCY and OKX_SUP_COLLATERAL_AMT")
		}
		svc.SupCollateral([]okx.FinanceFlexibleLoanSupCollateral{{Ccy: supCcy, Amt: supAmt}})
	}

	resp, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("borrowCcy=%s maxLoan=%s notionalUsd=%s remainingQuota=%s", resp.BorrowCcy, resp.MaxLoan, resp.NotionalUsd, resp.RemainingQuota)
}
