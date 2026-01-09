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
	if os.Getenv("OKX_CONFIRM_ACCOUNT_MMP_RESET") != "YES" {
		log.Fatal("refusing to reset mmp; set OKX_CONFIRM_ACCOUNT_MMP_RESET=YES to continue")
	}

	instFamily := os.Getenv("OKX_INST_FAMILY")
	if instFamily == "" {
		log.Fatal("missing env: OKX_INST_FAMILY")
	}

	instType := os.Getenv("OKX_INST_TYPE")
	if instType == "" {
		instType = "OPTION"
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

	ack, err := c.NewAccountMMPResetService().InstType(instType).InstFamily(instFamily).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_mmp_reset: result=%v", ack.Result)
}
