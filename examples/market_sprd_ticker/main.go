package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	sprdId := os.Getenv("OKX_SPRD_ID")
	if sprdId == "" {
		sprdId = "BTC-USDT_BTC-USDT-SWAP"
	}

	c := okx.NewClient()

	ticker, err := c.NewMarketSprdTickerService().SprdId(sprdId).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("sprdId=%s last=%s lastSz=%s ask=%s/%s bid=%s/%s ts=%d\n",
		ticker.SprdId, ticker.Last, ticker.LastSz, ticker.AskPx, ticker.AskSz, ticker.BidPx, ticker.BidSz, ticker.TS)
}
