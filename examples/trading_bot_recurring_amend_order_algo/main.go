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
		log.Fatal("refusing to amend tradingBot recurring order; set OKX_CONFIRM=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	algoId := os.Getenv("OKX_ALGO_ID")
	stgyName := os.Getenv("OKX_STGY_NAME")
	if algoId == "" || stgyName == "" {
		log.Fatal("missing env: OKX_ALGO_ID / OKX_STGY_NAME")
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

	ack, err := c.NewTradingBotRecurringAmendOrderAlgoService().AlgoId(algoId).StgyName(stgyName).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("tradingBot recurring amend-order-algo: algoId=%s algoClOrdId=%s sCode=%s sMsg=%s", ack.AlgoId, ack.AlgoClOrdId, ack.SCode, ack.SMsg)
}
