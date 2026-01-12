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
		instId = "BTC-USD-SWAP"
	}

	c := okx.NewClient()

	tk, err := c.NewMarketBlockTickerService().InstId(instId).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("instType=%s instId=%s volCcy24h=%s vol24h=%s ts=%d\n", tk.InstType, tk.InstId, tk.VolCcy24h, tk.Vol24h, tk.TS)
}
