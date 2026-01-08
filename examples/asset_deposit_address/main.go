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

	ccy := os.Getenv("OKX_CCY")
	if ccy == "" {
		log.Fatal("missing env: OKX_CCY")
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

	addrs, err := c.NewAssetDepositAddressService().Ccy(ccy).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("deposit addresses=%d", len(addrs))
	for i := 0; i < len(addrs) && i < 5; i++ {
		a := addrs[i]
		log.Printf("addr[%d]: ccy=%s chain=%s to=%s addr=%s tag=%s selected=%v", i, a.Ccy, a.Chain, a.To, a.Addr, a.Tag, a.Selected)
	}
}
