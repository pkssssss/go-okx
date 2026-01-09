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
	if os.Getenv("OKX_CONFIRM_SET_AUTO_REPAY") != "YES" {
		log.Fatal("refusing to set auto repay; set OKX_CONFIRM_SET_AUTO_REPAY=YES to continue")
	}

	autoRepayStr := os.Getenv("OKX_AUTO_REPAY")
	if autoRepayStr == "" {
		log.Fatal("missing env: OKX_AUTO_REPAY")
	}
	autoRepay, err := strconv.ParseBool(autoRepayStr)
	if err != nil {
		log.Fatalf("invalid env OKX_AUTO_REPAY: %v", err)
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

	ack, err := c.NewAccountSetAutoRepayService().AutoRepay(autoRepay).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_set_auto_repay: autoRepay=%v", ack.AutoRepay)
}
