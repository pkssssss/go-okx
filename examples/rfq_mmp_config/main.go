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

	cfg, err := c.NewRFQMMPConfigService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("rfq mmp config=%d", len(cfg))
	if len(cfg) > 0 {
		log.Printf("first: timeInterval=%s frozenInterval=%s countLimit=%s mmpFrozen=%v", cfg[0].TimeInterval, cfg[0].FrozenInterval, cfg[0].CountLimit, cfg[0].MMPFrozen)
	}
}
