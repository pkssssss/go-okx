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

	sz := 5
	if v := os.Getenv("OKX_BOOKS_SZ"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid OKX_BOOKS_SZ: %v", err)
		}
		sz = n
	}

	c := okx.NewClient()

	ob, err := c.NewMarketBooksService().InstId(instId).Sz(sz).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var bestAskPx, bestAskSz, bestBidPx, bestBidSz string
	if len(ob.Asks) > 0 {
		bestAskPx, bestAskSz = ob.Asks[0].Px, ob.Asks[0].Sz
	}
	if len(ob.Bids) > 0 {
		bestBidPx, bestBidSz = ob.Bids[0].Px, ob.Bids[0].Sz
	}

	fmt.Printf("instId=%s asks=%d bids=%d bestAsk=%s/%s bestBid=%s/%s ts=%d\n", instId, len(ob.Asks), len(ob.Bids), bestAskPx, bestAskSz, bestBidPx, bestBidSz, ob.TS)
}
