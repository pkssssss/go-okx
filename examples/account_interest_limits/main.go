package main

import (
	"context"
	"log"
	"os"

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

	svc := c.NewAccountInterestLimitsService()
	if borrowType != "" {
		svc.Type(borrowType)
	}
	if ccy != "" {
		svc.Ccy(ccy)
	}

	data, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_interest_limits: items=%d", len(data))
	for _, it := range data {
		log.Printf("debt=%s interest=%s nextDiscountTime=%d nextInterestTime=%d records=%d", it.Debt, it.Interest, it.NextDiscountTime, it.NextInterestTime, len(it.Records))
	}
}
