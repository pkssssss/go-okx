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
		log.Fatal("refusing to place tradingBot grid order; set OKX_CONFIRM=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	instId := os.Getenv("OKX_INST_ID")
	if instId == "" {
		log.Fatal("missing env: OKX_INST_ID")
	}
	algoOrdType := os.Getenv("OKX_ALGO_ORD_TYPE")
	if algoOrdType == "" {
		algoOrdType = "grid"
	}
	maxPx := os.Getenv("OKX_MAX_PX")
	if maxPx == "" {
		log.Fatal("missing env: OKX_MAX_PX")
	}
	minPx := os.Getenv("OKX_MIN_PX")
	if minPx == "" {
		log.Fatal("missing env: OKX_MIN_PX")
	}
	gridNum := os.Getenv("OKX_GRID_NUM")
	if gridNum == "" {
		log.Fatal("missing env: OKX_GRID_NUM")
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

	svc := c.NewTradingBotGridOrderAlgoService().
		InstId(instId).
		AlgoOrdType(algoOrdType).
		MaxPx(maxPx).
		MinPx(minPx).
		GridNum(gridNum)

	if v := os.Getenv("OKX_RUN_TYPE"); v != "" {
		svc.RunType(v)
	}

	switch algoOrdType {
	case "grid":
		quoteSz := os.Getenv("OKX_QUOTE_SZ")
		baseSz := os.Getenv("OKX_BASE_SZ")
		if quoteSz == "" && baseSz == "" {
			log.Fatal("missing env: OKX_QUOTE_SZ or OKX_BASE_SZ")
		}
		if quoteSz != "" {
			svc.QuoteSz(quoteSz)
		}
		if baseSz != "" {
			svc.BaseSz(baseSz)
		}
	case "contract_grid":
		sz := os.Getenv("OKX_SZ")
		direction := os.Getenv("OKX_DIRECTION")
		lever := os.Getenv("OKX_LEVER")
		if sz == "" || direction == "" || lever == "" {
			log.Fatal("missing env: OKX_SZ / OKX_DIRECTION / OKX_LEVER (required for contract_grid)")
		}
		svc.Sz(sz).Direction(direction).Lever(lever)
	default:
		log.Fatalf("unsupported OKX_ALGO_ORD_TYPE=%q (expected grid or contract_grid)", algoOrdType)
	}

	ack, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("tradingBot grid order placed: algoId=%s algoClOrdId=%s", ack.AlgoId, ack.AlgoClOrdId)
}
