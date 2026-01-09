package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}

	ccy := os.Getenv("OKX_CCY") // optional, "BTC,ETH"

	var enabledPtr *bool
	enabledStr := os.Getenv("OKX_COLLATERAL_ENABLED") // optional
	if enabledStr != "" {
		v, err := strconv.ParseBool(enabledStr)
		if err != nil {
			log.Fatalf("invalid env OKX_COLLATERAL_ENABLED: %v", err)
		}
		enabledPtr = &v
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

	svc := c.NewAccountCollateralAssetsService()
	if ccy != "" {
		svc.Ccy(ccy)
	}
	if enabledPtr != nil {
		svc.CollateralEnabled(*enabledPtr)
	}

	data, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_collateral_assets: items=%d", len(data))
	for _, it := range data {
		log.Printf("ccy=%s enabled=%v", it.Ccy, it.CollateralEnabled)
	}
}
