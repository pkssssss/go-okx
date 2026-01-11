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

	svc := c.NewSprdTradesService()
	if v := os.Getenv("OKX_SPRD_ID"); v != "" {
		svc.SprdId(v)
	}
	if v := os.Getenv("OKX_TRADE_ID"); v != "" {
		svc.TradeId(v)
	}
	if v := os.Getenv("OKX_ORD_ID"); v != "" {
		svc.OrdId(v)
	}
	if v := os.Getenv("OKX_BEGIN_ID"); v != "" {
		svc.BeginId(v)
	}
	if v := os.Getenv("OKX_END_ID"); v != "" {
		svc.EndId(v)
	}
	if v := os.Getenv("OKX_BEGIN"); v != "" {
		svc.Begin(v)
	}
	if v := os.Getenv("OKX_END"); v != "" {
		svc.End(v)
	}
	if v := os.Getenv("OKX_LIMIT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid OKX_LIMIT: %v", err)
		}
		svc.Limit(n)
	}

	items, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("sprd trades=%d", len(items))
	if len(items) > 0 {
		tr := items[0]
		log.Printf("first sprdId=%s tradeId=%s ordId=%s side=%s fillPx=%s fillSz=%s ts=%d legs=%d", tr.SprdId, tr.TradeId, tr.OrdId, tr.Side, tr.FillPx, tr.FillSz, tr.TS, len(tr.Legs))
	}
}
