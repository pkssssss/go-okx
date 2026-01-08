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
	if os.Getenv("OKX_CONFIRM_CANCEL_WITHDRAWAL") != "YES" {
		log.Fatal("refusing to cancel withdrawal; set OKX_CONFIRM_CANCEL_WITHDRAWAL=YES to continue")
	}

	wdId := os.Getenv("OKX_WD_ID")
	if wdId == "" {
		log.Fatal("missing env: OKX_WD_ID")
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

	ack, err := c.NewAssetCancelWithdrawalService().WdId(wdId).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("cancel withdrawal: wdId=%s", ack.WdId)
}
