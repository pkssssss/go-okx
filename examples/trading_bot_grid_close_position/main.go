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
		log.Fatal("refusing to close tradingBot grid position; set OKX_CONFIRM=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	algoId := os.Getenv("OKX_ALGO_ID")
	mktCloseStr := os.Getenv("OKX_MKT_CLOSE")
	if algoId == "" || mktCloseStr == "" {
		log.Fatal("missing env: OKX_ALGO_ID / OKX_MKT_CLOSE (true/false)")
	}
	mktClose, err := strconv.ParseBool(mktCloseStr)
	if err != nil {
		log.Fatalf("invalid env OKX_MKT_CLOSE=%q: %v", mktCloseStr, err)
	}

	sz := os.Getenv("OKX_SZ")
	px := os.Getenv("OKX_PX")
	if !mktClose && (sz == "" || px == "") {
		log.Fatal("missing env: OKX_SZ / OKX_PX (required when OKX_MKT_CLOSE=false)")
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

	svc := c.NewTradingBotGridClosePositionService().AlgoId(algoId).MktClose(mktClose)
	if !mktClose {
		svc.Sz(sz).Px(px)
	}

	ack, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("tradingBot grid close-position: algoId=%s algoClOrdId=%s ordId=%s tag=%s", ack.AlgoId, ack.AlgoClOrdId, ack.OrdId, ack.Tag)
}
