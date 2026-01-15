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
		log.Fatal("refusing to place tradingBot signal sub-order; set OKX_CONFIRM=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	algoId := os.Getenv("OKX_ALGO_ID")
	if algoId == "" {
		log.Fatal("missing env: OKX_ALGO_ID")
	}
	instId := os.Getenv("OKX_INST_ID")
	if instId == "" {
		log.Fatal("missing env: OKX_INST_ID")
	}
	side := os.Getenv("OKX_SIDE")
	if side == "" {
		log.Fatal("missing env: OKX_SIDE (buy/sell)")
	}
	sz := os.Getenv("OKX_SZ")
	if sz == "" {
		log.Fatal("missing env: OKX_SZ")
	}
	ordType := os.Getenv("OKX_ORD_TYPE")
	if ordType == "" {
		ordType = "limit"
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

	svc := c.NewTradingBotSignalSubOrderService().
		AlgoId(algoId).
		InstId(instId).
		Side(side).
		OrdType(ordType).
		Sz(sz)

	if ordType == "limit" {
		px := os.Getenv("OKX_PX")
		if px == "" {
			log.Fatal("missing env: OKX_PX (required for limit)")
		}
		svc.Px(px)
	}
	if os.Getenv("OKX_REDUCE_ONLY") == "true" {
		svc.ReduceOnly(true)
	}

	if err := svc.Do(context.Background()); err != nil {
		log.Fatal(err)
	}

	log.Printf("tradingBot signal sub-order submitted: algoId=%s instId=%s side=%s ordType=%s sz=%s", algoId, instId, side, ordType, sz)
}
