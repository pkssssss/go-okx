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
	if os.Getenv("OKX_CONFIRM_SET_ISOLATED_MODE") != "YES" {
		log.Fatal("refusing to set isolated mode; set OKX_CONFIRM_SET_ISOLATED_MODE=YES to continue")
	}

	isoMode := os.Getenv("OKX_ISO_MODE") // auto_transfers_ccy / automatic
	bizType := os.Getenv("OKX_BIZ_TYPE") // MARGIN / CONTRACTS
	if isoMode == "" || bizType == "" {
		log.Fatal("missing env: OKX_ISO_MODE / OKX_BIZ_TYPE (MARGIN / CONTRACTS)")
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

	ack, err := c.NewAccountSetIsolatedModeService().IsoMode(isoMode).Type(bizType).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_set_isolated_mode: isoMode=%s", ack.IsoMode)
}
