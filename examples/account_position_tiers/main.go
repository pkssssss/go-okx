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
		instType = "SWAP"
	}

	instFamily := os.Getenv("OKX_INST_FAMILY")
	uly := os.Getenv("OKX_ULY")
	if instFamily == "" && uly == "" {
		log.Fatal("missing env: OKX_INST_FAMILY (or OKX_ULY)")
	}

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

	svc := c.NewAccountPositionTiersService().InstType(instType)
	if instFamily != "" {
		svc.InstFamily(instFamily)
	} else {
		svc.Uly(uly)
	}

	data, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_position_tiers: items=%d", len(data))
	for _, it := range data {
		log.Printf("instFamily=%s uly=%s maxSz=%s posType=%s", it.InstFamily, it.Uly, it.MaxSz, it.PosType)
	}
}
