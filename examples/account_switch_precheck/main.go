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

	res, err := c.NewAccountSwitchPrecheckService().AcctLv(acctLv).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(
		"account_switch_precheck: sCode=%s curAcctLv=%s acctLv=%s unmatched=%d pos=%d posTier=%d",
		res.SCode, res.CurAcctLv, res.AcctLv, len(res.UnmatchedInfoCheck), len(res.PosList), len(res.PosTierCheck),
	)
}
