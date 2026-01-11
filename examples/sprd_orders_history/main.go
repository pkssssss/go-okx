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

	svc := c.NewSprdOrdersHistoryService()
	if v := os.Getenv("OKX_SPRD_ID"); v != "" {
		svc.SprdId(v)
	}
	if v := os.Getenv("OKX_ORD_TYPE"); v != "" {
		svc.OrdType(v)
	}
	if v := os.Getenv("OKX_STATE"); v != "" {
		svc.State(v)
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

	log.Printf("sprd orders history=%d", len(items))
	if len(items) > 0 {
		o := items[0]
		log.Printf("first sprdId=%s ordId=%s side=%s ordType=%s px=%s sz=%s state=%s", o.SprdId, o.OrdId, o.Side, o.OrdType, o.Px, o.Sz, o.State)
	}
}
