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

	instFamily := os.Getenv("OKX_INST_FAMILY") // optional
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

	svc := c.NewAccountMMPConfigService()
	if instFamily != "" {
		svc.InstFamily(instFamily)
	}

	data, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_mmp_config: items=%d", len(data))
	for _, it := range data {
		log.Printf("instFamily=%s frozen=%v frozenUntil=%s timeInterval=%s frozenInterval=%s qtyLimit=%s", it.InstFamily, it.MMPFrozen, it.MMPFrozenUntil, it.TimeInterval, it.FrozenInterval, it.QtyLimit)
	}
}
