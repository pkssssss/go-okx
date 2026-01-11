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
	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to set rfq mmp config; set OKX_CONFIRM=YES to continue")
	}

	timeInterval := os.Getenv("OKX_TIME_INTERVAL")
	frozenInterval := os.Getenv("OKX_FROZEN_INTERVAL")
	countLimit := os.Getenv("OKX_COUNT_LIMIT")
	if timeInterval == "" || frozenInterval == "" || countLimit == "" {
		log.Fatal("missing env: OKX_TIME_INTERVAL / OKX_FROZEN_INTERVAL / OKX_COUNT_LIMIT")
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

	cfg, err := c.NewRFQSetMMPConfigService().
		TimeInterval(timeInterval).
		FrozenInterval(frozenInterval).
		CountLimit(countLimit).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("set rfq mmp config: timeInterval=%s frozenInterval=%s countLimit=%s", cfg.TimeInterval, cfg.FrozenInterval, cfg.CountLimit)
}
