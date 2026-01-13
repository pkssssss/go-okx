package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	ccy := os.Getenv("OKX_CCY")

	c := okx.NewClient()

	svc := c.NewFinanceSavingsLendingRateSummaryService()
	if ccy != "" {
		svc.Ccy(ccy)
	}

	items, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("lending-rate-summary=%d", len(items))
}
