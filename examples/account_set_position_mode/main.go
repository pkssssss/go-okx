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
	if os.Getenv("OKX_CONFIRM_SET_POSITION_MODE") != "YES" {
		log.Fatal("refusing to set position mode; set OKX_CONFIRM_SET_POSITION_MODE=YES to continue")
	}

	posMode := os.Getenv("OKX_POS_MODE")
	if posMode == "" {
		log.Fatal("missing env: OKX_POS_MODE (long_short_mode / net_mode)")
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

	ack, err := c.NewAccountSetPositionModeService().PosMode(posMode).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("set_position_mode: posMode=%s", ack.PosMode)
}
