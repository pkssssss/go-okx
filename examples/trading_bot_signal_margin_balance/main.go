package main

import (
	"context"
	"log"
	"os"
	"strconv"

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
		log.Fatal("refusing to adjust tradingBot signal margin balance; set OKX_CONFIRM=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	algoId := os.Getenv("OKX_ALGO_ID")
	typ := os.Getenv("OKX_TYPE")
	amt := os.Getenv("OKX_AMT")
	if algoId == "" || typ == "" || amt == "" {
		log.Fatal("missing env: OKX_ALGO_ID / OKX_TYPE / OKX_AMT")
	}

	allowReinvestRaw := os.Getenv("OKX_ALLOW_REINVEST")
	var allowReinvest *bool
	if allowReinvestRaw != "" {
		b, err := strconv.ParseBool(allowReinvestRaw)
		if err != nil {
			log.Fatalf("invalid env OKX_ALLOW_REINVEST=%q: %v", allowReinvestRaw, err)
		}
		allowReinvest = &b
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

	svc := c.NewTradingBotSignalMarginBalanceService().AlgoId(algoId).Type(typ).Amt(amt)
	if allowReinvest != nil {
		svc.AllowReinvest(*allowReinvest)
	}

	ack, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("tradingBot signal margin-balance: algoId=%s", ack.AlgoId)
}
