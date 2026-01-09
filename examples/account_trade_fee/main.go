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

	demo := os.Getenv("OKX_DEMO") == "1"

	instType := os.Getenv("OKX_INST_TYPE")
	if instType == "" {
		instType = "SPOT"
	}

	instId := os.Getenv("OKX_INST_ID")
	instFamily := os.Getenv("OKX_INST_FAMILY")
	groupId := os.Getenv("OKX_GROUP_ID")

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

	svc := c.NewAccountTradeFeeService().InstType(instType)
	if instId != "" {
		svc.InstId(instId)
	}
	if instFamily != "" {
		svc.InstFamily(instFamily)
	}
	if groupId != "" {
		svc.GroupId(groupId)
	}

	data, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_trade_fee: items=%d", len(data))
	for _, it := range data {
		log.Printf("instType=%s level=%s maker=%s taker=%s ts=%d feeGroup=%d", it.InstType, it.Level, it.Maker, it.Taker, it.TS, len(it.FeeGroup))
	}
}
