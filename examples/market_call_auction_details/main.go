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
		instId = "ONDO-USDC"
	}

	c := okx.NewClient()

	d, err := c.NewMarketCallAuctionDetailsService().InstId(instId).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("instId=%s state=%s eqPx=%s matchedSz=%s unmatchedSz=%s auctionEndTime=%d ts=%d\n",
		d.InstId, d.State, d.EqPx, d.MatchedSz, d.UnmatchedSz, d.AuctionEndTime, d.TS)
}
