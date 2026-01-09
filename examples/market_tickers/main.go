package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	instType := os.Getenv("OKX_INST_TYPE")
	if instType == "" {
		instType = "SPOT"
	}

	instFamily := os.Getenv("OKX_INST_FAMILY")

	c := okx.NewClient()

	svc := c.NewMarketTickersService().InstType(instType)
	if instFamily != "" {
		svc.InstFamily(instFamily)
	}

	items, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("instType=%s instFamily=%s count=%d\n", instType, instFamily, len(items))
	if len(items) == 0 {
		return
	}

	it := items[0]
	fmt.Printf("first instId=%s last=%s bid=%s/%s ask=%s/%s ts=%d\n", it.InstId, it.Last, it.BidPx, it.BidSz, it.AskPx, it.AskSz, it.TS)
}
