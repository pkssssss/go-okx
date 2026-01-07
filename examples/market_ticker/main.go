package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pkssssss/go-okx"
)

func main() {
	instId := os.Getenv("OKX_INST_ID")
	if instId == "" {
		instId = "BTC-USDT"
	}

	c := okx.NewClient()

	tk, err := c.NewMarketTickerService().InstId(instId).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("instId=%s last=%s bid=%s/%s ask=%s/%s ts=%d\n", tk.InstId, tk.Last, tk.BidPx, tk.BidSz, tk.AskPx, tk.AskSz, tk.TS)
}
