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
		instType = "SWAP"
	}

	instFamily := os.Getenv("OKX_INST_FAMILY")

	c := okx.NewClient()

	svc := c.NewMarketBlockTickersService().InstType(instType)
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
	fmt.Printf("first instType=%s instId=%s volCcy24h=%s vol24h=%s ts=%d\n", it.InstType, it.InstId, it.VolCcy24h, it.Vol24h, it.TS)
}
