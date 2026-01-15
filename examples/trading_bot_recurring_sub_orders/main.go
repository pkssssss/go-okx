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

	svc := c.NewTradingBotRecurringSubOrdersService().AlgoId(algoId)
	if v := os.Getenv("OKX_ORD_ID"); v != "" {
		svc.OrdId(v)
	}

	items, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("algoId=%s count=%d\n", algoId, len(items))
	if len(items) == 0 {
		return
	}
	it := items[0]
	fmt.Printf("first instId=%s ordId=%s state=%s side=%s ordType=%s px=%s sz=%s cTime=%d uTime=%d\n",
		it.InstId, it.OrdId, it.State, it.Side, it.OrdType, it.Px, it.Sz, it.CTime, it.UTime,
	)
}
