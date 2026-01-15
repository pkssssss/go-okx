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

	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to compute tradingBot grid margin balance; set OKX_CONFIRM=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	algoId := os.Getenv("OKX_ALGO_ID")
	typ := os.Getenv("OKX_TYPE")
	amt := os.Getenv("OKX_AMT")
	if algoId == "" || typ == "" || amt == "" {
		log.Fatal("missing env: OKX_ALGO_ID / OKX_TYPE / OKX_AMT")
	}

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

	res, err := c.NewTradingBotGridComputeMarginBalanceService().AlgoId(algoId).Type(typ).Amt(amt).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("tradingBot grid compute-margin-balance: maxAmt=%s lever=%s", res.MaxAmt, res.Lever)
}
