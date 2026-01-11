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
		log.Fatal("refusing to cancel sprd order; set OKX_CONFIRM=YES to continue")
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

	ack, err := c.NewSprdCancelOrderService().
		OrdId(ordId).
		ClOrdId(clOrdId).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("sprd cancel order ack: clOrdId=%s ordId=%s sCode=%s sMsg=%s ts=%d", ack.ClOrdId, ack.OrdId, ack.SCode, ack.SMsg, ack.TS)
}
