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

	demo := os.Getenv("OKX_DEMO") == "1"

	ordId := os.Getenv("OKX_ORD_ID")
	clOrdId := os.Getenv("OKX_CL_ORD_ID")
	if ordId == "" && clOrdId == "" {
		log.Fatal("missing env: OKX_ORD_ID or OKX_CL_ORD_ID")
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

	o, err := c.NewSprdGetOrderService().
		OrdId(ordId).
		ClOrdId(clOrdId).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("sprd order: sprdId=%s ordId=%s clOrdId=%s side=%s ordType=%s px=%s sz=%s state=%s accFillSz=%s avgPx=%s uTime=%d",
		o.SprdId, o.OrdId, o.ClOrdId, o.Side, o.OrdType, o.Px, o.Sz, o.State, o.AccFillSz, o.AvgPx, o.UTime,
	)
}
