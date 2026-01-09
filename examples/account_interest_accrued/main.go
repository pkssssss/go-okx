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

	borrowType := os.Getenv("OKX_BORROW_TYPE")
	ccy := os.Getenv("OKX_CCY")
	instId := os.Getenv("OKX_INST_ID")
	mgnMode := os.Getenv("OKX_MGN_MODE")
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

	svc := c.NewAccountInterestAccruedService()
	if borrowType != "" {
		svc.Type(borrowType)
	}
	if ccy != "" {
		svc.Ccy(ccy)
	}
	if instId != "" {
		svc.InstId(instId)
	}
	if mgnMode != "" {
		svc.MgnMode(mgnMode)
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

	log.Printf("account_interest_accrued: items=%d", len(data))
	for _, it := range data {
		log.Printf("type=%s ccy=%s instId=%s mgnMode=%s interest=%s interestRate=%s liab=%s ts=%d", it.Type, it.Ccy, it.InstId, it.MgnMode, it.Interest, it.InterestRate, it.Liab, it.TS)
	}
}
