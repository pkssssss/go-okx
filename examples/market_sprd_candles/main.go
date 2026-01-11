package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	sprdId := os.Getenv("OKX_SPRD_ID")
	if sprdId == "" {
		sprdId = "BTC-USDT_BTC-USDT-SWAP"
	}

	bar := os.Getenv("OKX_BAR")
	if bar == "" {
		bar = "1m"
	}

	limit := 100
	if v := os.Getenv("OKX_LIMIT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid OKX_LIMIT: %v", err)
		}
		limit = n
	}

	c := okx.NewClient()

	candles, err := c.NewMarketSprdCandlesService().
		SprdId(sprdId).
		Bar(bar).
		Limit(limit).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("sprdId=%s bar=%s candles=%d\n", sprdId, bar, len(candles))
	if len(candles) > 0 {
		k := candles[0]
		fmt.Printf("first: ts=%d o=%s h=%s l=%s c=%s vol=%s confirm=%s\n", k.TS, k.Open, k.High, k.Low, k.Close, k.Vol, k.Confirm)
	}
}
