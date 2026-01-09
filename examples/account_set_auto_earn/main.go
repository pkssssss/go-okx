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
	if os.Getenv("OKX_CONFIRM_SET_AUTO_EARN") != "YES" {
		log.Fatal("refusing to set auto earn; set OKX_CONFIRM_SET_AUTO_EARN=YES to continue")
	}

	ccy := os.Getenv("OKX_CCY")
	action := os.Getenv("OKX_ACTION") // turn_on / turn_off
	if ccy == "" || action == "" {
		log.Fatal("missing env: OKX_CCY / OKX_ACTION (turn_on / turn_off)")
	}

	earnType := os.Getenv("OKX_EARN_TYPE") // optional: 0/1
	apr := os.Getenv("OKX_APR")            // optional (deprecated)
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

	svc := c.NewAccountSetAutoEarnService().Ccy(ccy).Action(action)
	if earnType != "" {
		svc.EarnType(earnType)
	}
	if apr != "" {
		svc.Apr(apr)
	}

	ack, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_set_auto_earn: earnType=%s ccy=%s action=%s apr=%s", ack.EarnType, ack.Ccy, ack.Action, ack.Apr)
}
