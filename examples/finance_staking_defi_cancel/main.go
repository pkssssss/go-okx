package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to cancel staking-defi order; set OKX_CONFIRM=YES to continue")
	}

	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}

	ordId := os.Getenv("OKX_ORD_ID")
	protocolType := os.Getenv("OKX_PROTOCOL_TYPE")
	if ordId == "" || protocolType == "" {
		log.Fatal("missing env: OKX_ORD_ID / OKX_PROTOCOL_TYPE")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	c := okx.NewClient(
		okx.WithCredentials(okx.Credentials{APIKey: apiKey, SecretKey: secretKey, Passphrase: passphrase}),
		okx.WithDemoTrading(demo),
	)

	if _, err := c.SyncTime(context.Background()); err != nil {
		log.Fatal(err)
	}

	ack, err := c.NewFinanceStakingDefiCancelService().OrdId(ordId).ProtocolType(protocolType).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ordId=%s tag=%s", ack.OrdId, ack.Tag)
}
