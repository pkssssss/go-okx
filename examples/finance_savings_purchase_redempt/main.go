package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to purchase/redempt savings; set OKX_CONFIRM=YES to continue")
	}

	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}

	ccy := os.Getenv("OKX_CCY")
	amt := os.Getenv("OKX_AMT")
	side := os.Getenv("OKX_SIDE")
	if ccy == "" || amt == "" || side == "" {
		log.Fatal("missing env: OKX_CCY / OKX_AMT / OKX_SIDE")
	}
	rate := os.Getenv("OKX_RATE")

	demo := os.Getenv("OKX_DEMO") == "1"

	c := okx.NewClient(
		okx.WithCredentials(okx.Credentials{APIKey: apiKey, SecretKey: secretKey, Passphrase: passphrase}),
		okx.WithDemoTrading(demo),
	)

	if _, err := c.SyncTime(context.Background()); err != nil {
		log.Fatal(err)
	}

	svc := c.NewFinanceSavingsPurchaseRedemptService().
		Ccy(ccy).
		Amt(amt).
		Side(side)
	if rate != "" {
		svc.Rate(rate)
	}

	ack, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ccy=%s amt=%s side=%s rate=%s", ack.Ccy, ack.Amt, ack.Side, ack.Rate)
}
