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

	svc := c.NewAssetValuationService()
	if ccy != "" {
		svc.Ccy(ccy)
	}

	v, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("asset_valuation: unit=%s totalBal=%s ts=%d funding=%s trading=%s earn=%s classic=%s", ccy, v.TotalBal, v.TS, v.Details.Funding, v.Details.Trading, v.Details.Earn, v.Details.Classic)
}
