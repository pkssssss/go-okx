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
	if os.Getenv("OKX_CONFIRM_ACCOUNT_LEVEL_SWITCH_PRESET") != "YES" {
		log.Fatal("refusing to preset account level switch; set OKX_CONFIRM_ACCOUNT_LEVEL_SWITCH_PRESET=YES to continue")
	}

	acctLv := os.Getenv("OKX_ACCT_LV")
	if acctLv == "" {
		log.Fatal("missing env: OKX_ACCT_LV")
	}

	lever := os.Getenv("OKX_LEVER")                     // optional
	riskOffsetType := os.Getenv("OKX_RISK_OFFSET_TYPE") // deprecated/optional

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

	svc := c.NewAccountLevelSwitchPresetService().AcctLv(acctLv)
	if lever != "" {
		svc.Lever(lever)
	}
	if riskOffsetType != "" {
		svc.RiskOffsetType(riskOffsetType)
	}

	ack, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_level_switch_preset: curAcctLv=%s acctLv=%s lever=%s", ack.CurAcctLv, ack.AcctLv, ack.Lever)
}
