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
	if os.Getenv("OKX_CONFIRM_SET_TRADING_CONFIG") != "YES" {
		log.Fatal("refusing to set trading config; set OKX_CONFIRM_SET_TRADING_CONFIG=YES to continue")
	}

	typ := os.Getenv("OKX_TYPE")
	if typ == "" {
		log.Fatal("missing env: OKX_TYPE")
	}
	stgyType := os.Getenv("OKX_STGY_TYPE") // required when type=stgyType
	if typ == "stgyType" && stgyType == "" {
		log.Fatal("missing env: OKX_STGY_TYPE (required when OKX_TYPE=stgyType)")
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

	svc := c.NewAccountSetTradingConfigService().Type(typ)
	if stgyType != "" {
		svc.StgyType(stgyType)
	}

	ack, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_set_trading_config: type=%s stgyType=%s", ack.Type, ack.StgyType)
}
