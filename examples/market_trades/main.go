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

	limit := 20
	if v := os.Getenv("OKX_TRADES_LIMIT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid OKX_TRADES_LIMIT: %v", err)
		}
		limit = n
	}

	c := okx.NewClient()

	items, err := c.NewMarketTradesService().InstId(instId).Limit(limit).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("instId=%s count=%d\n", instId, len(items))
	if len(items) == 0 {
		return
	}

	it := items[0]
	fmt.Printf("first tradeId=%s px=%s sz=%s side=%s ts=%d\n", it.TradeId, it.Px, it.Sz, it.Side, it.TS)
}
