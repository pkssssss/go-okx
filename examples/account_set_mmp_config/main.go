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
	if os.Getenv("OKX_CONFIRM_ACCOUNT_SET_MMP_CONFIG") != "YES" {
		log.Fatal("refusing to set mmp config; set OKX_CONFIRM_ACCOUNT_SET_MMP_CONFIG=YES to continue")
	}

	instFamily := os.Getenv("OKX_INST_FAMILY")
	timeInterval := os.Getenv("OKX_TIME_INTERVAL")
	frozenInterval := os.Getenv("OKX_FROZEN_INTERVAL")
	qtyLimit := os.Getenv("OKX_QTY_LIMIT")
	if instFamily == "" || timeInterval == "" || frozenInterval == "" || qtyLimit == "" {
		log.Fatal("missing env: OKX_INST_FAMILY / OKX_TIME_INTERVAL / OKX_FROZEN_INTERVAL / OKX_QTY_LIMIT")
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

	ack, err := c.NewAccountSetMMPConfigService().
		InstFamily(instFamily).
		TimeInterval(timeInterval).
		FrozenInterval(frozenInterval).
		QtyLimit(qtyLimit).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_set_mmp_config: instFamily=%s timeInterval=%s frozenInterval=%s qtyLimit=%s", ack.InstFamily, ack.TimeInterval, ack.FrozenInterval, ack.QtyLimit)
}
