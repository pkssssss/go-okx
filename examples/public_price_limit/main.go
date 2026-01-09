package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	instId := os.Getenv("OKX_INST_ID")
	if instId == "" {
		instId = "BTC-USDT-SWAP"
	}

	c := okx.NewClient()

	items, err := c.NewPublicPriceLimitService().InstId(instId).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("instId=%s count=%d\n", instId, len(items))
	if len(items) == 0 {
		return
	}

	it := items[0]
	fmt.Printf("instType=%s buyLmt=%s sellLmt=%s enabled=%v ts=%d\n", it.InstType, it.BuyLmt, it.SellLmt, it.Enabled, it.TS)
}
