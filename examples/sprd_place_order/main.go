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
		log.Fatal("refusing to place sprd order; set OKX_CONFIRM=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	sprdId := os.Getenv("OKX_SPRD_ID")
	side := os.Getenv("OKX_SIDE")
	ordType := os.Getenv("OKX_ORD_TYPE")
	sz := os.Getenv("OKX_SZ")
	if sprdId == "" || side == "" || ordType == "" || sz == "" {
		log.Fatal("missing env: OKX_SPRD_ID / OKX_SIDE / OKX_ORD_TYPE / OKX_SZ")
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

	svc := c.NewSprdPlaceOrderService().
		SprdId(sprdId).
		Side(side).
		OrdType(ordType).
		Sz(sz)

	if v := os.Getenv("OKX_PX"); v != "" {
		svc.Px(v)
	}
	if v := os.Getenv("OKX_CL_ORD_ID"); v != "" {
		svc.ClOrdId(v)
	}
	if v := os.Getenv("OKX_TAG"); v != "" {
		svc.Tag(v)
	}

	ack, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("sprd place order ack: clOrdId=%s ordId=%s sCode=%s sMsg=%s ts=%d", ack.ClOrdId, ack.OrdId, ack.SCode, ack.SMsg, ack.TS)
}
