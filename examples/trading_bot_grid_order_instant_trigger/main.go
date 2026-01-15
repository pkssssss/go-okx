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
		log.Fatal("refusing to trigger tradingBot grid order; set OKX_CONFIRM=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	algoId := os.Getenv("OKX_ALGO_ID")
	if algoId == "" {
		log.Fatal("missing env: OKX_ALGO_ID")
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

	svc := c.NewTradingBotGridOrderInstantTriggerService().AlgoId(algoId)
	if v := os.Getenv("OKX_TOP_UP_AMT"); v != "" {
		svc.TopUpAmt(v)
	}

	ack, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("tradingBot grid order-instant-trigger: algoId=%s algoClOrdId=%s sCode=%s sMsg=%s", ack.AlgoId, ack.AlgoClOrdId, ack.SCode, ack.SMsg)
}
