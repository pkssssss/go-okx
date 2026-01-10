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

	svc := c.NewOrdersHistoryArchiveService().InstType(instType)
	if v := os.Getenv("OKX_INST_FAMILY"); v != "" {
		svc.InstFamily(v)
	}
	if v := os.Getenv("OKX_INST_ID"); v != "" {
		svc.InstId(v)
	}
	if v := os.Getenv("OKX_ORD_TYPE"); v != "" {
		svc.OrdType(v)
	}
	if v := os.Getenv("OKX_ORDER_STATE"); v != "" {
		svc.State(v)
	}
	if v := os.Getenv("OKX_CATEGORY"); v != "" {
		svc.Category(v)
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

	log.Printf("instType=%s orders=%d", instType, len(items))
	if len(items) > 0 {
		o := items[0]
		log.Printf("first instId=%s ordId=%s side=%s ordType=%s px=%s sz=%s state=%s", o.InstId, o.OrdId, o.Side, o.OrdType, o.Px, o.Sz, o.State)
	}
}
