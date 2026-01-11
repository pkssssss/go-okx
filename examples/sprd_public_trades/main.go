package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	sprdId := os.Getenv("OKX_SPRD_ID") // 可选，例如 BTC-USDT_BTC-USDT-SWAP

	c := okx.NewClient()

	svc := c.NewSprdPublicTradesService()
	if sprdId != "" {
		svc = svc.SprdId(sprdId)
	}

	trades, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("sprd public trades count=%d", len(trades))
	if len(trades) > 0 {
		t := trades[0]
		log.Printf("first: sprdId=%s tradeId=%s px=%s sz=%s side=%s ts=%d", t.SprdId, t.TradeId, t.Px, t.Sz, t.Side, t.TS)
	}
}
