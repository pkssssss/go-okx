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

	rl, err := c.NewTradeAccountRateLimitService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(
		"trade_account_rate_limit: accRateLimit=%s fillRatio=%s mainFillRatio=%s nextAccRateLimit=%s ts=%d",
		rl.AccRateLimit,
		rl.FillRatio,
		rl.MainFillRatio,
		rl.NextAccRateLimit,
		rl.TS,
	)
}
