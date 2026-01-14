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

	algoId := os.Getenv("OKX_ALGO_ID")
	if algoId == "" {
		log.Fatal("missing env: OKX_ALGO_ID")
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

	items, err := c.NewTradingBotSignalEventHistoryService().AlgoId(algoId).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("algoId=%s count=%d\n", algoId, len(items))
	if len(items) == 0 {
		return
	}
	it := items[0]
	fmt.Printf("first eventType=%s eventStatus=%s eventCtime=%d triggeredOrdData=%d\n", it.EventType, it.EventStatus, it.EventCtime, len(it.TriggeredOrdData))
}
