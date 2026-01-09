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
	eventType := os.Getenv("OKX_TYPE") // auto_borrow/auto_repay/manual_borrow/manual_repay
	after := os.Getenv("OKX_AFTER")
	before := os.Getenv("OKX_BEFORE")
	limitStr := os.Getenv("OKX_LIMIT")

	c := okx.NewClient(
		okx.WithCredentials(okx.Credentials{
			APIKey:     apiKey,
			SecretKey:  secretKey,
			Passphrase: passphrase,
		}),
		okx.WithDemoTrading(demo),
	)

	if _, err := c.SyncTime(context.Background()); err != nil {
		log.Fatal(err)
	}

	svc := c.NewAccountSpotBorrowRepayHistoryService()
	if ccy != "" {
		svc.Ccy(ccy)
	}
	if eventType != "" {
		svc.Type(eventType)
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
			log.Fatalf("invalid env OKX_LIMIT: %v", err)
		}
		svc.Limit(limit)
	}

	data, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_spot_borrow_repay_history: items=%d", len(data))
	for _, it := range data {
		log.Printf("ccy=%s type=%s amt=%s accBorrowed=%s ts=%d", it.Ccy, it.Type, it.Amt, it.AccBorrowed, it.TS)
	}
}
