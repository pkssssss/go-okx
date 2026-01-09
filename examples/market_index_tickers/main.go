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
	quoteCcy := os.Getenv("OKX_QUOTE_CCY")
	if instId == "" && quoteCcy == "" {
		instId = "BTC-USDT"
	}

	c := okx.NewClient()

	svc := c.NewMarketIndexTickersService()
	if quoteCcy != "" {
		svc.QuoteCcy(quoteCcy)
	}
	if instId != "" {
		svc.InstId(instId)
	}

	items, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("instId=%s quoteCcy=%s count=%d\n", instId, quoteCcy, len(items))
	if len(items) == 0 {
		return
	}

	it := items[0]
	fmt.Printf("first instId=%s idxPx=%s high24h=%s low24h=%s open24h=%s ts=%d\n", it.InstId, it.IdxPx, it.High24h, it.Low24h, it.Open24h, it.TS)
}
