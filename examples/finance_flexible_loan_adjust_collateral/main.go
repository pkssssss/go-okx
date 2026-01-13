package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to adjust flexible-loan collateral; set OKX_CONFIRM=YES to continue")
	}

	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}

	typ := os.Getenv("OKX_TYPE")
	if typ == "" {
		log.Fatal("missing env: OKX_TYPE (add/reduce)")
	}
	collateralCcy := os.Getenv("OKX_COLLATERAL_CCY")
	if collateralCcy == "" {
		log.Fatal("missing env: OKX_COLLATERAL_CCY")
	}
	collateralAmt := os.Getenv("OKX_COLLATERAL_AMT")
	if collateralAmt == "" {
		log.Fatal("missing env: OKX_COLLATERAL_AMT")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	c := okx.NewClient(
		okx.WithCredentials(okx.Credentials{APIKey: apiKey, SecretKey: secretKey, Passphrase: passphrase}),
		okx.WithDemoTrading(demo),
	)

	if _, err := c.SyncTime(context.Background()); err != nil {
		log.Fatal(err)
	}

	if err := c.NewFinanceFlexibleLoanAdjustCollateralService().
		Type(typ).
		CollateralCcy(collateralCcy).
		CollateralAmt(collateralAmt).
		Do(context.Background()); err != nil {
		log.Fatal(err)
	}

	log.Printf("flexible-loan adjust-collateral accepted")
}
