package main

import (
	"context"
	"fmt"
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

	signalSourceType := os.Getenv("OKX_SIGNAL_SOURCE_TYPE")
	if signalSourceType == "" {
		signalSourceType = "1"
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

	svc := c.NewTradingBotSignalSignalsService().SignalSourceType(signalSourceType)
	if v := os.Getenv("OKX_SIGNAL_CHAN_ID"); v != "" {
		svc.SignalChanId(v)
	}

	items, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("signalSourceType=%s count=%d\n", signalSourceType, len(items))
	if len(items) == 0 {
		return
	}
	it := items[0]
	fmt.Printf("first signalChanId=%s signalChanName=%s signalSourceType=%s\n", it.SignalChanId, it.SignalChanName, it.SignalSourceType)
}
