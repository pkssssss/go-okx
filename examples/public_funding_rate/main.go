package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	instId := os.Getenv("OKX_INST_ID")
	if instId == "" {
		instId = "BTC-USDT-SWAP"
	}

	c := okx.NewClient()

	fr, err := c.NewPublicFundingRateService().InstId(instId).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	nextFunding := time.UnixMilli(fr.NextFundingTime).UTC()
	funding := time.UnixMilli(fr.FundingTime).UTC()

	fmt.Printf("instId=%s fundingRate=%s fundingTime=%s nextFundingTime=%s ts=%d\n", fr.InstId, fr.FundingRate, funding.Format(time.RFC3339), nextFunding.Format(time.RFC3339), fr.TS)
}
