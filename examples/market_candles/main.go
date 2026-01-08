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
	instId := os.Getenv("OKX_INST_ID")
	if instId == "" {
		instId = "BTC-USDT"
	}

	bar := os.Getenv("OKX_CANDLE_BAR")
	if bar == "" {
		bar = "1m"
	}

	limit := 5
	if v := os.Getenv("OKX_CANDLE_LIMIT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid OKX_CANDLE_LIMIT: %v", err)
		}
		limit = n
	}

	c := okx.NewClient()

	items, err := c.NewMarketCandlesService().InstId(instId).Bar(bar).Limit(limit).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("instId=%s bar=%s count=%d\n", instId, bar, len(items))
	if len(items) == 0 {
		return
	}

	it := items[0]
	fmt.Printf("first ts=%d o=%s h=%s l=%s c=%s vol=%s confirm=%s\n", it.TS, it.Open, it.High, it.Low, it.Close, it.Vol, it.Confirm)
}
