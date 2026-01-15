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
		log.Fatal("refusing to amend tradingBot grid basic params; set OKX_CONFIRM=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	algoId := os.Getenv("OKX_ALGO_ID")
	minPx := os.Getenv("OKX_MIN_PX")
	maxPx := os.Getenv("OKX_MAX_PX")
	gridNum := os.Getenv("OKX_GRID_NUM")
	if algoId == "" || minPx == "" || maxPx == "" || gridNum == "" {
		log.Fatal("missing env: OKX_ALGO_ID / OKX_MIN_PX / OKX_MAX_PX / OKX_GRID_NUM")
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

	res, err := c.NewTradingBotGridAmendAlgoBasicParamService().
		AlgoId(algoId).
		MinPx(minPx).
		MaxPx(maxPx).
		GridNum(gridNum).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("tradingBot grid amend-algo-basic-param: algoId=%s requiredTopupAmount=%s", res.AlgoId, res.RequiredTopupAmount)
}
