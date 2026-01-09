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

	year := os.Getenv("OKX_YEAR")
	quarter := os.Getenv("OKX_QUARTER") // Q1/Q2/Q3/Q4
	if year == "" || quarter == "" {
		log.Fatal("missing env: OKX_YEAR / OKX_QUARTER")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

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

	ack, err := c.NewAccountBillsHistoryArchiveApplyService().Year(year).Quarter(quarter).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_bills_history_archive_apply: year=%s quarter=%s result=%s ts=%d", year, quarter, ack.Result, ack.TS)
}
