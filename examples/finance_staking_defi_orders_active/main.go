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

	productId := os.Getenv("OKX_PRODUCT_ID")
	protocolType := os.Getenv("OKX_PROTOCOL_TYPE")
	ccy := os.Getenv("OKX_CCY")
	state := os.Getenv("OKX_STATE")
	demo := os.Getenv("OKX_DEMO") == "1"

	c := okx.NewClient(
		okx.WithCredentials(okx.Credentials{APIKey: apiKey, SecretKey: secretKey, Passphrase: passphrase}),
		okx.WithDemoTrading(demo),
	)

	if _, err := c.SyncTime(context.Background()); err != nil {
		log.Fatal(err)
	}

	svc := c.NewFinanceStakingDefiOrdersActiveService()
	if productId != "" {
		svc.ProductId(productId)
	}
	if protocolType != "" {
		svc.ProtocolType(protocolType)
	}
	if ccy != "" {
		svc.Ccy(ccy)
	}
	if state != "" {
		svc.State(state)
	}

	items, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("orders-active=%d", len(items))
}
