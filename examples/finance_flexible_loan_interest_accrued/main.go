package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}

	demo := os.Getenv("OKX_DEMO") == "1"
	ccy := os.Getenv("OKX_CCY")
	after := os.Getenv("OKX_AFTER")
	before := os.Getenv("OKX_BEFORE")
	limitStr := os.Getenv("OKX_LIMIT")

	c := okx.NewClient(
		okx.WithCredentials(okx.Credentials{APIKey: apiKey, SecretKey: secretKey, Passphrase: passphrase}),
		okx.WithDemoTrading(demo),
	)

	if _, err := c.SyncTime(context.Background()); err != nil {
		log.Fatal(err)
	}

	svc := c.NewFinanceFlexibleLoanInterestAccruedService()
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
	log.Printf("interest-accrued=%d", len(items))
}
