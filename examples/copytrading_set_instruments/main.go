package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to set instruments; set OKX_CONFIRM=YES to continue")
	}

	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}

	instId := os.Getenv("OKX_INST_ID")
	if instId == "" {
		log.Fatal("missing env: OKX_INST_ID (comma-separated list)")
	}

	instType := os.Getenv("OKX_INST_TYPE")
	if instType == "" {
		instType = "SWAP"
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

	items, err := c.NewCopyTradingSetInstrumentsService().
		InstType(instType).
		InstId(instId).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("copytrading_set_instruments: items=%d", len(items))
	for i := 0; i < len(items) && i < 5; i++ {
		it := items[i]
		log.Printf("item[%d]: instId=%s enabled=%v", i, it.InstId, it.Enabled)
	}
}
