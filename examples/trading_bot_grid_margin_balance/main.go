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
		log.Fatal("refusing to adjust tradingBot grid margin balance; set OKX_CONFIRM=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	algoId := os.Getenv("OKX_ALGO_ID")
	typ := os.Getenv("OKX_TYPE")
	amt := os.Getenv("OKX_AMT")
	percent := os.Getenv("OKX_PERCENT")
	if algoId == "" || typ == "" {
		log.Fatal("missing env: OKX_ALGO_ID / OKX_TYPE")
	}
	if amt == "" && percent == "" {
		log.Fatal("missing env: OKX_AMT or OKX_PERCENT")
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

	svc := c.NewTradingBotGridMarginBalanceService().AlgoId(algoId).Type(typ)
	if amt != "" {
		svc.Amt(amt)
	}
	if percent != "" {
		svc.Percent(percent)
	}

	ack, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("tradingBot grid margin-balance: algoId=%s algoClOrdId=%s sCode=%s sMsg=%s", ack.AlgoId, ack.AlgoClOrdId, ack.SCode, ack.SMsg)
}
