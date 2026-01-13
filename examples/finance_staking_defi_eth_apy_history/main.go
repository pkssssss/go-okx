package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	days := os.Getenv("OKX_DAYS")
	if days == "" {
		days = "7"
	}

	c := okx.NewClient()

	items, err := c.NewFinanceStakingDefiETHAPYHistoryService().Days(days).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("eth apy-history=%d days=%s", len(items), days)
}
