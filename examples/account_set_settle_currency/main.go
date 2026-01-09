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
	if os.Getenv("OKX_CONFIRM_SET_SETTLE_CURRENCY") != "YES" {
		log.Fatal("refusing to set settle currency; set OKX_CONFIRM_SET_SETTLE_CURRENCY=YES to continue")
	}

	settleCcy := os.Getenv("OKX_SETTLE_CCY")
	if settleCcy == "" {
		log.Fatal("missing env: OKX_SETTLE_CCY")
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

	ack, err := c.NewAccountSetSettleCurrencyService().SettleCcy(settleCcy).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_set_settle_currency: settleCcy=%s", ack.SettleCcy)
}
