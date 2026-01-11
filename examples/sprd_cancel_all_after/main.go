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
	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to set sprd cancel-all-after; set OKX_CONFIRM=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	timeOut := os.Getenv("OKX_TIME_OUT")
	if timeOut == "" {
		log.Fatal("missing env: OKX_TIME_OUT (seconds; 0 or 10-120)")
	}

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

	ack, err := c.NewSprdCancelAllAfterService().TimeOut(timeOut).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("sprd_cancel_all_after: timeOut=%s triggerTime=%d ts=%d", timeOut, ack.TriggerTime, ack.TS)
}
