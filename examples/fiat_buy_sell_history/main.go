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

	svc := c.NewFiatBuySellHistoryService().
		OrdId(os.Getenv("OKX_ORD_ID")).
		ClOrdId(os.Getenv("OKX_CL_ORD_ID")).
		State(os.Getenv("OKX_STATE")).
		Begin(os.Getenv("OKX_BEGIN")).
		End(os.Getenv("OKX_END"))

	if v := os.Getenv("OKX_LIMIT"); v != "" {
		limit, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid env OKX_LIMIT: %v", err)
		}
		svc.Limit(limit)
	}

	history, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("fiat_buy_sell_history: n=%d", len(history))
	for i, o := range history {
		log.Printf("order[%d]: ordId=%s clOrdId=%s state=%s cTime=%d uTime=%d", i, o.OrdId, o.ClOrdId, o.State, o.CTime, o.UTime)
	}
}
