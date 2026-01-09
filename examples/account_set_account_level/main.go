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
	if os.Getenv("OKX_CONFIRM_SET_ACCOUNT_LEVEL") != "YES" {
		log.Fatal("refusing to set account level; set OKX_CONFIRM_SET_ACCOUNT_LEVEL=YES to continue")
	}

	acctLv := os.Getenv("OKX_ACCT_LV")
	if acctLv == "" {
		log.Fatal("missing env: OKX_ACCT_LV")
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

	ack, err := c.NewAccountSetAccountLevelService().AcctLv(acctLv).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_set_account_level: acctLv=%s", ack.AcctLv)
}
