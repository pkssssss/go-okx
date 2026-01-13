package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to amend profit sharing ratio; set OKX_CONFIRM=YES to continue")
	}

	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}

	profitSharingRatio := os.Getenv("OKX_PROFIT_SHARING_RATIO")
	if profitSharingRatio == "" {
		log.Fatal("missing env: OKX_PROFIT_SHARING_RATIO")
	}

	instType := os.Getenv("OKX_INST_TYPE")
	if instType == "" {
		instType = "SWAP"
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	c := okx.NewClient(
		okx.WithCredentials(okx.Credentials{
			APIKey:     apiKey,
			SecretKey:  secretKey,
			Passphrase: passphrase,
		}),
		okx.WithDemoTrading(demo),
	)

	if _, err := c.SyncTime(context.Background()); err != nil {
		log.Fatal(err)
	}

	res, err := c.NewCopyTradingAmendProfitSharingRatioService().
		InstType(instType).
		ProfitSharingRatio(profitSharingRatio).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("copytrading_amend_profit_sharing_ratio: result=%v", res.Result)
}
