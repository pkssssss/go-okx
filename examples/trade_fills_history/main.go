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
	instType := os.Getenv("OKX_INST_TYPE")
	if instType == "" {
		instType = "SPOT"
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

	svc := c.NewTradeFillsHistoryService().InstType(instType)
	if v := os.Getenv("OKX_INST_FAMILY"); v != "" {
		svc.InstFamily(v)
	}
	if v := os.Getenv("OKX_INST_ID"); v != "" {
		svc.InstId(v)
	}
	if v := os.Getenv("OKX_ORD_ID"); v != "" {
		svc.OrdId(v)
	}
	if v := os.Getenv("OKX_SUB_TYPE"); v != "" {
		svc.SubType(v)
	}
	if v := os.Getenv("OKX_AFTER"); v != "" {
		svc.After(v)
	}
	if v := os.Getenv("OKX_BEFORE"); v != "" {
		svc.Before(v)
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

	log.Printf("instType=%s fills=%d", instType, len(items))
	if len(items) > 0 {
		f := items[0]
		log.Printf("first instId=%s ordId=%s tradeId=%s side=%s fillPx=%s fillSz=%s fee=%s %s", f.InstId, f.OrdId, f.TradeId, f.Side, f.FillPx, f.FillSz, f.Fee, f.FeeCcy)
	}
}
