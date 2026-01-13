package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	ccy := os.Getenv("OKX_CCY")
	after := os.Getenv("OKX_AFTER")
	before := os.Getenv("OKX_BEFORE")
	limitStr := os.Getenv("OKX_LIMIT")

	c := okx.NewClient()

	svc := c.NewFinanceSavingsLendingRateHistoryService()
	if ccy != "" {
		svc.Ccy(ccy)
	}
	if after != "" {
		svc.After(after)
	}
	if before != "" {
		svc.Before(before)
	}
	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			log.Fatal(err)
		}
		svc.Limit(limit)
	}

	items, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("lending-rate-history=%d", len(items))
}
