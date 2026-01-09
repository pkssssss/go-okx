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
	if os.Getenv("OKX_CONFIRM_SET_FEE_TYPE") != "YES" {
		log.Fatal("refusing to set fee type; set OKX_CONFIRM_SET_FEE_TYPE=YES to continue")
	}

	feeType := os.Getenv("OKX_FEE_TYPE") // 0/1
	if feeType == "" {
		log.Fatal("missing env: OKX_FEE_TYPE (0 / 1)")
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

	ack, err := c.NewAccountSetFeeTypeService().FeeType(feeType).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_set_fee_type: feeType=%s", ack.FeeType)
}
